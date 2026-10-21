// Author: HyperparameterSI

package chromium

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func TestNavigationURL(t *testing.T) {
	for _, raw := range []string{"", "about:blank", "https://go.dev", "http://127.0.0.1:3000/path?q=1"} {
		if _, err := navigationURL(raw); err != nil {
			t.Errorf("%q: %v", raw, err)
		}
	}
	for _, raw := range []string{"javascript:alert(1)", "file:///etc/passwd", "https:///path", "https://user:secret@example.com", "data:text/html,x", "--no-sandbox", "https://go.dev/" + strings.Repeat("a", 16384)} {
		if _, err := navigationURL(raw); err == nil {
			t.Errorf("accepted %q", raw)
		}
	}
}

func TestLaunchMissingHost(t *testing.T) {
	_, err := Launch(context.Background(), Options{HostPath: filepath.Join(t.TempDir(), "missing")})
	if err == nil || !strings.Contains(err.Error(), "not installed") {
		t.Fatalf("got %v", err)
	}
}

func TestLaunchCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Launch(ctx, Options{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
}

// A stalled host must not turn a caller's cancelled context into an unbounded
// pipe write. The bounded writer queue is intentionally full here.
func TestCancelledCallWithFullQueue(t *testing.T) {
	b := &Browser{writes: make(chan []byte, 1), pending: make(map[uint64]chan message), done: make(chan struct{})}
	b.writes <- []byte("occupied")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := b.Reload(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
	if len(b.pending) != 0 {
		t.Fatalf("pending requests leaked: %d", len(b.pending))
	}
}
