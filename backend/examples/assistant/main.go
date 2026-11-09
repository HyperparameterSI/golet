// Author: HyperparameterSI

// Assistant is a pure Go native chat app. Without DEEPSEEK_API_KEY it runs an
// explicitly labeled offline demonstration; it never requires Node or Python.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/plugins/dsh-go"
	"github.com/HyperparameterSI/golet/backend/plugins/dsh-go/native"
	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
	"github.com/JopenChen/dsh-go/pkg/llm"
)

func main() {
	var adapter llm.LLMAdapter = offline{}
	model := os.Getenv("DEEPSEEK_MODEL")
	if model == "" {
		model = "deepseek-chat"
	}
	title := "Golet Assistant · 离线演示"
	if key := os.Getenv("DEEPSEEK_API_KEY"); key != "" {
		var err error
		adapter, err = dshgo.NewDeepSeek(dshgo.DeepSeekOptions{APIKey: key, BaseURL: os.Getenv("DEEPSEEK_BASE_URL")})
		if err != nil {
			log.Fatal(err)
		}
		title = "Golet Assistant"
	}
	panel, err := native.New(dshgo.Options{Adapter: adapter, Model: model, System: "你是一个简洁友好的桌面助手。需要当前时间时调用 current_time。", Tools: []dshgo.Tool{{Name: "current_time", Description: "读取本机当前时间", RequiresApproval: true, Execute: func(ctx context.Context, _ map[string]any) (any, error) {
		return time.Now().Format(time.RFC3339), ctx.Err()
	}}}})
	if err != nil {
		log.Fatal(err)
	}
	defer panel.Close()
	golet.App.WhenReady(func() {
		window := golet.NewWindow(golet.WindowOptions{Title: title, Width: 1000, Height: 760, MinWidth: 620, MinHeight: 480, TitleBarStyle: golet.TitleBarHiddenInset, TrafficLightPosition: &golet.Point{X: 12, Y: 11}, Content: ui.View(func(c *ui.Context) { c.SetTheme(forma.Theme(c.Theme().Dark)); panel.View(c) })})
		panel.SetWindow(window)
	})
	if err := golet.App.Run(); err != nil {
		log.Fatal(err)
	}
}

type offline struct{}

func (offline) Name() string { return "offline-demo" }
func (offline) Chat(ctx context.Context, r llm.ChatRequest, cb func(llm.StreamChunk)) (llm.Usage, error) {
	if err := ctx.Err(); err != nil {
		return llm.Usage{}, err
	}
	last := r.Messages[len(r.Messages)-1]
	if len(last.Content) > 0 && last.Content[0].ToolResult != nil {
		cb(llm.StreamChunk{Kind: llm.ChunkText, Text: "**离线演示，未调用模型。**\n\n工具返回：" + last.Content[0].ToolResult.Content + "\n\n设置 `DEEPSEEK_API_KEY` 后重新运行即可连接真实模型。"})
	} else {
		cb(llm.StreamChunk{Kind: llm.ChunkToolCall, ToolCall: &llm.ToolCall{ID: "demo-time", Name: "current_time", Input: map[string]any{}}})
	}
	cb(llm.StreamChunk{Kind: llm.ChunkDone})
	return llm.Usage{}, nil
}
