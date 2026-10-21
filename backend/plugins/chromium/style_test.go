// Author: HyperparameterSI
package chromium

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestBrowserStyle(t *testing.T) {
	for _, tc := range []struct {
		name    string
		opts    Options
		args    []string
		invalid bool
	}{
		{name: "default"},
		{name: "app", opts: Options{AppMode: true}, args: []string{"--golet-app-mode"}},
		{name: "app chrome", opts: Options{AppMode: true, Style: StyleChrome}, invalid: true},
		{name: "app legacy chrome", opts: Options{AppMode: true, ChromeStyle: true}, invalid: true},
		{name: "app extensions", opts: Options{AppMode: true, Extensions: []string{"unused"}}, invalid: true},
		{name: "app assistant", opts: Options{AppMode: true, AssistantURL: "http://localhost"}, invalid: true},
		{name: "alloy", opts: Options{Style: StyleAlloy}},
		{name: "chrome", opts: Options{Style: StyleChrome}, args: []string{"--golet-chrome-style"}},
		{name: "legacy", opts: Options{ChromeStyle: true}, args: []string{"--golet-chrome-style"}},
		{name: "consistent", opts: Options{Style: StyleChrome, ChromeStyle: true}, args: []string{"--golet-chrome-style"}},
		{name: "conflict", opts: Options{Style: StyleAlloy, ChromeStyle: true}, invalid: true},
		{name: "unknown", opts: Options{Style: "unknown"}, invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args, err := extensionArguments(tc.opts)
			if (err != nil) != tc.invalid {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(args, tc.args) {
				t.Fatalf("got %v want %v", args, tc.args)
			}
		})
	}
}

func TestInvalidStyleDoesNotCreateProfile(t *testing.T) {
	profile := t.TempDir() + "/profile"
	_, err := Launch(context.Background(), Options{Style: "invalid", ProfileDir: profile})
	if err == nil {
		t.Fatal("invalid style accepted")
	}
	if _, statErr := os.Stat(profile); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("profile created: %v", statErr)
	}
	// Style validation precedes host lookup and profile creation.
	if err.Error() != `chromium: unsupported browser style "invalid"; use StyleAlloy or StyleChrome` {
		t.Fatal(err)
	}
}
