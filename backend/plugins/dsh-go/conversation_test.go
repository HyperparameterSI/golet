// Author: HyperparameterSI
package dshgo

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/JopenChen/dsh-go/pkg/llm"
	"github.com/JopenChen/dsh-go/pkg/session"
)

type modelFunc func(context.Context, llm.ChatRequest, func(llm.StreamChunk)) (llm.Usage, error)

func (modelFunc) Name() string { return "test" }
func (f modelFunc) Chat(c context.Context, r llm.ChatRequest, emit func(llm.StreamChunk)) (llm.Usage, error) {
	return f(c, r, emit)
}
func mustConversation(t *testing.T, o Options) *Conversation {
	t.Helper()
	if o.Model == "" {
		o.Model = "custom"
	}
	c, e := New(o)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { c.Close() })
	return c
}

func TestToolsContinueWithProtocolAndHistory(t *testing.T) {
	calls := 0
	executed := 0
	adapter := modelFunc(func(_ context.Context, r llm.ChatRequest, emit func(llm.StreamChunk)) (llm.Usage, error) {
		calls++
		if r.Model != "custom" || len(r.Tools) != 1 {
			t.Fatalf("request: %+v", r)
		}
		if calls == 1 {
			emit(llm.StreamChunk{Kind: llm.ChunkToolCall, ToolCall: &llm.ToolCall{ID: "a", Name: "sum", Input: map[string]any{"n": float64(2)}}})
		} else {
			if calls == 2 {
				last := r.Messages[len(r.Messages)-1].Content[0].ToolResult
				if last == nil || last.ToolUseID != "a" || last.Content != "3" {
					t.Fatalf("tool result: %+v", last)
				}
			}
			emit(llm.StreamChunk{Kind: llm.ChunkText, Text: "done"})
		}
		return llm.Usage{PromptTokens: 3, CompletionTokens: 1}, nil
	})
	c := mustConversation(t, Options{Adapter: adapter, Tools: []Tool{{Name: "sum", Parameters: map[string]any{"type": "object", "properties": map[string]any{"n": map[string]any{"type": "number"}}, "required": []string{"n"}}, Execute: func(_ context.Context, in map[string]any) (any, error) { executed++; return in["n"].(float64) + 1, nil }}}})
	r, e := c.Run(context.Background(), "add", nil)
	if e != nil {
		t.Fatal(e)
	}
	if r.Text != "done" || r.Steps != 2 || executed != 1 || r.Usage.PromptTokens != 6 {
		t.Fatalf("%+v executions=%d", r, executed)
	}
	if failures := session.VerifyInvariants(r.Events); len(failures) > 0 {
		t.Fatal(failures)
	}
	history := c.History()
	history[0].Content[0].Text = "mutated"
	if c.History()[0].Content[0].Text != "add" {
		t.Fatal("history aliases caller")
	}
	if _, e = c.Run(context.Background(), "next", nil); e != nil {
		t.Fatal(e)
	}
}

func TestToolGates(t *testing.T) {
	for _, kind := range []string{"approval", "schema", "unknown", "panic"} {
		t.Run(kind, func(t *testing.T) {
			executed := 0
			step := 0
			adapter := modelFunc(func(_ context.Context, r llm.ChatRequest, emit func(llm.StreamChunk)) (llm.Usage, error) {
				step++
				if step == 1 {
					name := "tool"
					if kind == "unknown" {
						name = "missing"
					}
					in := map[string]any{"n": float64(1)}
					if kind == "schema" {
						in["n"] = "bad"
					}
					emit(llm.StreamChunk{Kind: llm.ChunkToolCall, ToolCall: &llm.ToolCall{ID: "a", Name: name, Input: in}})
				} else {
					if !r.Messages[len(r.Messages)-1].Content[0].ToolResult.IsError {
						t.Fatal("failed tool not marked")
					}
					emit(llm.StreamChunk{Kind: llm.ChunkText, Text: "denied"})
				}
				return llm.Usage{}, nil
			})
			c := mustConversation(t, Options{Adapter: adapter, Tools: []Tool{{Name: "tool", RequiresApproval: kind == "approval", Parameters: map[string]any{"type": "object", "properties": map[string]any{"n": map[string]any{"type": "number"}}}, Execute: func(context.Context, map[string]any) (any, error) { executed++; panic("private failure") }}}})
			r, e := c.Run(context.Background(), "hi", nil)
			if e != nil {
				t.Fatal(e)
			}
			want := 0
			if kind == "panic" {
				want = 1
			}
			if executed != want {
				t.Fatal("unexpected execution", executed)
			}
			if failures := session.VerifyInvariants(r.Events); len(failures) > 0 {
				t.Fatal(failures)
			}
		})
	}
}

func TestCancelBusyAndClose(t *testing.T) {
	started := make(chan struct{})
	done := make(chan error, 1)
	c := mustConversation(t, Options{Adapter: modelFunc(func(ctx context.Context, _ llm.ChatRequest, _ func(llm.StreamChunk)) (llm.Usage, error) {
		close(started)
		<-ctx.Done()
		return llm.Usage{}, ctx.Err()
	})})
	go func() { _, e := c.Run(context.Background(), "hello", nil); done <- e }()
	<-started
	if _, e := c.Run(context.Background(), "other", nil); !errors.Is(e, ErrBusy) {
		t.Fatal(e)
	}
	c.Close()
	if e := <-done; !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if len(c.History()) != 0 {
		t.Fatal("cancelled turn committed")
	}
	if _, e := c.Run(context.Background(), "again", nil); !errors.Is(e, ErrClosed) {
		t.Fatal(e)
	}
}

func TestLimitsAndRecovery(t *testing.T) {
	var requests atomic.Int32
	c := mustConversation(t, Options{MaxSteps: 2, Adapter: modelFunc(func(_ context.Context, _ llm.ChatRequest, emit func(llm.StreamChunk)) (llm.Usage, error) {
		id := requests.Add(1)
		emit(llm.StreamChunk{Kind: llm.ChunkToolCall, ToolCall: &llm.ToolCall{ID: string(rune('a' + id)), Name: "unknown", Input: map[string]any{}}})
		return llm.Usage{}, nil
	})})
	r, e := c.Run(context.Background(), "loop", nil)
	if !errors.Is(e, ErrLimit) || requests.Load() != 2 {
		t.Fatalf("%v requests %d", e, requests.Load())
	}
	if len(c.History()) != 0 {
		t.Fatal("failed turn committed")
	}
	if failures := session.VerifyInvariants(r.Events); len(failures) > 0 {
		t.Fatal(failures)
	}
	bounded := mustConversation(t, Options{HistoryTurns: 1, Adapter: modelFunc(func(_ context.Context, _ llm.ChatRequest, emit func(llm.StreamChunk)) (llm.Usage, error) {
		emit(llm.StreamChunk{Kind: llm.ChunkText, Text: "ok"})
		return llm.Usage{}, nil
	})})
	for _, v := range []string{"first", "second"} {
		if _, e := bounded.Run(context.Background(), v, nil); e != nil {
			t.Fatal(e)
		}
	}
	if h := bounded.History(); len(h) != 2 || h[0].Content[0].Text != "second" {
		t.Fatal(h)
	}
}

func TestOutputBudgetAndEmptyResponse(t *testing.T) {
	for _, kind := range []string{"empty", "text", "tool"} {
		t.Run(kind, func(t *testing.T) {
			c := mustConversation(t, Options{MaxTurnBytes: 256, Adapter: modelFunc(func(_ context.Context, _ llm.ChatRequest, emit func(llm.StreamChunk)) (llm.Usage, error) {
				switch kind {
				case "text":
					emit(llm.StreamChunk{Kind: llm.ChunkText, Text: strings.Repeat("x", 300)})
				case "tool":
					emit(llm.StreamChunk{Kind: llm.ChunkToolCall, ToolCall: &llm.ToolCall{ID: "large", Name: "large", Input: map[string]any{}}})
				}
				return llm.Usage{}, nil
			}), Tools: []Tool{{Name: "large", Execute: func(context.Context, map[string]any) (any, error) { return strings.Repeat("x", 300), nil }}}})
			r, err := c.Run(context.Background(), "test", nil)
			if err == nil || (kind != "empty" && !errors.Is(err, ErrLimit)) {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(c.History()) != 0 {
				t.Fatal("failed turn committed")
			}
			if failures := session.VerifyInvariants(r.Events); len(failures) != 0 {
				t.Fatal(failures)
			}
		})
	}
}
