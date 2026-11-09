// Author: HyperparameterSI
package dshgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/herlanai/omnix/pkg/llm"
)

func TestDeepSeekToolsOverSSE(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != "/chat/completions" || r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("request metadata")
		}
		var req struct {
			Model    string        `json:"model"`
			Messages []wireMessage `json:"messages"`
			Tools    []any         `json:"tools"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
		}
		if req.Model != "custom" || len(req.Tools) != 1 {
			t.Error("missing model/tool schema")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		if requests == 1 {
			fmt.Fprint(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call\",\"function\":{\"name\":\"echo\",\"arguments\":\"{\\\"n\\\":\"}}]}}]}\n\n")
			fmt.Fprint(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"function\":{\"arguments\":\"2}\"}}]},\"finish_reason\":\"tool_calls\"}]}\n\n")
		} else {
			last := req.Messages[len(req.Messages)-1]
			prior := req.Messages[len(req.Messages)-2]
			if last.Role != "tool" || last.CallID != "call" || last.Content != "2" || len(prior.Calls) != 1 {
				t.Errorf("invalid wire tool exchange: %+v", req.Messages)
			}
			fmt.Fprint(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"完成\"},\"finish_reason\":\"stop\"}]}\n\n")
		}
		fmt.Fprint(w, "data: {\"choices\":[],\"usage\":{\"prompt_tokens\":3,\"completion_tokens\":1}}\n\ndata: [DONE]\n\n")
	}))
	defer server.Close()
	d, e := NewDeepSeek(DeepSeekOptions{APIKey: "test-key", BaseURL: server.URL})
	if e != nil {
		t.Fatal(e)
	}
	c := mustConversation(t, Options{Adapter: d, Tools: []Tool{{Name: "echo", Execute: func(_ context.Context, in map[string]any) (any, error) { return in["n"], nil }}}})
	result, e := c.Run(context.Background(), "hello", nil)
	if e != nil {
		t.Fatal(e)
	}
	if requests != 2 || result.Text != "完成" || result.Usage.PromptTokens != 6 {
		t.Fatalf("%+v", result)
	}
}

func TestDeepSeekStreamFailures(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		status     int
	}{{"truncated", "data: {\"choices\":[]}\n\n", 200}, {"invalid", "data: broken\n\n", 200}, {"quota", "secret provider body", 429}, {"redirect", "", 302}, {"tool", "data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"x\",\"function\":{\"name\":\"tool\",\"arguments\":\"invalid\"}}]}}]}\n\ndata: [DONE]\n\n", 200}} {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(tc.status); fmt.Fprint(w, tc.body) }))
			defer srv.Close()
			d, _ := NewDeepSeek(DeepSeekOptions{BaseURL: srv.URL})
			_, err := d.Chat(context.Background(), llm.ChatRequest{Model: "test"}, nil)
			if err == nil {
				t.Fatal("accepted failed response")
			}
			if strings.Contains(err.Error(), "secret") {
				t.Fatal("leaked response")
			}
			if tc.name == "truncated" && !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Fatal(err)
			}
		})
	}
}
