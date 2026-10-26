// Author: HyperparameterSI
package forma

import (
	"fmt"
	"math"

	"github.com/HyperparameterSI/golet/backend/ui"
)

// Attachment is a file presentation. Loading and error states do not start I/O.
func Attachment(c *ui.Context, name, description, state string, action func()) ui.Element {
	e := ui.Row(c).Gap(8).Padding(8, 10).Radius(radius(c, 16)).Border(1, c.Theme().Border).Background(c.Theme().Surface).AlignItems(ui.Center)
	if state == "error" {
		e.Border(1, c.Theme().Danger.Alpha(.3))
	}
	e.Children(func() {
		if state == "uploading" || state == "processing" {
			Spinner(c)
		} else {
			ui.Text(c, "↗").FontSize(fontSize(c, 18))
		}
		ui.Column(c).Gap(2).Children(func() { Label(c, name); CardDescription(c, description).FontSize(fontSize(c, 12)) })
		if action != nil {
			Button(c, "×", ButtonOptions{Variant: Ghost, Size: IconXS}).Label("Remove attachment").OnClick(action)
		}
	})
	return e
}
func Bubble(c *ui.Context, text string, variant Variant) ui.Element {
	bg, fg, _, border := colors(c, variant)
	return ui.Column(c).Padding(10, 14).Radius(radius(c, 24)).Background(bg).TextColor(fg).Border(1, border).Children(func() { ui.Text(c, text).FontSize(fontSize(c, 14)) })
}
func Message(c *ui.Context, author string, end bool, build func()) ui.Element {
	e := ui.Column(c).Gap(8)
	if end {
		e.AlignItems(ui.End)
	}
	return e.Children(func() { ui.Text(c, author).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted); build() })
}

// MessageScroller follows new content only while already at the bottom.
func MessageScroller(c *ui.Context, state *ui.ScrollState, build func()) ui.Element {
	if state.Y >= state.MaxY-1 {
		state.Y = math.MaxFloat32
	}
	return ui.Scroll(c).TrackScroll(state).Gap(32).Children(build)
}
func Marker(c *ui.Context, text string) ui.Element {
	return ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() { Separator(c).Grow(1); CardDescription(c, text); Separator(c).Grow(1) })
}
func Trace(c *ui.Context, title string, open *bool, build func()) ui.Element {
	return ui.Column(c).Gap(8).Children(func() {
		Collapsible(c, title, open, func() {
			ui.Column(c).MarginStart(8).PaddingStart(16).Border(1, c.Theme().Border).Gap(8).Children(build)
		})
	})
}
func CodeBlock(c *ui.Context, language, code string) ui.Element {
	return ui.Column(c).Radius(radius(c, 12)).Border(1, c.Theme().Border).Clip().Children(func() {
		ui.Row(c).Padding(12, 16).AlignItems(ui.Center).Children(func() {
			CardDescription(c, language).Grow(1)
			Button(c, "Copy", ButtonOptions{Variant: Ghost, Size: Small}).OnClick(func() { c.WriteClipboard(code) })
		})
		ui.ScrollBoth(c).MaxHeight(320).Padding(16).Children(func() { ui.Text(c, code).Font("monospace").FontSize(fontSize(c, 14)).Selectable() })
	})
}

// Question is one single-choice step. Answer is an application-owned binding.
type Question struct {
	Title, Description string
	Choices            []string
	Answer             *string
	Required           bool
}

func Questionnaire(c *ui.Context, index *int, questions []Question, onComplete func()) ui.Element {
	*index = max(0, min(*index, max(0, len(questions)-1)))
	return ui.Column(c).Gap(24).Children(func() {
		if len(questions) == 0 {
			return
		}
		q := questions[*index]
		ui.Text(c, fmt.Sprintf("%d / %d", *index+1, len(questions))).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted)
		CardTitle(c, q.Title)
		CardDescription(c, q.Description)
		if q.Answer != nil {
			RadioGroup(c, func() {
				for _, choice := range q.Choices {
					Radio(c, q.Answer, choice, choice)
				}
			})
		}
		ui.Row(c).Gap(8).Children(func() {
			Button(c, "Previous", ButtonOptions{Variant: Outline, Disabled: *index == 0}).OnClick(func() { *index-- })
			label := "Next"
			if *index == len(questions)-1 {
				label = "Submit"
			}
			Button(c, label, ButtonOptions{Disabled: q.Required && (q.Answer == nil || *q.Answer == "")}).OnClick(func() {
				if *index < len(questions)-1 {
					*index++
				} else if onComplete != nil {
					onComplete()
				}
			})
		})
	})
}
