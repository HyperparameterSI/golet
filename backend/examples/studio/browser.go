// Author: HyperparameterSI
package main

import (
	"context"
	"errors"
	"github.com/HyperparameterSI/golet/backend/plugins/chromium"
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
}
type BrowserStatus struct {
	Installed bool           `json:"installed"`
	Supported bool           `json:"supported"`
	Running   bool           `json:"running"`
	Host      string         `json:"host"`
	State     chromium.State `json:"state"`
}

func (b *BrowserDemo) Status() BrowserStatus {
	b.mu.Lock()
	defer b.mu.Unlock()
	host, _ := chromium.DefaultHostPath()
	info, err := os.Stat(host)
	status := BrowserStatus{Host: host, Installed: err == nil && !info.IsDir(), Supported: runtime.GOOS == "darwin" || runtime.GOOS == "windows"}
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
func (b *BrowserDemo) Start(style string, target string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.browser != nil {
		select {
		case <-b.browser.Done():
		default:
			return errors.New("请先关闭当前 Browser，再切换风格")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	browser, err := chromium.Launch(ctx, chromium.Options{Style: chromium.BrowserStyle(style), URL: target, EphemeralDemo: true})
	if err != nil {
		cancel()
		return err
	}
	b.browser = browser
	b.cancel = cancel
	return nil
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
