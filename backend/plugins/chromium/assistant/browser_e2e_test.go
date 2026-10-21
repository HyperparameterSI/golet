// Author: HyperparameterSI
package assistant

import (
	"context"
	"fmt"
	"github.com/HyperparameterSI/golet/backend/plugins/chromium"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestAssistantActionsE2E(t *testing.T) {
	if os.Getenv("GOLET_CHROMIUM_E2E") != "1" || runtime.GOOS != "darwin" {
		t.Skip("requires native CEF")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!doctype html><title>Assistant fixture</title><p>Article text</p><input aria-label="Search"><input type="password" value="secret"><button onclick="this.textContent='Done'">Submit</button>`)
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	b, err := chromium.Launch(ctx, chromium.Options{URL: server.URL, ProfileDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		c, stop := context.WithTimeout(context.Background(), time.Minute)
		defer stop()
		if err := b.Close(c); err != nil {
			t.Error(err)
		}
	}()
	for b.State().Title != "Assistant fixture" || b.State().Loading {
		select {
		case <-b.Events():
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		}
	}
	id := b.State().BrowserID
	page, err := Snapshot(ctx, b, id)
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Elements) != 2 {
		t.Fatalf("expected safe controls: %+v", page)
	}
	if err := Apply(ctx, b, id, page, Action{Kind: "fill", Element: 0, Value: "Golet"}); err != nil {
		t.Fatal(err)
	}
	var input string
	if err := b.Evaluate(ctx, id, "document.querySelector('input').value", &input); err != nil || input != "Golet" {
		t.Fatalf("fill %q %v", input, err)
	}
	if err := Apply(ctx, b, id, page, Action{Kind: "click", Element: 1}); err != nil {
		t.Fatal(err)
	}
	if err := Apply(ctx, b, id, page, Action{Kind: "click", Element: 1}); err == nil {
		t.Fatal("stale element accepted")
	}
	ext := Extension{Name: "fixture", Origins: []string{server.URL}, Script: "document.title='Extension done';return true;"}
	if err := ext.Run(ctx, b, id); err != nil {
		t.Fatal(err)
	}
	ext.Origins = []string{"https://go.dev"}
	if err := ext.Run(ctx, b, id); err == nil {
		t.Fatal("ungranted origin accepted")
	}
}
