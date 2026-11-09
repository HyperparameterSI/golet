// Author: HyperparameterSI
package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/examples/internal/demotheme"
	"github.com/HyperparameterSI/golet/backend/plugins/dsh-go"
	"github.com/HyperparameterSI/golet/backend/plugins/dsh-go/native"
	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
	"github.com/herlanai/omnix/pkg/llm"
)

func studioAssistantOptions() (dshgo.Options, string, error) {
	return assistantConfigFromEnv().options()
}

type assistantConfig struct{ key, baseURL, model string }

func assistantConfigFromEnv() assistantConfig {
	c := assistantConfig{key: os.Getenv("DEEPSEEK_API_KEY"), baseURL: os.Getenv("DEEPSEEK_BASE_URL"), model: os.Getenv("DEEPSEEK_MODEL")}
	if c.baseURL == "" {
		c.baseURL = "https://api.deepseek.com"
	}
	if c.model == "" {
		c.model = "deepseek-chat"
	}
	return c
}
func (config assistantConfig) options() (dshgo.Options, string, error) {
	var adapter llm.LLMAdapter = studioOfflineAssistant{}
	model := "offline-demo"
	notice := "离线演示 · 不调用模型、不联网。点击右上角「模型设置」填写 DeepSeek API Key 即可切换在线模式。"
	if key := strings.TrimSpace(config.key); key != "" {
		var err error
		adapter, err = dshgo.NewDeepSeek(dshgo.DeepSeekOptions{APIKey: key, BaseURL: strings.TrimSpace(config.baseURL)})
		if err != nil {
			return dshgo.Options{}, "", err
		}
		model = strings.TrimSpace(config.model)
		if model == "" {
			model = "deepseek-chat"
		}
		notice = "DeepSeek 在线模式 · 当前模型：" + model + " · 密钥仅保存在 Go 后端。附件和思考强度尚未接入此演示。"
	}
	return dshgo.Options{Adapter: adapter, Model: model, System: "你是 Golet Studio 的原生 AI 助手。使用 Markdown 和 Go 代码回答。需要当前时间时调用 current_time。", Tools: []dshgo.Tool{{Name: "current_time", Description: "读取本机当前时间", RequiresApproval: true, Execute: func(ctx context.Context, _ map[string]any) (any, error) {
		return time.Now().Format(time.RFC3339), ctx.Err()
	}}}}, notice, nil
}

func openAssistant() (*golet.Window, error) {
	options, notice, err := studioAssistantOptions()
	if err != nil {
		return nil, err
	}
	panel, err := native.New(options)
	if err != nil {
		return nil, err
	}
	panel.SetNotice(notice)
	settings := &assistantSettings{config: assistantConfigFromEnv()}
	panel.SetHeaderActions(func(c *ui.Context) {
		if forma.Button(c, "模型设置", forma.ButtonOptions{Variant: forma.Ghost}).Clicked() {
			settings.open = true
			settings.err = ""
		}
	})
	window := golet.NewWindow(golet.WindowOptions{Title: "Studio · DSH-Go AI Assistant", Width: 1000, Height: 820, MinWidth: 680, MinHeight: 600, TitleBarStyle: golet.TitleBarHiddenInset, TrafficLightPosition: &golet.Point{X: 12, Y: 11}, Content: ui.View(func(c *ui.Context) { demotheme.Apply(c); panel.View(c); settings.view(c, panel) })})
	panel.SetWindow(window)
	return window, nil
}

// Settings stay in this window's Go memory; no IPC, disk, or log contains keys.
type assistantSettings struct {
	config assistantConfig
	open   bool
	err    string
}

func (s *assistantSettings) apply(panel *native.Panel, offline bool) {
	c := s.config
	if offline {
		c.key = ""
	} else if strings.TrimSpace(c.key) == "" {
		s.err = "请输入 DeepSeek API Key"
		return
	}
	options, notice, err := c.options()
	if err == nil {
		err = panel.Configure(options)
	}
	if err != nil {
		s.err = err.Error()
		return
	}
	panel.SetNotice(notice)
	s.config = c
	s.err = ""
	s.open = false
}
func (s *assistantSettings) view(c *ui.Context, panel *native.Panel) {
	forma.Dialog(c, &s.open, "DeepSeek 模型设置", func() {
		forma.Field(c, "API Key", func() {
			forma.Input(c, &s.config.key, forma.InputOptions{Placeholder: "sk-…"}).Password().Label("DeepSeek API Key")
		})
		forma.Field(c, "API 地址", func() { forma.Input(c, &s.config.baseURL).Label("DeepSeek API URL") })
		forma.Field(c, "模型", func() {
			forma.Input(c, &s.config.model, forma.InputOptions{Placeholder: "deepseek-chat"}).Label("DeepSeek model")
		})
		ui.Text(c, "密钥仅保存在当前窗口内存中，关闭后不保存。应用设置会清空当前会话；配置有效性将在首次发送时由服务端验证。").FontSize(12).TextColor(c.Theme().TextMuted)
		if panel.Busy() {
			ui.Text(c, "请先停止当前回复，再应用设置。").TextColor(c.Theme().TextMuted)
		}
		if s.err != "" {
			ui.Text(c, s.err).TextColor(c.Theme().Danger)
		}
		ui.Row(c).Gap(8).Children(func() {
			if forma.Button(c, "离线演示", forma.ButtonOptions{Variant: forma.Outline, Disabled: panel.Busy()}).Clicked() {
				s.apply(panel, true)
			}
			if forma.Button(c, "应用并新建会话", forma.ButtonOptions{Disabled: panel.Busy()}).Clicked() {
				s.apply(panel, false)
			}
		})
	})
}

type studioOfflineAssistant struct{}

func (studioOfflineAssistant) Name() string { return "Studio offline" }
func (studioOfflineAssistant) Chat(ctx context.Context, r llm.ChatRequest, emit func(llm.StreamChunk)) (llm.Usage, error) {
	if err := ctx.Err(); err != nil {
		return llm.Usage{}, err
	}
	if len(r.Messages) == 0 {
		return llm.Usage{}, nil
	}
	last := r.Messages[len(r.Messages)-1]
	if len(last.Content) > 0 && last.Content[0].ToolResult != nil {
		result := last.Content[0].ToolResult
		text := "工具返回：" + result.Content
		if result.IsError {
			text = "工具未执行或返回错误：" + result.Content
		}
		for _, part := range []string{"**这是离线演示，未调用真实模型。**\n\n", text, "\n\n以下是原生控件的使用示例：\n\n```go\nforma.AIMessage(c, forma.MessageAssistant, reply)\nforma.PromptInput(c, &draft, options)\n```\n\n你可以复制回复、评价、编辑问题，或清空会话。"} {
			if err := ctx.Err(); err != nil {
				return llm.Usage{}, err
			}
			emit(llm.StreamChunk{Kind: llm.ChunkText, Text: part})
		}
	} else {
		emit(llm.StreamChunk{Kind: llm.ChunkToolCall, ToolCall: &llm.ToolCall{ID: "studio-time", Name: "current_time", Input: map[string]any{}}})
	}
	emit(llm.StreamChunk{Kind: llm.ChunkDone})
	return llm.Usage{}, nil
}
