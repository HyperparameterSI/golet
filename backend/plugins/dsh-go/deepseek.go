// Author: HyperparameterSI

package dshgo

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/herlanai/omnix/pkg/llm"
)

// DeepSeekOptions configures the Go HTTP/SSE adapter. BaseURL may point to a
// compatible local model server. HTTPS is required except for loopback hosts.
type DeepSeekOptions struct {
	APIKey, BaseURL string
	Client          *http.Client
}

// DeepSeek implements dsh-go's LLMAdapter with native function-call messages.
type DeepSeek struct {
	key, endpoint string
	client        *http.Client
}

// NewDeepSeek creates an adapter without making a request. Supplied clients
// retain their transport and timeout, but redirects are refused to protect keys.
func NewDeepSeek(o DeepSeekOptions) (*DeepSeek, error) {
	if o.BaseURL == "" {
		o.BaseURL = "https://api.deepseek.com"
	}
	u, err := url.Parse(o.BaseURL)
	if err != nil {
		return nil, errors.New("assistant: invalid API base URL")
	}
	local := u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1" || u.Hostname() == "::1"
	if u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Scheme != "https" && !(u.Scheme == "http" && local)) {
		return nil, errors.New("assistant: use an HTTPS base URL or a loopback HTTP endpoint")
	}
	client := http.Client{Timeout: 2 * time.Minute}
	if o.Client != nil {
		client = *o.Client
	}
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &DeepSeek{o.APIKey, strings.TrimRight(o.BaseURL, "/") + "/chat/completions", &client}, nil
}
func (*DeepSeek) Name() string { return "deepseek" }

// HTTPError reports status without copying provider response bodies or secrets.
type HTTPError struct{ StatusCode int }

func (e *HTTPError) Error() string {
	return fmt.Sprintf("assistant: model HTTP status %d", e.StatusCode)
}

type wireCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function wireFunction `json:"function"`
}
type wireFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
type wireMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	Reasoning string     `json:"reasoning_content,omitempty"`
	Calls     []wireCall `json:"tool_calls,omitempty"`
	CallID    string     `json:"tool_call_id,omitempty"`
}

// Chat maps dsh-go's content blocks to the provider's native tool-call protocol.
// In particular tool results are separate role=tool messages, never plain text.
func (d *DeepSeek) Chat(ctx context.Context, r llm.ChatRequest, cb func(llm.StreamChunk)) (usage llm.Usage, err error) {
	var messages []wireMessage
	if r.System != "" {
		messages = append(messages, wireMessage{Role: "system", Content: r.System})
	}
	for _, m := range r.Messages {
		w := wireMessage{Role: string(m.Role)}
		var results []wireMessage
		for _, b := range m.Content {
			switch b.Kind {
			case llm.BlockText:
				w.Content += b.Text
			case llm.BlockReasoning:
				w.Reasoning += b.Text
			case llm.BlockToolUse:
				if b.ToolCall == nil {
					return usage, errors.New("assistant: missing tool call")
				}
				args, e := json.Marshal(b.ToolCall.Input)
				if e != nil {
					return usage, e
				}
				w.Calls = append(w.Calls, wireCall{b.ToolCall.ID, "function", wireFunction{b.ToolCall.Name, string(args)}})
			case llm.BlockToolResult:
				if b.ToolResult == nil {
					return usage, errors.New("assistant: missing tool result")
				}
				results = append(results, wireMessage{Role: "tool", Content: b.ToolResult.Content, CallID: b.ToolResult.ToolUseID})
			default:
				return usage, errors.New("assistant: this adapter supports text and tools only")
			}
		}
		if w.Content != "" || w.Reasoning != "" || len(w.Calls) > 0 || len(results) == 0 {
			messages = append(messages, w)
		}
		messages = append(messages, results...)
	}
	body := map[string]any{"model": r.Model, "messages": messages, "stream": true, "stream_options": map[string]any{"include_usage": true}}
	if r.MaxTokens > 0 {
		body["max_tokens"] = r.MaxTokens
	}
	if len(r.Tools) > 0 {
		defs := make([]any, 0, len(r.Tools))
		for _, t := range r.Tools {
			defs = append(defs, map[string]any{"type": "function", "function": map[string]any{"name": t.Name, "description": t.Description, "parameters": t.Parameters}})
		}
		body["tools"] = defs
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return usage, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.endpoint, bytes.NewReader(raw))
	if err != nil {
		return usage, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	if d.key != "" {
		req.Header.Set("Authorization", "Bearer "+d.key)
	}
	resp, err := d.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return usage, ctx.Err()
		}
		return usage, errors.New("assistant: model connection failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return usage, &HTTPError{resp.StatusCode}
	}
	if cb == nil {
		cb = func(llm.StreamChunk) {}
	}
	const maxResponse = 2 << 20
	scanner := bufio.NewScanner(io.LimitReader(resp.Body, maxResponse+1))
	scanner.Buffer(make([]byte, 4096), 256<<10)
	pending := map[int]*wireCall{}
	total := 0
	finished := false
	var data []string
	consume := func() error {
		if len(data) == 0 {
			return nil
		}
		payload := strings.Join(data, "\n")
		data = nil
		if payload == "[DONE]" {
			finished = true
			return nil
		}
		var frame struct {
			Error   json.RawMessage `json:"error"`
			Choices []struct {
				Index  int     `json:"index"`
				Finish *string `json:"finish_reason"`
				Delta  struct {
					Content   string `json:"content"`
					Reasoning string `json:"reasoning_content"`
					Calls     []struct {
						Index    int          `json:"index"`
						ID       string       `json:"id"`
						Function wireFunction `json:"function"`
					} `json:"tool_calls"`
				} `json:"delta"`
			} `json:"choices"`
			Usage *struct {
				Prompt     int `json:"prompt_tokens"`
				Completion int `json:"completion_tokens"`
				Hit        int `json:"prompt_cache_hit_tokens"`
				Miss       int `json:"prompt_cache_miss_tokens"`
			} `json:"usage"`
		}
		if e := json.Unmarshal([]byte(payload), &frame); e != nil {
			return errors.New("assistant: malformed model stream")
		}
		if len(frame.Error) > 0 && string(frame.Error) != "null" {
			return errors.New("assistant: model stream reported an error")
		}
		if frame.Usage != nil {
			usage = llm.Usage{PromptTokens: frame.Usage.Prompt, CompletionTokens: frame.Usage.Completion, PromptCacheHitTokens: frame.Usage.Hit, PromptCacheMissTokens: frame.Usage.Miss}
		}
		for _, choice := range frame.Choices {
			if choice.Index != 0 {
				continue
			}
			if choice.Finish != nil && *choice.Finish != "stop" && *choice.Finish != "tool_calls" {
				return fmt.Errorf("assistant: incomplete model response (%s)", *choice.Finish)
			}
			if choice.Delta.Content != "" {
				cb(llm.StreamChunk{Kind: llm.ChunkText, Text: choice.Delta.Content})
			}
			if choice.Delta.Reasoning != "" {
				cb(llm.StreamChunk{Kind: llm.ChunkReasoning, Reasoning: choice.Delta.Reasoning})
			}
			for _, part := range choice.Delta.Calls {
				if part.Index < 0 || part.Index >= 64 {
					return ErrLimit
				}
				v := pending[part.Index]
				if v == nil {
					v = &wireCall{}
					pending[part.Index] = v
				}
				v.ID += part.ID
				v.Function.Name += part.Function.Name
				v.Function.Arguments += part.Function.Arguments
			}
		}
		return ctx.Err()
	}
	for scanner.Scan() {
		line := scanner.Text()
		total += len(line) + 1
		if total > maxResponse {
			return usage, ErrLimit
		}
		if line == "" {
			if err = consume(); err != nil {
				return usage, err
			}
			if finished {
				break
			}
		} else if strings.HasPrefix(line, "data:") {
			data = append(data, strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
		}
	}
	if err = scanner.Err(); err != nil {
		if ctx.Err() != nil {
			return usage, ctx.Err()
		}
		return usage, errors.New("assistant: model stream read failed")
	}
	if !finished {
		if err = consume(); err != nil {
			return usage, err
		}
	}
	if !finished {
		return usage, io.ErrUnexpectedEOF
	}
	indices := make([]int, 0, len(pending))
	for i := range pending {
		indices = append(indices, i)
	}
	sort.Ints(indices)
	for _, i := range indices {
		v := pending[i]
		var input map[string]any
		if v.ID == "" || v.Function.Name == "" || json.Unmarshal([]byte(v.Function.Arguments), &input) != nil || input == nil {
			return usage, errors.New("assistant: malformed tool arguments")
		}
		cb(llm.StreamChunk{Kind: llm.ChunkToolCall, ToolCall: &llm.ToolCall{ID: v.ID, Name: v.Function.Name, Input: input}})
	}
	cb(llm.StreamChunk{Kind: llm.ChunkDone})
	return usage, ctx.Err()
}
