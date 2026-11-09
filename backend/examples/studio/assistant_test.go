// Author: HyperparameterSI
package main

import (
	"context"
	"github.com/HyperparameterSI/golet/backend/plugins/dsh-go"
	"github.com/HyperparameterSI/golet/backend/plugins/dsh-go/native"
	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
	"github.com/JopenChen/dsh-go/pkg/llm"
	"strings"
	"testing"
)

func TestStudioAssistantOffline(t *testing.T) {
	t.Setenv("DEEPSEEK_API_KEY", "")
	for _, allow := range []bool{true, false} {
		options, notice, err := studioAssistantOptions()
		if err != nil {
			t.Fatal(err)
		}
		if options.Model != "offline-demo" || !strings.Contains(notice, "离线演示") {
			t.Fatal("offline mode not labeled")
		}
		approvals := 0
		options.Approve = func(context.Context, llm.ToolCall) (bool, error) { approvals++; return allow, nil }
		conversation, err := dshgo.New(options)
		if err != nil {
			t.Fatal(err)
		}
		result, err := conversation.Run(context.Background(), "现在几点？", nil)
		conversation.Close()
		if err != nil || approvals != 1 || result.Steps != 2 {
			t.Fatalf("%v approvals=%d steps=%d", err, approvals, result.Steps)
		}
		if !strings.Contains(result.Text, "离线演示") || !strings.Contains(result.Text, "forma.PromptInput") {
			t.Fatal(result.Text)
		}
		if !allow && !strings.Contains(result.Text, "工具未执行") {
			t.Fatal("denial presented as success")
		}
	}
}
func TestStudioAssistantThemeAndControls(t *testing.T) {
	t.Setenv("DEEPSEEK_API_KEY", "")
	for _, dark := range []bool{false, true} {
		options, notice, err := studioAssistantOptions()
		if err != nil {
			t.Fatal(err)
		}
		p, err := native.New(options)
		if err != nil {
			t.Fatal(err)
		}
		defer p.Close()
		p.SetNotice(notice)
		tt := ui.NewTester(func(c *ui.Context) { c.SetTheme(forma.Theme(dark)); p.View(c) }, 1000, 820)
		for _, label := range []string{"Prompt input", "Prompt text", "Model", "Send prompt"} {
			if _, ok := tt.Find(label); !ok {
				t.Fatalf("dark=%v missing %s", dark, label)
			}
		}
		if !tt.HasText(notice) {
			t.Fatal("mode notice missing")
		}
	}
}

func TestAssistantSettingsApply(t *testing.T) {
	t.Setenv("DEEPSEEK_API_KEY", "")
	options, _, err := studioAssistantOptions()
	if err != nil {
		t.Fatal(err)
	}
	p, err := native.New(options)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	settings := assistantSettings{open: true, config: assistantConfig{baseURL: "https://api.deepseek.com", model: "deepseek-chat"}}
	settings.apply(p, false)
	if settings.err == "" || !settings.open {
		t.Fatal("blank key accepted")
	}
	settings.config.key = "test-key-only"
	settings.config.baseURL = "http://example.com"
	settings.apply(p, false)
	if settings.err == "" || !settings.open {
		t.Fatal("insecure endpoint accepted")
	}
	settings.config.baseURL = "https://api.deepseek.com"
	settings.apply(p, false)
	if settings.err != "" || settings.open {
		t.Fatalf("apply failed: %s", settings.err)
	}
	tt := ui.NewTester(func(c *ui.Context) { p.View(c) }, 1000, 820)
	if !tt.HasText("deepseek-chat") {
		t.Fatal("configured model not updated")
	}
	settings.open = true
	settings.apply(p, true)
	if settings.config.key != "" || settings.open {
		t.Fatal("offline mode did not clear key")
	}
	tt.Frame()
	if !tt.HasText("offline-demo") {
		t.Fatal("offline model missing")
	}
}
