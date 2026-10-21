// Author: HyperparameterSI

package chromium

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"testing"
	"time"
)

// TestMain provides a tiny native-host stand-in for protocol/lifecycle tests;
// it does not simulate or claim to test Chromium rendering.
func TestMain(m *testing.M) {
	if os.Getenv("GOLET_CHROMIUM_TEST_HOST") == "1" {
		out := json.NewEncoder(os.Stdout)
		_ = out.Encode(message{Event: "ready", State: State{BrowserID: 1, URL: "about:blank"}})
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			var request message
			if json.Unmarshal(scanner.Bytes(), &request) != nil {
				continue
			}
			if request.Method == "navigate" {
				_ = out.Encode(message{Event: "state", State: State{BrowserID: 1, URL: request.URL}})
			}
			_ = out.Encode(message{ID: request.ID})
			if request.Method == "close" {
				os.Exit(0)
			}
		}
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestBrowserProtocolLifecycle(t *testing.T) {
	t.Setenv("GOLET_CHROMIUM_TEST_HOST", "1")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	browser, err := Launch(ctx, Options{HostPath: executable, ProfileDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	defer browser.Close(ctx)
	if err := browser.Navigate(ctx, "https://go.dev/"); err != nil {
		t.Fatal(err)
	}
	if state := browser.State(); state.URL != "https://go.dev/" {
		t.Fatalf("state=%+v", state)
	}
	var group sync.WaitGroup
	for range 20 {
		group.Go(func() {
			if err := browser.Reload(ctx); err != nil {
				t.Error(err)
			}
		})
	}
	group.Wait()
	if err := browser.Close(ctx); err != nil {
		t.Fatal(err)
	}
	if err := browser.Close(ctx); err != nil {
		t.Fatal(err)
	}
	if err := browser.Reload(ctx); !errors.Is(err, ErrClosed) {
		t.Fatalf("closed reload: %v", err)
	}
	browser.mu.Lock()
	defer browser.mu.Unlock()
	if len(browser.pending) != 0 {
		t.Fatalf("leaked pending requests: %d", len(browser.pending))
	}
}

func TestBrowserContextClosesHost(t *testing.T) {
	t.Setenv("GOLET_CHROMIUM_TEST_HOST", "1")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	browser, err := Launch(ctx, Options{HostPath: executable, ProfileDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	guard, stop := context.WithTimeout(context.Background(), time.Minute)
	defer stop()
	select {
	case <-browser.Done():
	case <-guard.Done():
		t.Fatal("cancelled browser did not exit")
	}
	if err := browser.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestCloseCancelledTerminatesHost(t *testing.T) {
	t.Setenv("GOLET_CHROMIUM_TEST_HOST", "1")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	lifetime, stop := context.WithTimeout(context.Background(), time.Minute)
	defer stop()
	browser, err := Launch(lifetime, Options{HostPath: executable, ProfileDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	closeCtx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := browser.Close(closeCtx); err == nil {
		t.Fatal("forced close returned no error")
	}
	select {
	case <-browser.Done():
	case <-lifetime.Done():
		t.Fatal("forced host did not exit")
	}
}
