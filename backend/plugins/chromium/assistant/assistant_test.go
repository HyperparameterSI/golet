// Author: HyperparameterSI
package assistant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestProviderStructuredResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("missing server credential")
		}
		var body struct{ Messages []struct{ Content string } }
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			t.Fatal("invalid JSON")
		}
		if len(body.Messages) != 2 || !strings.Contains(body.Messages[1].Content, "question") {
			t.Error("missing context")
		}
		json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"message": map[string]string{"content": `{"reply":"Ready","actions":[{"kind":"click","element":0}]}`}}}})
	}))
	defer server.Close()
	p := Provider{Endpoint: server.URL, APIKey: "test-key", Model: "test"}
	answer, err := p.Ask(context.Background(), "help", &Page{URL: "https://go.dev"})
	if err != nil || len(answer.Actions) != 1 {
		t.Fatalf("%+v %v", answer, err)
	}
	answer, err = p.Ask(context.Background(), "help", nil)
	if err != nil || len(answer.Actions) != 0 {
		t.Fatal("actions allowed without a shared page")
	}
}
func TestLocalServerAuthorization(t *testing.T) {
	s, err := Start(fstest.MapFS{"index.html": {Data: []byte("assistant")}}, Provider{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	for _, tc := range []struct {
		token, origin string
		code          int
	}{{"", "", 403}, {s.token, "https://evil.example", 403}, {s.token, s.origin, 503}} {
		req, _ := http.NewRequest("GET", s.origin+"/api/state", nil)
		req.Header.Set("Authorization", "Bearer "+tc.token)
		req.Header.Set("Origin", tc.origin)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != tc.code {
			t.Fatalf("got %d want %d", res.StatusCode, tc.code)
		}
	}
	req, _ := http.NewRequest("GET", s.origin+"/", nil)
	req.Host = "evil.example"
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != 403 {
		t.Fatal("DNS rebinding host accepted")
	}
}
func TestExtensionScopeAndPaths(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "content.js"), []byte("document.title='test'"), 0600)
	for _, tc := range []struct {
		script, origin string
		valid          bool
	}{{"content.js", "https://go.dev", true}, {"../escape.js", "https://go.dev", false}, {"content.js", "*", false}, {"content.js", "https://go.dev/path", false}} {
		raw, _ := json.Marshal(map[string]any{"name": "test", "script": tc.script, "origins": []string{tc.origin}})
		os.WriteFile(filepath.Join(dir, "golet-extension.json"), raw, 0600)
		_, err := LoadExtension(dir)
		if (err == nil) != tc.valid {
			t.Fatalf("%+v: %v", tc, err)
		}
	}
}
