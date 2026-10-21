// Author: HyperparameterSI

// Browser launches the optional CEF host and a local assistant workspace.
package main

import (
	"context"
	"flag"
	"github.com/HyperparameterSI/golet/backend/plugins/chromium"
	"github.com/HyperparameterSI/golet/backend/plugins/chromium/assistant"
	"log"
	"os"
	"os/signal"
	"strings"
)

type paths []string

func (p *paths) String() string         { return strings.Join(*p, ",") }
func (p *paths) Set(value string) error { *p = append(*p, value); return nil }
func main() {
	demo := flag.Bool("ephemeral-demo", false, "disposable demo profile without macOS Keychain protection; do not use real credentials")
	target := flag.String("url", "", "initial URL; empty opens the new-tab page")
	host := flag.String("host", "", "native host executable")
	profile := flag.String("profile", "", "dedicated profile directory")
	openAssistant := flag.Bool("assistant", false, "open assistant workspace on startup")
	style := flag.String("style", "", "window style: alloy (default) or chrome")
	chrome := flag.Bool("chrome-style", false, "deprecated alias for -style chrome")
	assets := flag.String("assistant-assets", "frontend/examples/browser/assistant", "assistant frontend directory")
	var userExtensions, chromeExtensions paths
	flag.Var(&userExtensions, "user-extension", "trusted Golet extension directory (repeatable)")
	flag.Var(&chromeExtensions, "extension", "trusted unpacked Manifest V3 directory; requires -style chrome (repeatable)")
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	var extensions []assistant.Extension
	for _, path := range userExtensions {
		ext, err := assistant.LoadExtension(path)
		if err != nil {
			log.Fatal(err)
		}
		extensions = append(extensions, ext)
	}
	if _, err := os.Stat(*assets + "/index.html"); err != nil {
		log.Fatal("assistant assets not found; run from repository root or set -assistant-assets")
	}
	server, err := assistant.Start(os.DirFS(*assets), assistant.Provider{Endpoint: os.Getenv("GOLET_AI_ENDPOINT"), APIKey: os.Getenv("GOLET_AI_API_KEY"), Model: os.Getenv("GOLET_AI_MODEL")}, extensions)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	browser, err := chromium.Launch(ctx, chromium.Options{EphemeralDemo: *demo, HostPath: *host, URL: *target, ProfileDir: *profile, Style: chromium.BrowserStyle(*style), ChromeStyle: *chrome, Extensions: chromeExtensions, AssistantURL: server.URL()})
	if err != nil {
		log.Fatal(err)
	}
	server.Bind(browser)
	if *chrome || *style == string(chromium.StyleChrome) || *openAssistant {
		if err := browser.NewTab(ctx, server.URL()); err != nil {
			log.Print(err)
		}
	}
	log.Print("Golet Browser ready. Use the sparkle button for AI and user extensions.")
	for range browser.Events() {
	}
	if err := browser.Wait(); err != nil && ctx.Err() == nil {
		log.Print(err)
	}
}
