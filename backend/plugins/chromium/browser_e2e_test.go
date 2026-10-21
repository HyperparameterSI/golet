// Author: HyperparameterSI

package chromium

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

// This test opens a real CEF window. It waits on state events, never sleeps to
// infer success. The one-minute context is solely a hang guard.
func TestBrowserE2E(t *testing.T) {
	if os.Getenv("GOLET_CHROMIUM_E2E") != "1" || runtime.GOOS != "darwin" {
		t.Skip("requires an installed CEF host and macOS desktop")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<!doctype html><title>Golet %s</title><h1>Golet Browser</h1><p>%s</p>", r.URL.Path, r.URL.Path)
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	browser, err := Launch(ctx, Options{URL: server.URL + "/first", EphemeralDemo: true})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		closeCtx, stop := context.WithTimeout(context.Background(), time.Minute)
		defer stop()
		if err := browser.Close(closeCtx); err != nil {
			t.Error(err)
		}
	}()
	wait := func(path string) {
		t.Helper()
		for {
			state := browser.State()
			if !state.Loading && state.URL == server.URL+path && state.Title == "Golet "+path {
				return
			}
			select {
			case _, ok := <-browser.Events():
				if !ok {
					t.Fatalf("browser exited: %v", browser.Wait())
				}
			case <-ctx.Done():
				t.Fatalf("waiting for %s: state=%+v", path, browser.State())
			}
		}
	}
	wait("/first")
	firstID := browser.State().BrowserID
	var title string
	if err := browser.Evaluate(ctx, firstID, "document.title", &title); err != nil || title != "Golet /first" {
		t.Fatalf("evaluate: %q %v", title, err)
	}
	if err := browser.Evaluate(ctx, firstID, "throw new Error('expected')", nil); err == nil {
		t.Fatal("script error was ignored")
	}
	if _, err := browser.DevTools(ctx, firstID, "Golet.invalidMethod", nil); err == nil {
		t.Fatal("invalid CDP method accepted")
	}
	if err := browser.Evaluate(ctx, firstID, "document.body.insertAdjacentHTML('beforeend', '<button id=action>Click</button>'); document.querySelector('#action').onclick=()=>document.querySelector('#action').textContent='Clicked'; true", nil); err != nil {
		t.Fatal(err)
	}
	if err := browser.Evaluate(ctx, firstID, "document.querySelector('#action').click();document.querySelector('#action').textContent", &title); err != nil || title != "Clicked" {
		t.Fatalf("DOM action: %q %v", title, err)
	}
	if err := browser.Navigate(ctx, server.URL+"/second"); err != nil {
		t.Fatal(err)
	}
	wait("/second")
	if !browser.State().CanGoBack {
		t.Fatal("missing back history")
	}
	if err := browser.Back(ctx); err != nil {
		t.Fatal(err)
	}
	wait("/first")
	if !browser.State().CanGoForward {
		t.Fatal("missing forward history")
	}
	if err := browser.Forward(ctx); err != nil {
		t.Fatal(err)
	}
	wait("/second")
	if err := browser.NewTab(ctx, server.URL+"/third"); err != nil {
		t.Fatal(err)
	}
	wait("/third")
	thirdID := browser.State().BrowserID
	if thirdID == firstID || len(browser.State().Tabs) != 2 {
		t.Fatalf("new tab state: %+v", browser.State())
	}
	if err := browser.SelectTab(ctx, firstID); err != nil {
		t.Fatal(err)
	}
	wait("/second")
	if err := browser.CloseTab(ctx, thirdID); err != nil {
		t.Fatal(err)
	}
	for len(browser.State().Tabs) != 1 {
		select {
		case _, ok := <-browser.Events():
			if !ok {
				t.Fatal("closing a background tab closed the host")
			}
		case <-ctx.Done():
			t.Fatal("tab did not close")
		}
	}
	if err := browser.SelectTab(ctx, thirdID); err == nil {
		t.Fatal("closed tab still selectable")
	}
	if err := browser.NewTab(ctx, ""); err != nil {
		t.Fatal(err)
	}
	for browser.State().Title != "新标签页" || browser.State().Loading {
		select {
		case _, ok := <-browser.Events():
			if !ok {
				t.Fatal("new-tab page closed the host")
			}
		case <-ctx.Done():
			t.Fatal("new-tab page did not load")
		}
	}
	if err := browser.CloseTab(ctx, browser.State().BrowserID); err != nil {
		t.Fatal(err)
	}
	wait("/second")
	if err := browser.Reload(ctx); err != nil {
		t.Fatal(err)
	}
	if err := browser.Stop(ctx); err != nil {
		t.Fatal(err)
	}
}

// Chrome extensions execute only in the explicit Chrome-style compatibility mode.
func TestChromeExtensionE2E(t *testing.T) {
	if os.Getenv("GOLET_CHROMIUM_E2E") != "1" || runtime.GOOS != "darwin" {
		t.Skip("requires real CEF host")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<!doctype html><title>Extension test</title><h1>Fixture</h1>")
	}))
	defer server.Close()
	ext := t.TempDir()
	if err := os.WriteFile(ext+"/manifest.json", []byte(`{"manifest_version":3,"name":"Golet test extension","version":"1.0","content_scripts":[{"matches":["http://127.0.0.1/*"],"js":["content.js"],"run_at":"document_idle"}]}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ext+"/content.js", []byte(`document.title="Golet extension active"`), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	b, err := Launch(ctx, Options{URL: server.URL, ProfileDir: t.TempDir(), Style: StyleChrome, Extensions: []string{ext}})
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
	for b.State().Title != "Golet extension active" {
		select {
		case _, ok := <-b.Events():
			if !ok {
				t.Fatalf("host closed: %v", b.Wait())
			}
		case <-ctx.Done():
			t.Fatalf("extension did not run: %+v", b.State())
		}
	}
}

func TestEphemeralDemoE2E(t *testing.T) {
	if os.Getenv("GOLET_CHROMIUM_E2E") != "1" || runtime.GOOS != "darwin" {
		t.Skip("requires real CEF host")
	}
	for _, style := range []BrowserStyle{StyleAlloy, StyleChrome} {
		t.Run(string(style), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			b, err := Launch(ctx, Options{Style: style, EphemeralDemo: true})
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := b.Close(ctx); err != nil {
					t.Error(err)
				}
			}()
			var profile string
			for _, arg := range b.cmd.Args {
				if strings.HasPrefix(arg, "--cache-path=") {
					profile = strings.TrimPrefix(arg, "--cache-path=")
				}
			}
			if profile == "" {
				t.Fatal("missing profile")
			}
			if err := b.Close(ctx); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(profile); !os.IsNotExist(err) {
				t.Fatalf("temporary profile not removed: %v", err)
			}
		})
	}
}

func TestAppModeE2E(t *testing.T) {
	if os.Getenv("GOLET_CHROMIUM_E2E") != "1" || runtime.GOOS != "darwin" {
		t.Skip("requires an installed CEF host and macOS desktop")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<!doctype html><title>Mini app</title><h1>Application content</h1>")
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	b, err := Launch(ctx, Options{AppMode: true, URL: server.URL, EphemeralDemo: true})
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close(context.Background())
	for b.State().Loading || b.State().Title != "Mini app" {
		select {
		case _, ok := <-b.Events():
			if !ok {
				t.Fatalf("host exited: %v", b.Wait())
			}
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		}
	}
	var height int
	if err := b.Evaluate(ctx, b.State().BrowserID, "innerHeight", &height); err != nil || height != 820 {
		t.Fatalf("content should fill the 820px client area, got %d: %v", height, err)
	}
	if err := b.NewTab(ctx, server.URL); err == nil || !strings.Contains(err.Error(), "AppMode") {
		t.Fatalf("new tab should be rejected: %v", err)
	}
	var blocked bool
	if err := b.Evaluate(ctx, b.State().BrowserID, "window.open(location.href) === null", &blocked); err != nil || !blocked {
		t.Fatalf("popup should be rejected: %v %v", blocked, err)
	}
	if len(b.State().Tabs) != 1 {
		t.Fatal("app must have one page")
	}
}

func TestProductBrandingE2E(t *testing.T) {
	if os.Getenv("GOLET_CHROMIUM_E2E") != "1" || runtime.GOOS != "darwin" {
		t.Skip("requires real CEF host")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	name := "feo & <测试>"
	b, err := Launch(ctx, Options{AppName: name, EphemeralDemo: true})
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close(context.Background())
	for b.State().Loading || b.State().Title != "新标签页" {
		select {
		case _, ok := <-b.Events():
			if !ok {
				t.Fatal(b.Wait())
			}
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		}
	}
	var got string
	if err := b.Evaluate(ctx, b.State().BrowserID, `document.querySelector('.brand').innerText.trim()`, &got); err != nil || got != name {
		t.Fatalf("brand=%q: %v", got, err)
	}
	var valid bool
	if err := b.Evaluate(ctx, b.State().BrowserID, `document.querySelector('.brand img').alt === new URLSearchParams(location.search).get('appName') && document.querySelector('.brand strong').children.length === 0`, &valid); err != nil || !valid {
		t.Fatalf("unsafe branding: %v %v", valid, err)
	}
}
