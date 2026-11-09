// Author: HyperparameterSI

// Package native provides a reusable Forma chat panel, rendered entirely in Go.
package native

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/plugins/dsh-go"
	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
	"github.com/JopenChen/dsh-go/pkg/llm"
)

type message struct {
	role     forma.MessageRole
	text     string
	feedback forma.MessageFeedback
}
type approval struct {
	name      string
	arguments string
	answer    chan bool
}

// Panel owns a conversation. View and SetWindow run on the UI thread; model and
// tools run on a worker. Close is nonblocking and safe from any goroutine.
type Panel struct {
	conversation  *dshgo.Conversation
	window        *golet.Window
	status        string
	draft         forma.PromptInputState
	models        []forma.ModelOption
	notice        string
	messages      []*message
	list          ui.ListState
	busy          bool
	approval      *approval
	headerActions func(*ui.Context)
}

// New constructs a panel. When no approval callback is supplied, protected
// tools require an explicit choice in the native panel before they execute.
func New(o dshgo.Options) (*Panel, error) {
	p := &Panel{draft: forma.PromptInputState{Model: o.Model}, models: []forma.ModelOption{{ID: o.Model, Name: o.Model, Provider: "Configured model"}}}
	if o.Approve == nil {
		o.Approve = p.approve
	}
	var err error
	p.conversation, err = dshgo.New(o)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// SetWindow connects UI updates and cancels work on window close. Call once,
// immediately after NewWindow; constructing the first frame needs no attachment.
func (p *Panel) SetWindow(w *golet.Window) { p.window = w; w.OnClosed(func() { _ = p.Close() }) }
func (p *Panel) Close() error              { return p.conversation.Close() }

// SetNotice sets a persistent mode description. Call on the UI thread.
func (p *Panel) SetNotice(text string) { p.notice = text }

// SetHeaderActions adds non-draggable application controls on the UI thread.
func (p *Panel) SetHeaderActions(build func(*ui.Context)) { p.headerActions = build }

// Busy reports whether a turn is running. Call on the UI thread.
func (p *Panel) Busy() bool { return p.busy }

// Configure replaces an idle conversation and clears its history. Validation
// happens before replacing the old conversation. Call on the UI thread.
func (p *Panel) Configure(o dshgo.Options) error {
	if p.busy {
		return dshgo.ErrBusy
	}
	if o.Approve == nil {
		o.Approve = p.approve
	}
	conversation, err := dshgo.New(o)
	if err != nil {
		return err
	}
	_ = p.conversation.Close()
	p.conversation = conversation
	p.draft = forma.PromptInputState{Model: o.Model}
	p.models = []forma.ModelOption{{ID: o.Model, Name: o.Model, Provider: "Configured model"}}
	p.messages = nil
	p.status = ""
	p.approval = nil
	return nil
}

func (p *Panel) approve(ctx context.Context, call llm.ToolCall) (bool, error) {
	if p.window == nil {
		return false, errors.New("assistant: panel has no window")
	}
	raw, _ := json.MarshalIndent(call.Input, "", "  ")
	a := &approval{name: call.Name, arguments: string(raw), answer: make(chan bool, 1)}
	p.window.Update(func() { p.approval = a })
	defer p.window.Update(func() {
		if p.approval == a {
			p.approval = nil
		}
	})
	select {
	case answer := <-a.answer:
		return answer, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

// View builds a chat area which follows the surrounding Forma/UI theme.
func (p *Panel) View(c *ui.Context) {
	t := c.Theme()
	p.list.FollowEnd = true
	ui.Column(c).Fill().Background(t.Background).Children(func() {
		ui.Row(c).FillWidth().Height(max(c.TitleBar().Height, 42)).Padding(8, c.TitleBar().Right+16, 8, c.TitleBar().Left+16).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Golet Assistant").Bold().Grow(1).DragWindow()
			if p.headerActions != nil {
				p.headerActions(c)
			}
		})
		if p.notice != "" {
			ui.Text(c, p.notice).FontSize(12).TextColor(t.TextMuted).Padding(8, 20)
		}
		if len(p.messages) == 0 {
			ui.Text(c, "试试：现在几点？或请帮我写一段 Go 代码。").TextColor(t.TextMuted).Padding(12, 20)
		}
		ui.List(c, &p.list, len(p.messages), func(i int) {
			m := p.messages[i]
			ui.Column(c).FillWidth().Padding(12).Children(func() {
				o := forma.AIMessageOptions{Streaming: p.busy && i == len(p.messages)-1 && m.role == forma.MessageAssistant}
				if m.role == forma.MessageAssistant {
					o.Feedback = &m.feedback
				} else if !p.busy {
					o.OnEdit = func(text string) { p.draft.Text = text }
				}
				forma.AIMessage(c, m.role, m.text, o)
			})
		}).Grow(1).Padding(12, 20)
		if p.approval != nil {
			a := p.approval
			forma.Card(c).Children(func() {
				ui.Text(c, "允许执行工具："+a.name+"？").Bold()
				ui.Scroll(c).MaxHeight(160).Children(func() { ui.Text(c, a.arguments) })
				ui.Text(c, "仅允许本次调用").TextColor(t.TextMuted)
				ui.Row(c).Gap(8).Children(func() {
					for _, v := range []struct {
						label string
						allow bool
					}{{"允许", true}, {"拒绝", false}} {
						if forma.Button(c, v.label).Clicked() {
							select {
							case a.answer <- v.allow:
							default:
							}
							p.approval = nil
						}
					}
				})
			})
		}
		if p.status != "" {
			ui.Text(c, p.status).TextColor(t.TextMuted).Padding(4, 20)
		}
		ui.Column(c).FillWidth().Gap(8).Padding(12, 20).Children(func() {
			status := forma.PromptReady
			if p.busy {
				status = forma.PromptStreaming
			}
			forma.PromptInput(c, &p.draft, forma.PromptInputOptions{Placeholder: "输入问题，Enter 发送，Shift+Enter 换行…", Models: p.models, Status: status, Disabled: p.window == nil, OnSubmit: func(request forma.PromptSubmission) { p.send(request.Text) }, OnStop: func() { p.conversation.Cancel() }})
			if !p.busy {
				if forma.Button(c, "清空", forma.ButtonOptions{Variant: forma.Ghost}).Clicked() {
					if p.conversation.Reset() == nil {
						p.messages = nil
						p.status = ""
					}
				}
			}
		})
	})
}

func (p *Panel) send(input string) {
	if p.busy || p.window == nil {
		return
	}
	p.draft.Text = ""
	p.busy = true
	p.status = "正在回复…"
	answer := &message{role: forma.MessageAssistant}
	p.messages = append(p.messages, &message{role: forma.MessageUser, text: input}, answer)
	if len(p.messages) > 40 {
		p.messages = p.messages[len(p.messages)-40:]
	}
	go func() {
		result, err := p.conversation.Run(context.Background(), input, func(e dshgo.Event) {
			p.window.Update(func() {
				switch e.Kind {
				case "text":
					answer.text += e.Text
				case "tool_start":
					p.status = "调用工具：" + e.Tool
				case "tool_end":
					p.status = "工具完成：" + e.Tool
				}
			})
		})
		p.window.Update(func() {
			p.busy = false
			p.status = fmt.Sprintf("%d 次模型步骤 · 输入 %d / 输出 %d tokens", result.Steps, result.Usage.PromptTokens, result.Usage.CompletionTokens)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					p.status = "已停止"
				} else {
					p.status = err.Error()
				}
			}
		})
	}()
}
