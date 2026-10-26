// Author: HyperparameterSI
package main

import (
	"context"
	"errors"
	"github.com/HyperparameterSI/golet/backend/plugins/chromium"
	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
	formabrowser "github.com/HyperparameterSI/golet/backend/ui/forma/browser"
	"os"
	"runtime"
	"sync"
	"time"
)

// BrowserDemo controls one opt-in CEF host owned by Studio.
type BrowserDemo struct {
	mu      sync.Mutex
	browser *chromium.Browser
	cancel  context.CancelFunc
	style   chromium.BrowserStyle
}
type BrowserStatus struct {
	ThemeSupported bool           `json:"themeSupported"`
	Style          string         `json:"style"`
	Presets        []forma.Preset `json:"presets"`
	Installed      bool           `json:"installed"`
	Supported      bool           `json:"supported"`
	Running        bool           `json:"running"`
	Host           string         `json:"host"`
	State          chromium.State `json:"state"`
}

func (b *BrowserDemo) Status() BrowserStatus {
	b.mu.Lock()
	defer b.mu.Unlock()
	host, _ := chromium.DefaultHostPath()
	info, err := os.Stat(host)
	status := BrowserStatus{ThemeSupported: runtime.GOOS == "darwin", Style: string(b.style), Presets: forma.Presets(), Host: host, Installed: err == nil && !info.IsDir(), Supported: runtime.GOOS == "darwin" || runtime.GOOS == "windows"}
	if b.browser != nil {
		select {
		case <-b.browser.Done():
		default:
			status.Running = true
			status.State = b.browser.State()
		}
	}
	return status
}

// BrowserThemeSettings configures the native Forma browser shell.
type BrowserThemeSettings struct {
	Preset   string  `json:"preset"`
	Dark     bool    `json:"dark"`
	Radius   float32 `json:"radius"`
	FontSize float32 `json:"fontSize"`
}

func (settings BrowserThemeSettings) theme() (*ui.Theme, error) {
	valid := false
	for _, preset := range forma.Presets() {
		if string(preset) == settings.Preset {
			valid = true
		}
	}
	if !valid || !(settings.Radius >= 0 && settings.Radius <= 21) || !(settings.FontSize >= 10 && settings.FontSize <= 20) {
		return nil, errors.New("请选择有效的 Forma 主题，圆角 0–21，字号 10–20")
	}
	return forma.NewTheme(forma.ThemeOptions{Preset: forma.Preset(settings.Preset), Dark: settings.Dark, Radius: &settings.Radius, FontSize: settings.FontSize}), nil
}

func (b *BrowserDemo) Start(style string, target string) error {
	return b.StartWithTheme(style, target, BrowserThemeSettings{Preset: "Blue", Radius: 16, FontSize: 14})
}

// StartWithTheme applies Forma styling where the native shell supports it.
func (b *BrowserDemo) StartWithTheme(style string, target string, settings BrowserThemeSettings) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.browser != nil {
		select {
		case <-b.browser.Done():
		default:
			return errors.New("请先关闭当前 Browser，再切换风格")
		}
	}
	theme, err := settings.theme()
	if err != nil {
		return err
	}
	opts := chromium.Options{Style: chromium.BrowserStyle(style), URL: target, EphemeralDemo: true}
	ctx, cancel := context.WithCancel(context.Background())
	var browser *chromium.Browser
	themed := runtime.GOOS == "darwin" && (style == "" || style == "alloy")
	if themed {
		browser, err = formabrowser.Launch(ctx, opts, theme)
		if err == nil {
			// Verify the installed host understands themes; old hosts ignore startup switches.
			check, stop := context.WithTimeout(ctx, 5*time.Second)
			err = formabrowser.SetTheme(check, browser, theme)
			stop()
			if err != nil {
				err = errors.New("浏览器宿主未能应用主题，请重新运行 pnpm setup:chromium：" + err.Error())
			}
		}
	} else {
		browser, err = chromium.Launch(ctx, opts)
	}
	if err != nil {
		cancel()
		return err
	}
	b.browser = browser
	b.cancel = cancel
	b.style = opts.Style
	if b.style == "" {
		b.style = chromium.StyleAlloy
	}
	return nil
}

// ApplyTheme updates the running shell without closing its tabs.
func (b *BrowserDemo) ApplyTheme(settings BrowserThemeSettings) error {
	theme, err := settings.theme()
	if err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.browser == nil {
		return errors.New("请先启动 Browser")
	}
	if runtime.GOOS != "darwin" || b.style != chromium.StyleAlloy {
		return errors.New("Forma Browser 主题当前仅支持 macOS Alloy")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return formabrowser.SetTheme(ctx, b.browser, theme)
}

func (b *BrowserDemo) Command(action string, target string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.browser == nil {
		return errors.New("请先启动 Browser")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	switch action {
	case "navigate":
		return b.browser.Navigate(ctx, target)
	case "back":
		return b.browser.Back(ctx)
	case "forward":
		return b.browser.Forward(ctx)
	case "reload":
		return b.browser.Reload(ctx)
	case "newtab":
		return b.browser.NewTab(ctx, target)
	case "stop":
		return b.browser.Stop(ctx)
	case "close":
		err := b.browser.Close(ctx)
		if b.cancel != nil {
			b.cancel()
		}
		b.browser = nil
		return err
	default:
		return errors.New("未知浏览器命令")
	}
}
func (b *BrowserDemo) shutdown() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.cancel != nil {
		b.cancel()
	}
	if b.browser != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = b.browser.Close(ctx)
	}
}
