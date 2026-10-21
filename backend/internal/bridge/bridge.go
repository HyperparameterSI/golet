// Author: HyperparameterSI

// Package bridge embeds the renderer runtime that Golet injects into pages.
//
// bridge.js is built from ../../packages/bridge by `bun run build` and
// committed, so building a Golet application never requires Bun.
package bridge

import (
	_ "embed"
	"encoding/json"

	"github.com/HyperparameterSI/golet/backend/internal/platform"
)

//go:embed bridge.js
var source string

// Config is passed to the bridge as __GOLET_CONFIG__.
type Config struct {
	Platform string `json:"platform"`
	WindowID int    `json:"windowId"`
	Version  string `json:"version"`
	// Secret prefixes every message the bridge posts.
	Secret string `json:"secret"`
	// Locale is the language the app shows, e.g. "zh-Hans", and Dir the
	// direction it is written in, "ltr" or "rtl", which pages read as
	// golet.locale and golet.dir from document start.
	Locale string `json:"locale"`
	Dir    string `json:"dir"`
}

// Script returns the bridge source configured for one window.
func Script(cfg Config) string {
	b, _ := json.Marshal(cfg)
	return "(function(__GOLET_CONFIG__){\n" + source + "})(" + string(b) + ");"
}

// TitleBarEvent tells a page of a window with a hidden title bar the room
// the window controls take. The bridge keeps it in the --golet-titlebar-*
// CSS variables.
const TitleBarEvent = "golet:title-bar"

// LocaleEvent tells the pages of a window that the language the app shows
// changed, with a Locale. The bridge configuration tells pages the
// language at document start; the event is sent as a page loads too, for
// one that started loading as the language changed, and the bridge tells
// its listeners only of a language the page did not have.
const LocaleEvent = "golet:locale"

// Locale is the payload of LocaleEvent.
type Locale struct {
	Locale string `json:"locale"`
	Dir    string `json:"dir"`
}

// TitleBar is the payload of TitleBarEvent, in CSS pixels.
type TitleBar struct {
	Height float64 `json:"height"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
}

// NewTitleBar converts the room of the window controls to the CSS pixels
// of a page shown at zoom.
func NewTitleBar(tb platform.TitleBar, zoom float64) TitleBar {
	if zoom <= 0 {
		zoom = 1
	}
	return TitleBar{Height: float64(tb.Height) / zoom, Left: float64(tb.Left) / zoom, Right: float64(tb.Right) / zoom}
}

// TitleBarScript delivers TitleBarEvent at document start, after the
// bridge, so that a page lays out around the window controls before it
// first paints.
func TitleBarScript(tb platform.TitleBar, zoom float64) string {
	p, _ := json.Marshal(NewTitleBar(tb, zoom))
	return `window.__golet&&window.__golet.receive({"t":"event","n":"` + TitleBarEvent + `","p":` + string(p) + `});`
}

// PreferencesEvent tells a page the accent color and the accessibility
// settings of the desktop, as golet.ThemePreferences encodes them. The
// bridge keeps the latest for golet-runtime's preferences(), the accent in
// the --golet-accent CSS variable and the text scale in --golet-text-scale.
const PreferencesEvent = "golet:preferences"

// PreferencesScript delivers PreferencesEvent with the JSON prefs at
// document start, after the bridge has registered its handler.
func PreferencesScript(prefs []byte) string {
	return `window.__golet&&window.__golet.receive({"t":"event","n":"` + PreferencesEvent + `","p":` + string(prefs) + `});`
}
