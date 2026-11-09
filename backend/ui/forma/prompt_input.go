// Author: HyperparameterSI
package forma

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"slices"
	"strings"
)

// PromptAttachment is application-owned metadata, not an uploaded file. Path
// may be a local path or an application resource ID; the control never reads it.
type PromptAttachment struct{ ID, Name, Path, Description string }
type PromptInputState struct {
	Text, Model, Effort string
	Attachments         []PromptAttachment
}

// PromptSubmission is a snapshot; editing/removing attachments after sending
// cannot change the request handed to the application.
type PromptSubmission struct {
	Text, Model, Effort string
	Attachments         []PromptAttachment
}
type PromptStatus string

const (
	PromptReady     PromptStatus = ""
	PromptSubmitted PromptStatus = "submitted"
	PromptStreaming PromptStatus = "streaming"
)

type PromptInputOptions struct {
	Placeholder   string
	InputRef      *ui.Handle
	HideAttach    bool // Hide the optional attachment toolbar button.
	Disabled      bool
	Status        PromptStatus
	Error         string
	Models        []ModelOption
	OnAttach      func()
	OnDropFiles   func([]string)
	OnRemove      func(PromptAttachment)
	OnSubmit      func(PromptSubmission)
	OnStop        func()
	OnModelChange func(model, effort string)
	Tools         func()
}

// PromptInput is a native chat composer. Enter submits and Shift+Enter inserts
// a newline; IME composition is handled by the native text editor. Callbacks
// run on the UI thread: start network/file work on a worker and use Window.Update
// for results. The application owns clearing input and setting Status/Error.
func PromptInput(c *ui.Context, s *PromptInputState, options ...PromptInputOptions) ui.Element {
	o := PromptInputOptions{}
	if len(options) > 0 {
		o = options[0]
	}
	if o.Placeholder == "" {
		o.Placeholder = "Ask the agent to change your code…"
	}
	busy := o.Status == PromptSubmitted || o.Status == PromptStreaming
	canSend := func(text string) bool {
		if o.Disabled || busy || o.OnSubmit == nil || (strings.TrimSpace(text) == "" && len(s.Attachments) == 0) {
			return false
		}
		if len(o.Models) > 0 {
			for _, m := range o.Models {
				if m.ID == s.Model && !m.Disabled {
					return true
				}
			}
			return false
		}
		return true
	}
	submit := func(text string) {
		if canSend(text) {
			o.OnSubmit(PromptSubmission{text, s.Model, s.Effort, slices.Clone(s.Attachments)})
		}
	}
	root := ui.Column(c).FillWidth().MinWidth(0).Padding(12).Gap(12).Radius(radius(c, 20)).Background(token(c, "muted", c.Theme().SurfaceHover)).Role(ui.RoleGroup).Label("Prompt input")
	if o.Error != "" {
		root.Border(1, c.Theme().Danger.Alpha(.4))
	}
	if o.OnDropFiles != nil && !o.Disabled && !busy {
		if paths := root.DroppedFiles(); len(paths) > 0 {
			o.OnDropFiles(slices.Clone(paths))
		}
	}
	return root.Children(func() {
		if len(s.Attachments) > 0 {
			ui.Scroll(c).MaxHeight(144).FillWidth().Children(func() {
				remove := -1
				for i, a := range s.Attachments {
					ui.Row(c).FillWidth().Gap(8).Padding(6, 10).Radius(radius(c, 10)).Background(c.Theme().Background).AlignItems(ui.Center).Children(func() {
						ui.Column(c).Grow(1).MinWidth(0).Children(func() {
							ui.Text(c, a.Name).MaxLines(1).FontSize(fontSize(c, 13))
							if a.Description != "" {
								ui.Text(c, a.Description).FontSize(fontSize(c, 11)).TextColor(c.Theme().TextMuted)
							}
						})
						if Button(c, "×", ButtonOptions{Variant: Ghost, Size: IconXS, Disabled: o.Disabled || busy}).Label("Remove attachment " + a.Name).Clicked() {
							remove = i
						}
					})
				}
				if remove >= 0 {
					a := s.Attachments[remove]
					s.Attachments = slices.Delete(s.Attachments, remove, remove+1)
					if o.OnRemove != nil {
						o.OnRemove(a)
					}
				}
			})
		}
		input := ui.TextAreaBase(c, &s.Text).FillWidth().MinWidth(0).Height(56).Padding(2).FontSize(fontSize(c, 14)).TextColor(c.Theme().Text).Placeholder(o.Placeholder).Label("Prompt text").Disabled(o.Disabled || busy).FocusRing(false).OnTextKey(func(k ui.TextKey) bool {
			if k.Key == ui.KeyEnter && k.Mods == 0 {
				// TextKey.Text includes the last edit even before the binding updates.
				s.Text = k.Text()
				submit(s.Text)
				k.SetText(s.Text)
				return true
			}
			return false
		})
		if o.InputRef != nil {
			input.Ref(o.InputRef)
		}
		ui.Row(c).FillWidth().Gap(8).AlignItems(ui.Center).Children(func() {
			if !o.HideAttach {
				if modelCircle(c, "Attach files", modelControlIcons().add, false).Disabled(o.Disabled || busy || o.OnAttach == nil).Clicked() && o.OnAttach != nil {
					o.OnAttach()
				}
			}
			if o.Tools != nil {
				o.Tools()
			}
			ui.Spacer(c)
			if len(o.Models) > 0 {
				ModelSelector(c, &s.Model, &s.Effort, o.Models, ModelSelectorOptions{Disabled: o.Disabled || busy, OnChange: o.OnModelChange})
			}
			if busy {
				stop := ui.ButtonBase(c).Size(32, 32).Radius(16).Background(c.Theme().TextMuted).TextColor(c.Theme().Background).Label("Stop response").Disabled(o.Disabled || o.OnStop == nil).AlignItems(ui.Center).Justify(ui.Center)
				stop.Children(func() { ui.Box(c).Size(10, 10).Radius(2).Background(c.Theme().Background) })
				Tooltip(c, stop, "Stop response")
				if stop.Clicked() && o.OnStop != nil {
					o.OnStop()
				}
			} else {
				if modelCircle(c, "Send prompt", modelControlIcons().send, true).Disabled(!canSend(s.Text)).Clicked() {
					submit(s.Text)
				}
			}
		})
		if busy {
			label := "Waiting for response…"
			if o.Status == PromptStreaming {
				label = "Responding…"
			}
			ui.Text(c, label).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted)
		}
		if o.Error != "" {
			ui.Text(c, o.Error).FontSize(fontSize(c, 12)).TextColor(c.Theme().Danger)
		}
	})
}
