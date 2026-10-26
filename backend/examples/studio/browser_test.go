// Author: HyperparameterSI
package main

import "testing"

func TestBrowserDemoRequiresLaunch(t *testing.T) {
	demo := &BrowserDemo{}
	if err := demo.Command("navigate", "https://go.dev"); err == nil {
		t.Fatal("command accepted without host")
	}
	if err := demo.Start("invalid-style", "https://go.dev"); err == nil {
		demo.shutdown()
		t.Fatal("invalid style launched")
	}
	if demo.browser != nil {
		t.Fatal("failed launch retained host")
	}
}

func TestBrowserThemeSettings(t *testing.T) {
	settings := BrowserThemeSettings{Preset: "Violet", Dark: true, Radius: 8, FontSize: 16}
	theme, err := settings.theme()
	if err != nil || !theme.Dark || theme.Radius != 8 || theme.FontSize != 16 {
		t.Fatalf("theme=%+v error=%v", theme, err)
	}
	demo := &BrowserDemo{}
	if err := demo.ApplyTheme(settings); err == nil {
		t.Fatal("theme accepted without browser")
	}
	for _, invalid := range []BrowserThemeSettings{
		{Preset: "unknown", Radius: 8, FontSize: 14},
		{Preset: "Blue", Radius: -1, FontSize: 14},
		{Preset: "Blue", Radius: 8, FontSize: 100},
	} {
		if _, err := invalid.theme(); err == nil {
			t.Fatalf("invalid settings accepted: %+v", invalid)
		}
	}
}
