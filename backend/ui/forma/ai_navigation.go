// Author: HyperparameterSI
package forma

import (
	"fmt"
	"math"
	"net/url"
	"strings"

	"github.com/HyperparameterSI/golet/backend/ui"
)

// Chat follows new content only while at the live edge. Scroll state belongs to
// the caller; jumping to the latest message is an explicit user action.
func Chat(c *ui.Context, state *ui.ScrollState, build func()) ui.Element {
	return ui.Column(c).FillWidth().MinHeight(0).Children(func() {
		MessageScroller(c, state, build).Grow(1).FillWidth()
		if state.Y < state.MaxY-1 {
			if Button(c, "Jump to latest", ButtonOptions{Variant: Outline, Size: Small}).AlignSelf(ui.Center).Clicked() {
				state.Y = math.MaxFloat32
			}
		}
	})
}
func Checkpoint(c *ui.Context, label string, onRestore func()) ui.Element {
	return ui.Row(c).FillWidth().Gap(12).AlignItems(ui.Center).Children(func() {
		Separator(c).Grow(1)
		ui.Text(c, label).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted)
		if Button(c, "Restore", ButtonOptions{Variant: Ghost, Size: Small, Disabled: onRestore == nil}).Clicked() && onRestore != nil {
			onRestore()
		}
		Separator(c).Grow(1)
	})
}

type ConfirmationState string

const (
	ConfirmationPending  ConfirmationState = ""
	ConfirmationApproved ConfirmationState = "approved"
	ConfirmationDenied   ConfirmationState = "denied"
)

// Confirmation records one explicit decision. To ask again, the application
// must create/reset the state for a new approval request.
func Confirmation(c *ui.Context, title, details string, state *ConfirmationState, onDecision func(bool)) ui.Element {
	return ui.Column(c).FillWidth().Padding(16).Gap(12).Radius(radius(c, 16)).Border(1, c.Theme().Border).Background(c.Theme().Surface).Children(func() {
		ui.Text(c, title).Bold()
		if details != "" {
			ui.Scroll(c).MaxHeight(160).Children(func() { ui.Text(c, details).Selectable() })
		}
		if *state != ConfirmationPending {
			label := "Approved"
			if *state == ConfirmationDenied {
				label = "Denied"
			}
			Badge(c, label, Secondary)
			return
		}
		ui.Row(c).Gap(8).Children(func() {
			if Button(c, "拒绝", ButtonOptions{Variant: Outline, Disabled: onDecision == nil}).Clicked() {
				*state = ConfirmationDenied
				onDecision(false)
			}
			if Button(c, "允许", ButtonOptions{Disabled: onDecision == nil || *state != ConfirmationPending}).Clicked() {
				*state = ConfirmationApproved
				onDecision(true)
			}
		})
	})
}

// ContextUsage contains measured usage and an application-supplied model limit.
// Limit <= 0 is unknown; this component never invents costs or model limits.
type ContextUsage struct {
	Used, Limit, Input, Output int64
	Cost                       string
}

func ContextUsagePercent(usage ContextUsage) float64 {
	if usage.Limit <= 0 {
		return 0
	}
	return max(0, min(100, float64(usage.Used)/float64(usage.Limit)*100))
}
func AIContext(c *ui.Context, usage ContextUsage) ui.Element {
	label := "Context · unknown limit"
	if usage.Limit > 0 {
		label = fmt.Sprintf("Context · %.0f%%", ContextUsagePercent(usage))
	}
	trigger := Button(c, label, ButtonOptions{Variant: Ghost, Size: Small})
	HoverCard(c, trigger, func() {
		ui.Text(c, "Context usage").Bold()
		ui.Text(c, fmt.Sprintf("%d / %d tokens", max(0, usage.Used), max(0, usage.Limit)))
		if usage.Limit > 0 {
			Progress(c, ContextUsagePercent(usage)/100)
		}
		ui.Text(c, fmt.Sprintf("Input %d · Output %d", max(0, usage.Input), max(0, usage.Output))).TextColor(c.Theme().TextMuted)
		if usage.Cost != "" {
			ui.Text(c, usage.Cost)
		}
	})
	return trigger
}

type Source struct{ ID, Title, URL, Description string }

func sourceName(s Source) string {
	if s.Title != "" {
		return s.Title
	}
	if u, err := url.Parse(s.URL); err == nil && u.Host != "" {
		return u.Host
	}
	return s.URL
}
func sourceAllowed(s Source) bool {
	u, err := url.Parse(s.URL)
	return err == nil && u.Host != "" && (u.Scheme == "https" || u.Scheme == "http")
}
func uniqueSources(sources []Source) []Source {
	out := []Source{}
	seen := map[string]bool{}
	for _, s := range sources {
		if s.URL != "" && !seen[s.URL] {
			seen[s.URL] = true
			out = append(out, s)
		}
	}
	return out
}
func sourceButton(c *ui.Context, s Source, onOpen func(Source)) ui.Element {
	b := Button(c, sourceName(s), ButtonOptions{Variant: Ghost, Size: Small, Disabled: onOpen == nil || !sourceAllowed(s)})
	if b.Clicked() && onOpen != nil {
		onOpen(s)
	}
	return b
}
func Sources(c *ui.Context, sources []Source, open *bool, onOpen func(Source)) ui.Element {
	sources = uniqueSources(sources)
	return aiTrace(c, fmt.Sprintf("Used %d sources", len(sources)), AIComplete, open, func() {
		for _, s := range sources {
			sourceButton(c, s, onOpen)
			if s.Description != "" {
				ui.Text(c, s.Description).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted)
			}
		}
	})
}
func InlineCitation(c *ui.Context, label string, source Source, onOpen func(Source)) ui.Element {
	if label == "" {
		label = sourceName(source)
	}
	b := Button(c, label, ButtonOptions{Variant: Secondary, Size: XS})
	HoverCard(c, b, func() {
		ui.Text(c, sourceName(source)).Bold()
		ui.Text(c, source.Description)
		ui.Text(c, source.URL).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted)
		sourceButton(c, source, onOpen)
	})
	return b
}
func Suggestion(c *ui.Context, prompt string, onSelect func(string)) ui.Element {
	b := Button(c, prompt, ButtonOptions{Variant: Outline, Size: Small, Disabled: onSelect == nil}).Shrink(0)
	if b.Clicked() && onSelect != nil {
		onSelect(prompt)
	}
	return b
}
func Suggestions(c *ui.Context, prompts []string, onSelect func(string)) ui.Element {
	return ui.ScrollHorizontal(c).Children(func() {
		ui.Row(c).Gap(8).Children(func() {
			for _, prompt := range prompts {
				Suggestion(c, prompt, onSelect)
			}
		})
	})
}

type FileNode struct {
	ID, Name, Path, Change string
	Directory              bool
	Children               []FileNode
}
type FileTreeState struct {
	Open     map[string]bool
	Selected string
}

// FileTree renders metadata only. IDs must be unique; open callbacks decide how
// to read/open files. Depth is capped at 64 to bound malformed recursive data.
func FileTree(c *ui.Context, state *FileTreeState, nodes []FileNode, onOpen func(FileNode)) ui.Element {
	if state.Open == nil {
		state.Open = map[string]bool{}
	}
	var render func([]FileNode, int)
	render = func(nodes []FileNode, depth int) {
		if depth >= 64 {
			ui.Text(c, "Maximum folder depth reached")
			return
		}
		for _, n := range nodes {
			if n.Directory || len(n.Children) > 0 {
				expanded := state.Open[n.ID]
				ui.TreeItem(c, n.Name, &expanded, func() { render(n.Children, depth+1) })
				state.Open[n.ID] = expanded
			} else {
				label := n.Name
				if n.Change != "" {
					label += " · " + strings.ToUpper(n.Change)
				}
				b := Button(c, label, ButtonOptions{Variant: Ghost, Size: Small}).Label("File " + n.ID)
				if state.Selected == n.ID {
					b.Background(c.Theme().SurfaceHover)
				}
				if b.Clicked() {
					state.Selected = n.ID
					if onOpen != nil {
						onOpen(n)
					}
				}
			}
		}
	}
	return ui.Tree(c, func() { render(nodes, 0) })
}
