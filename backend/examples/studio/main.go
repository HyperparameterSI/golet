// Author: HyperparameterSI

package main

import (
	"log"
	"runtime"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/examples/internal/gallerydemo"
)

// Greeter is callable from the frontend: `golet generate` turns its methods
// into typed TypeScript functions in src/golet.ts.
type Greeter struct{}

// Greet returns a greeting for name.
func (Greeter) Greet(name string) string {
	if name == "" {
		name = "stranger"
	}
	return "Hello, " + name + "! This message comes from Go."
}

// SystemInfo describes the machine the app runs on.
type SystemInfo struct {
	OS               string `json:"os"`
	Arch             string `json:"arch"`
	GoVersion        string `json:"goVersion"`
	FrameworkVersion string `json:"frameworkVersion"`
}

// Info returns information about the system.
func (Greeter) Info() SystemInfo {
	return SystemInfo{OS: runtime.GOOS, Arch: runtime.GOARCH, GoVersion: runtime.Version(), FrameworkVersion: golet.Version}
}

func main() {
	gallerydemo.ConfigureLanguages()
	browser := &BrowserDemo{}
	defer browser.shutdown()
	golet.Bind(Greeter{})
	golet.Bind(&Showcase{})
	golet.Bind(browser)

	golet.App.WhenReady(func() {
		golet.NewWindow(golet.WindowOptions{
			Title:           "Golet Studio",
			Width:           1280,
			Height:          800,
			MinWidth:        760,
			MinHeight:       360,
			BackgroundColor: "light-dark(#fafafa, #0a0a0a)", // the page's --bg in src/style.css
			// Opens where the user left it last time.
			StateKey: "main",
			// The frontend: devUrl during `golet dev`, frontendDist
			// embedded by `golet build` otherwise (see golet.json).
			URL: "/",
		})
	})
	if err := golet.App.Run(); err != nil {
		log.Fatal(err)
	}
}
