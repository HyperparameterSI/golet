// Author: HyperparameterSI
package forma

import (
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/markdown"
)

type MessageRole string

const (
	MessageUser      MessageRole = "user"
	MessageAssistant MessageRole = "assistant"
)

type MessageFeedback string

const (
	FeedbackNone MessageFeedback = ""
	FeedbackGood MessageFeedback = "good"
	FeedbackBad  MessageFeedback = "bad"
)

// AIMessageOptions keeps effects explicit: retry/edit/feedback are application
// callbacks. Streaming disables actions until the reply finishes.
type AIMessageOptions struct {
	Author     string
	Streaming  bool
	HideCopy   bool
	Feedback   *MessageFeedback
	OnFeedback func(MessageFeedback)
	OnRetry    func()
	OnEdit     func(string)
	Footer     func()
}

func MessageContent(c *ui.Context) ui.Element { return ui.Column(c).FillWidth().MinWidth(0).Gap(10) }
func MessageHeader(c *ui.Context) ui.Element {
	return ui.Row(c).Gap(8).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted)
}
func MessageFooter(c *ui.Context, end bool) ui.Element {
	e := ui.Row(c).Gap(2).AlignItems(ui.Center)
	if end {
		e.AlignSelf(ui.End)
	}
	return e
}

// MessageResponse renders streamed Markdown in Go; fenced code uses Forma's
// CodeBlock. Supply the current accumulated text on each render.
func MessageResponse(c *ui.Context, text string, streaming ...bool) ui.Element {
	o := &markdown.Options{Theme: &markdown.Theme{FontSize: fontSize(c, 14), LineHeight: 1.6, BlockSpacing: 1, CodeSize: 1, CodeBackground: token(c, "muted", c.Theme().SurfaceHover), Code: c.Theme().Text}, CodeBlock: func(c *ui.Context, b *markdown.CodeBlockView) {
		CodeBlock(c, b.Language, b.Code(), CodeBlockOptions{LineNumbers: true, HideTitleIcon: true, Highlight: HighlightCode(b.Language), HighlightKey: b.Language})
	}}
	if len(streaming) > 0 && streaming[0] {
		return markdown.Stream(c, text, o)
	}
	return markdown.View(c, text, o)
}

type messageCopyState struct {
	at     time.Time
	failed bool
}

// AIMessage composes the native user bubble or assistant reply and its actions.
// Existing Message(c, author, end, build) remains available for custom layouts.
func AIMessage(c *ui.Context, role MessageRole, text string, options ...AIMessageOptions) ui.Element {
	o := AIMessageOptions{}
	if len(options) > 0 {
		o = options[0]
	}
	end := role == MessageUser
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(10).Role(ui.RoleGroup)
	if end {
		root.AlignItems(ui.End)
	}
	st := ui.Local(root, "forma.message-copy", func() messageCopyState { return messageCopyState{} })
	return root.Children(func() {
		if o.Author != "" {
			MessageHeader(c).Children(func() { ui.Text(c, o.Author) })
		}
		if end {
			ui.Column(c).MaxWidthPercent(80).Padding(12, 14).Radius(radius(c, 22)).Background(token(c, "muted", c.Theme().SurfaceHover)).Children(func() { ui.Text(c, text).FontSize(fontSize(c, 14)).LineHeight(1.6).Selectable() })
		} else {
			MessageContent(c).Children(func() { MessageResponse(c, text, o.Streaming) })
		}
		if o.Streaming {
			ui.Text(c, "Responding…").FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted)
		}
		MessageFooter(c, end).Children(func() {
			if !o.HideCopy {
				label := "Copy"
				if !st.at.IsZero() && c.Now().Sub(st.at) < 2*time.Second {
					label = "Copied"
					if st.failed {
						label = "Copy failed"
					}
					c.After(2*time.Second - c.Now().Sub(st.at))
				}
				icon := codeBlockIcons().copy
				if label == "Copied" {
					icon = codeBlockIcons().check
				}
				if label == "Copy failed" {
					icon = codeBlockIcons().failed
				}
				b := codeBlockAction(c, label, icon).Disabled(o.Streaming || text == "")
				if b.Clicked() {
					now := c.Now()
					c.WriteClipboardThen(text, func(err error) { st.at = now; st.failed = err != nil })
					c.After(2 * time.Second)
				}
			}
			if o.Feedback != nil && !end {
				for _, f := range []struct {
					value MessageFeedback
					label string
				}{{FeedbackGood, "Good response"}, {FeedbackBad, "Bad response"}} {
					on := *o.Feedback == f.value
					btn := ui.ToggleBase(c, &on).Label(f.label).Disabled(o.Streaming).Size(28, 28).Radius(radius(c, 6)).AlignItems(ui.Center).Justify(ui.Center).TextColor(c.Theme().TextMuted)
					if on || btn.Hovered() {
						btn.Background(c.Theme().Text.Alpha(.06)).TextColor(c.Theme().Text)
					}
					icon := messageActionIcons().good
					if f.value == FeedbackBad {
						icon = messageActionIcons().bad
					}
					btn.Children(func() { ui.Icon(c, icon).Size(16, 16) })
					Tooltip(c, btn, f.label)
					if on != (*o.Feedback == f.value) {
						*o.Feedback = FeedbackNone
						if on {
							*o.Feedback = f.value
						}
						if o.OnFeedback != nil {
							o.OnFeedback(*o.Feedback)
						}
					}
				}
			}
			if o.OnRetry != nil && !end {
				if codeBlockAction(c, "Retry", messageActionIcons().retry).Disabled(o.Streaming).Clicked() {
					o.OnRetry()
				}
			}
			if o.OnEdit != nil && end {
				if codeBlockAction(c, "Edit", messageActionIcons().edit).Disabled(o.Streaming).Clicked() {
					o.OnEdit(text)
				}
			}
			if o.Footer != nil {
				o.Footer()
			}
		})
	})
}
