// Author: HyperparameterSI

// Package dshgo embeds dsh-go's model, event and tool interfaces in a
// bounded, cancellable desktop conversation. It starts no process or server.
package dshgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/JopenChen/dsh-go/pkg/brand"
	"github.com/JopenChen/dsh-go/pkg/llm"
	"github.com/JopenChen/dsh-go/pkg/session"
	"github.com/JopenChen/dsh-go/pkg/tools"
)

var (
	ErrBusy   = errors.New("assistant: a turn is already running")
	ErrClosed = errors.New("assistant: conversation is closed")
	ErrLimit  = errors.New("assistant: turn or context limit reached")
)

// Tool exposes an explicitly registered Go function to the model. No shell,
// file access, desktop automation or other authority is registered implicitly.
type Tool struct {
	Name, Description string
	Parameters        map[string]any
	Execute           tools.ToolExecuteFunc
	RequiresApproval  bool
}

// Options is copied by New. The adapter must serialize callbacks within Chat
// and return only after its callbacks finish; adapters and tools must honor ctx.
type Options struct {
	Adapter       llm.LLMAdapter
	Model, System string
	Tools         []Tool
	// Approve runs on the turn's worker, before executing a protected tool.
	// A nil callback denies protected tools. It must respect cancellation.
	Approve func(context.Context, llm.ToolCall) (bool, error)

	MaxSteps, MaxToolCalls, MaxTokens, MaxTurnBytes, MaxHistoryBytes, HistoryTurns int
	Timeout                                                                        time.Duration
}

// Event reports incremental text and tool activity. Emit runs synchronously on
// Run's caller goroutine: use Window.Update before changing native UI state.
// Callbacks may call Cancel, but must not reenter Run or History.
type Event struct {
	Kind    string
	Text    string
	Tool    string
	CallID  string
	IsError bool
}

// Result includes a per-turn dsh-go event log, suitable for explicit persistence.
// Failed turns are not added to conversation history; tools may already have run.
type Result struct {
	Text   string
	Usage  llm.Usage
	Steps  int
	Events []session.SessionEvent
}

type registered struct {
	tool     Tool
	schema   *tools.JsonSchemaNode
	pipeline *tools.Pipeline
	impl     *tools.Tool
}

// Conversation serializes turns and owns bounded in-memory history. API keys
// belong to the adapter and are never included in History or event logs.
type Conversation struct {
	mu      sync.Mutex
	life    sync.Mutex
	cancel  context.CancelFunc
	closed  bool
	opts    Options
	tools   map[string]registered
	schemas []llm.ToolSchema
	turns   [][]llm.Message
}

// New creates a conversation without network traffic or background goroutines.
func New(o Options) (*Conversation, error) {
	if o.Adapter == nil || strings.TrimSpace(o.Model) == "" {
		return nil, errors.New("assistant: adapter and model are required")
	}
	if o.MaxSteps == 0 {
		o.MaxSteps = 8
	}
	if o.MaxToolCalls == 0 {
		o.MaxToolCalls = 16
	}
	if o.MaxTokens == 0 {
		o.MaxTokens = 2048
	}
	if o.MaxTurnBytes == 0 {
		o.MaxTurnBytes = 256 << 10
	}
	if o.MaxHistoryBytes == 0 {
		o.MaxHistoryBytes = 1 << 20
	}
	if o.HistoryTurns == 0 {
		o.HistoryTurns = 20
	}
	if o.Timeout == 0 {
		o.Timeout = 2 * time.Minute
	}
	if o.MaxSteps < 1 || o.MaxToolCalls < 1 || o.MaxTokens < 1 || o.MaxTurnBytes < 1 || o.MaxHistoryBytes < o.MaxTurnBytes || o.HistoryTurns < 1 || o.Timeout < 0 {
		return nil, errors.New("assistant: invalid limits")
	}
	c := &Conversation{opts: o, tools: make(map[string]registered)}
	for _, t := range o.Tools {
		if t.Name == "" || t.Execute == nil {
			return nil, errors.New("assistant: tool name and implementation required")
		}
		if _, ok := c.tools[t.Name]; ok {
			return nil, fmt.Errorf("assistant: duplicate tool %q", t.Name)
		}
		if t.Parameters == nil {
			t.Parameters = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		raw, err := json.Marshal(t.Parameters)
		if err != nil {
			return nil, err
		}
		var parameters map[string]any
		if err = json.Unmarshal(raw, &parameters); err != nil {
			return nil, err
		}
		t.Parameters = parameters
		schema, err := tools.Compile(raw)
		if err != nil {
			return nil, fmt.Errorf("tool %s: %w", t.Name, err)
		}
		impl := &tools.Tool{Name: t.Name, Description: t.Description, Schema: schema, Execute: t.Execute}
		c.tools[t.Name] = registered{t, schema, tools.NewPipeline().WithTool(impl), impl}
		c.schemas = append(c.schemas, llm.ToolSchema{Name: t.Name, Description: t.Description, Parameters: parameters})
	}
	c.opts.Tools = nil
	return c, nil
}

// Cancel cancels the current turn, if any, without blocking the UI thread.
func (c *Conversation) Cancel() {
	c.life.Lock()
	defer c.life.Unlock()
	if c.cancel != nil {
		c.cancel()
	}
}

// Close prevents new turns and cancels current work. Cooperative tools unwind
// asynchronously; Close never waits for a tool which might need the UI thread.
func (c *Conversation) Close() error {
	c.life.Lock()
	defer c.life.Unlock()
	c.closed = true
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

// Reset discards history when idle.
func (c *Conversation) Reset() error {
	if !c.mu.TryLock() {
		return ErrBusy
	}
	defer c.mu.Unlock()
	c.turns = nil
	return nil
}

// History returns an independent JSON-compatible copy of completed turns.
func (c *Conversation) History() []llm.Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	return clone(c.history())
}
func (c *Conversation) history() []llm.Message {
	var out []llm.Message
	for _, t := range c.turns {
		out = append(out, t...)
	}
	return out
}
func clone[T any](v T) T { b, _ := json.Marshal(v); var out T; _ = json.Unmarshal(b, &out); return out }
func size(v any) int     { b, _ := json.Marshal(v); return len(b) }

// Run performs a bounded model/tool loop. Call it from a worker goroutine.
// Only a successfully completed turn is committed; errors return partial output.
func (c *Conversation) Run(parent context.Context, input string, emit func(Event)) (out Result, err error) {
	if !c.mu.TryLock() {
		return out, ErrBusy
	}
	defer c.mu.Unlock()
	if strings.TrimSpace(input) == "" {
		return out, errors.New("assistant: empty input")
	}
	if len(input) > c.opts.MaxTurnBytes {
		return out, ErrLimit
	}
	c.life.Lock()
	if c.closed {
		c.life.Unlock()
		return out, ErrClosed
	}
	ctx, cancel := context.WithTimeout(parent, c.opts.Timeout)
	c.cancel = cancel
	c.life.Unlock()
	defer func() { cancel(); c.life.Lock(); c.cancel = nil; c.life.Unlock() }()
	if emit == nil {
		emit = func(Event) {}
	}
	log := session.NewSessionLog(brand.NewSessionID("golet-turn"))
	record := func(d session.EventData) {
		if _, e := log.Append(d); e != nil {
			panic(e)
		}
	}
	record(session.TurnStartData{Turn: 0})
	record(session.UserMessageData{Content: input})
	defer func() {
		reason := session.ReasonCompleted
		if err != nil {
			reason = session.ReasonError
			if ctx.Err() != nil {
				reason = session.ReasonAborted
			}
		}
		record(session.TurnEndData{Turn: 0, Reason: reason})
		out.Events = log.Events()
	}()
	turn := []llm.Message{llm.NewUserMessage(input)}
	used := len(input)
	calls := 0
	seen := map[string]bool{}
	for step := 1; step <= c.opts.MaxSteps; step++ {
		if err = ctx.Err(); err != nil {
			return out, err
		}
		out.Steps = step
		messages := append(c.history(), turn...)
		if size(messages)+len(c.opts.System) > c.opts.MaxHistoryBytes {
			return out, ErrLimit
		}
		record(session.StepStartData{Turn: 0, Step: uint64(step)})
		var content []llm.ContentBlock
		var pending []*llm.ToolCall
		var streamErr error
		usage, e := c.opts.Adapter.Chat(ctx, llm.ChatRequest{Model: c.opts.Model, System: c.opts.System, Messages: clone(messages), Tools: clone(c.schemas), MaxTokens: c.opts.MaxTokens}, func(chunk llm.StreamChunk) {
			if streamErr != nil {
				return
			}
			used += len(chunk.Text) + len(chunk.Reasoning)
			if chunk.ToolCall != nil {
				used += size(chunk.ToolCall)
			}
			if used > c.opts.MaxTurnBytes {
				streamErr = ErrLimit
				cancel()
				return
			}
			switch chunk.Kind {
			case llm.ChunkText:
				if len(content) > 0 && content[len(content)-1].Kind == llm.BlockText {
					content[len(content)-1].Text += chunk.Text
				} else {
					content = append(content, llm.Text(chunk.Text))
				}
				out.Text += chunk.Text
				emit(Event{Kind: "text", Text: chunk.Text})
			case llm.ChunkReasoning:
				if len(content) > 0 && content[len(content)-1].Kind == llm.BlockReasoning {
					content[len(content)-1].Text += chunk.Reasoning
				} else {
					content = append(content, llm.Reasoning(chunk.Reasoning))
				}
			case llm.ChunkToolCall:
				if chunk.ToolCall != nil {
					tc := clone(*chunk.ToolCall)
					pending = append(pending, &tc)
					content = append(content, llm.ToolUse(&tc))
				}
			}
		})
		out.Usage.PromptTokens += usage.PromptTokens
		out.Usage.CompletionTokens += usage.CompletionTokens
		out.Usage.PromptCacheHitTokens += usage.PromptCacheHitTokens
		out.Usage.PromptCacheMissTokens += usage.PromptCacheMissTokens
		if streamErr != nil {
			e = streamErr
		}
		if e == nil {
			e = ctx.Err()
		}
		if e != nil {
			record(session.StepEndData{Turn: 0, Step: uint64(step)})
			return out, e
		}
		turn = append(turn, llm.Message{Role: llm.RoleAssistant, Content: content})
		if len(pending) == 0 {
			if len(content) == 0 {
				record(session.StepEndData{Turn: 0, Step: uint64(step)})
				return out, errors.New("assistant: model returned no content")
			}
			record(session.AssistantMessageData{Content: out.Text})
			record(session.StepEndData{Turn: 0, Step: uint64(step)})
			if size(turn) > c.opts.MaxHistoryBytes {
				return out, ErrLimit
			}
			c.turns = append(c.turns, turn)
			for len(c.turns) > c.opts.HistoryTurns || size(c.history()) > c.opts.MaxHistoryBytes {
				c.turns = c.turns[1:]
			}
			return out, nil
		}
		for _, tc := range pending {
			if tc.ID == "" || seen[tc.ID] {
				record(session.StepEndData{Turn: 0, Step: uint64(step)})
				return out, errors.New("assistant: invalid tool call ID")
			}
			seen[tc.ID] = true
		}
		for _, tc := range pending {
			calls++
			if calls > c.opts.MaxToolCalls {
				record(session.StepEndData{Turn: 0, Step: uint64(step)})
				return out, ErrLimit
			}
			raw, _ := json.Marshal(tc.Input)
			id := brand.NewToolCallID(tc.ID)
			record(session.ToolCallData{CallID: id, Tool: tc.Name, Input: raw})
			emit(Event{Kind: "tool_start", Tool: tc.Name, CallID: tc.ID})
			value, failed := c.execute(ctx, *tc)
			oversized := len(value) > c.opts.MaxTurnBytes-used
			if oversized {
				value = "tool result exceeds remaining byte budget"
				failed = true
			}
			used += len(value)
			record(session.ToolResultData{CallID: id, IsError: failed, Output: value})
			emit(Event{Kind: "tool_end", Tool: tc.Name, CallID: tc.ID, IsError: failed})
			if oversized {
				record(session.StepEndData{Turn: 0, Step: uint64(step)})
				return out, ErrLimit
			}
			turn = append(turn, llm.Message{Role: llm.RoleUser, Content: []llm.ContentBlock{llm.ToolResult(&llm.ToolResultContent{ToolUseID: tc.ID, Content: value, IsError: failed})}})
		}
		record(session.StepEndData{Turn: 0, Step: uint64(step)})
	}
	return out, ErrLimit
}

func (c *Conversation) execute(ctx context.Context, tc llm.ToolCall) (text string, failed bool) {
	defer func() {
		if recover() != nil {
			text = "tool panicked"
			failed = true
		}
	}()
	t, ok := c.tools[tc.Name]
	if !ok {
		return "unknown tool", true
	}
	if err := t.schema.Validate(tc.Input); err != nil {
		return "invalid tool arguments: " + err.Error(), true
	}
	if t.tool.RequiresApproval {
		if c.opts.Approve == nil {
			return "tool approval required", true
		}
		allowed, err := c.opts.Approve(ctx, clone(tc))
		if err != nil {
			return "tool approval failed", true
		}
		if !allowed {
			return "tool approval denied", true
		}
	}
	if err := ctx.Err(); err != nil {
		return err.Error(), true
	}
	r := t.pipeline.Run(ctx, &tools.ToolCallRequest{CallID: brand.NewToolCallID(tc.ID), Tool: tc.Name, Input: tc.Input}, t.impl)
	if r.IsError {
		return r.Error, true
	}
	b, err := json.Marshal(r.Value)
	if err != nil {
		return "tool result is not JSON serializable", true
	}
	return string(b), false
}
