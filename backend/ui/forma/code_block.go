// Author: HyperparameterSI
package forma

import (
	"strconv"
	"strings"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
)

// CodeLineVariant decorates a source line without changing copied text.
type CodeLineVariant string

const (
	CodeLineDefault     CodeLineVariant = ""
	CodeLineAdded       CodeLineVariant = "added"
	CodeLineRemoved     CodeLineVariant = "removed"
	CodeLineHighlighted CodeLineVariant = "highlighted"
)

// CodeHighlighter tokenizes one line. State flows across lines; tokens must
// reproduce the original text. Incremental results are cached by the block.
// Change HighlightKey when the language or highlighter configuration changes.
type CodeHighlighter func(line string, state uint64, theme *ui.Theme) ([]ui.Span, uint64)

// CodeBlockOptions extends the original title/code API. Lines uses zero-based
// source indices, independent of StartLine. Zero MaxHeight defaults to 320;
// a negative height removes the cap. Footer and Actions build native children.
type CodeBlockOptions struct {
	Size        Size
	Variant     Variant
	LineNumbers bool
	StartLine   int
	Wrap        bool
	// WrapToggle optionally binds the header's wrap button to application state.
	WrapToggle *bool
	// OnDownload receives the original source. The application chooses where
	// to save it; no files are written by the control itself.
	OnDownload      func(string)
	HideCopy        bool
	HideTitleIcon   bool
	MaxHeight       float32
	CollapsedHeight float32
	Expanded        *bool
	Lines           map[int]CodeLineVariant
	Highlight       CodeHighlighter
	HighlightKey    string
	Actions         func()
	Footer          func()
}

type codeBlockState struct {
	copiedAt    time.Time
	copyFailed  bool
	scroll      ui.ScrollState
	source      string
	lines       []string
	spans       [][]ui.Span
	ends        []uint64
	theme       ui.Theme
	key         string
	highlighted bool
	measured    bool
}

// CodeBlock draws selectable code, logs or a diff entirely in Go. With no
// decoration or highlighter it renders a single text node, including streaming
// output. Keep the call's position stable (or key its parent) to retain state.
func CodeBlock(c *ui.Context, title, code string, options ...CodeBlockOptions) ui.Element {
	o := CodeBlockOptions{}
	if len(options) > 0 {
		o = options[0]
	}
	if o.WrapToggle != nil {
		o.Wrap = *o.WrapToggle
	}
	pad, size, corner := float32(16), float32(14), float32(20)
	if o.Size == Small {
		pad, size, corner = 12, 12, 14
	}
	th := c.Theme()
	fg, border := th.Text, th.Text.Alpha(.025)
	if o.Variant == Destructive {
		fg, border = th.Danger, th.Danger.Alpha(.35)
	}
	root := ui.Column(c).FillWidth().Radius(radius(c, corner)).Border(1, border).Background(token(c, "muted", th.SurfaceHover)).TextColor(fg).Clip().Role(ui.RoleGroup).Label("Code block")
	st := ui.Local(root, "forma.code-block", func() codeBlockState { return codeBlockState{} })
	if !st.measured {
		st.measured = true
		c.Invalidate()
	}
	st.update(code, o, th)
	return root.Children(func() {
		actions := func() {
			icons := codeBlockIcons()
			if o.WrapToggle != nil {
				button := codeBlockAction(c, "Wrap lines", icons.wrap)
				if *o.WrapToggle {
					button.TextColor(th.Text)
				}
				if button.Clicked() {
					*o.WrapToggle = !*o.WrapToggle
				}
			}
			if o.OnDownload != nil {
				if codeBlockAction(c, "Download code", icons.download).Clicked() {
					o.OnDownload(code)
				}
			}
			if o.Actions != nil {
				o.Actions()
			}
			if !o.HideCopy {
				copied := !st.copiedAt.IsZero() && c.Now().Sub(st.copiedAt) < 2*time.Second
				label := "Copy"
				if copied {
					c.After(2*time.Second - c.Now().Sub(st.copiedAt))
					label = "Copied"
					if st.copyFailed {
						label = "Copy failed"
					}
				}
				icon := icons.copy
				if copied {
					icon = icons.check
					if st.copyFailed {
						icon = icons.failed
					}
				}
				if codeBlockAction(c, label, icon).Clicked() {
					now := c.Now()
					c.WriteClipboardThen(code, func(err error) { st.copiedAt = now; st.copyFailed = err != nil })
					c.After(2 * time.Second)
				}
			}
		}
		header := title != "" || !o.HideCopy || o.Actions != nil || o.WrapToggle != nil || o.OnDownload != nil
		if header {
			ui.Row(c).FillWidth().Padding(pad, pad-4, 0, pad).Gap(4).AlignItems(ui.Center).Children(func() {
				if title != "" && !o.HideTitleIcon {
					ui.Icon(c, codeBlockIcons().terminal).Size(18, 18).MarginEnd(4).Unselectable()
				}
				ui.Text(c, title).FontSize(fontSize(c, size)).FontWeight(500).MaxLines(1).Grow(1).Unselectable()
				actions()
			})
		}
		height := o.MaxHeight
		if height == 0 {
			height = 320
		}
		if o.Expanded != nil && !*o.Expanded {
			height = o.CollapsedHeight
			if height <= 0 {
				height = 128
			}
		}
		content := ui.ScrollBoth(c)
		if o.Wrap {
			content = ui.Scroll(c)
		}
		top := pad
		if header {
			top = 8
		}
		content.Direction(ui.LeftToRight).FillWidth().Padding(top, pad, pad, pad).Font("monospace").FontSize(fontSize(c, size)).LineHeight(1.625).Selectable().Focusable().Label("Code content")
		if height > 0 {
			content.MaxHeight(height)
		}
		content.FocusRing(content.FocusVisible())
		for _, key := range []struct {
			key  ui.Key
			x, y float32
		}{
			{ui.KeyDown, 0, 40}, {ui.KeyUp, 0, -40}, {ui.KeyRight, 40, 0}, {ui.KeyLeft, -40, 0},
			{ui.KeyPageDown, 0, 200}, {ui.KeyPageUp, 0, -200},
		} {
			if content.Shortcut(0, key.key) {
				st.scroll.X += key.x
				st.scroll.Y += key.y
			}
		}
		content.TrackScroll(&st.scroll)
		content.Children(func() {
			if !o.LineNumbers && len(o.Lines) == 0 && o.Highlight == nil {
				text := ui.Text(c, code)
				if !o.Wrap {
					text.NoWrap()
				}
				return
			}
			start := o.StartLine
			if start < 1 {
				start = 1
			}
			gutter := float32(len(strconv.Itoa(start+len(st.lines)-1)))*fontSize(c, size)*.65 + 12
			col := ui.Column(c)
			if !o.Wrap {
				col.Shrink(0).MinWidth(max(0, root.Bounds().W-2*pad-2))
			} else {
				col.FillWidth()
			}
			col.Children(func() {
				for i, line := range st.lines {
					variant := o.Lines[i]
					row := ui.Row(c).FillWidth().MinHeight(fontSize(c, size) * 1.625)
					switch variant {
					case CodeLineAdded, CodeLineHighlighted:
						row.Background(th.Text.Alpha(.06))
					case CodeLineRemoved:
						row.Background(th.Danger.Alpha(.1))
					}
					row.Children(func() {
						if len(o.Lines) > 0 {
							marker := " "
							if variant == CodeLineAdded {
								marker = "+"
							}
							if variant == CodeLineRemoved {
								marker = "−"
							}
							ui.Text(c, marker).Width(14).TextColor(th.TextMuted).Unselectable()
						}
						if o.LineNumbers {
							ui.Text(c, strconv.Itoa(start+i)).Width(gutter).TextColor(th.TextMuted).Unselectable()
						}
						var text ui.Element
						if o.Highlight != nil && len(st.spans[i]) > 0 {
							text = ui.RichText(c, st.spans[i]...)
						} else {
							text = ui.Text(c, line)
							if line == "" {
								text.Height(fontSize(c, size) * 1.625)
							}
						}
						if o.Wrap {
							text.Grow(1).MinWidth(0)
						} else {
							text.NoWrap()
						}
					})
				}
			})
		})
		if o.Footer != nil || o.Expanded != nil {
			ui.Row(c).FillWidth().Padding(0, pad, pad, pad).Gap(8).Children(func() {
				if o.Footer != nil {
					o.Footer()
				}
				if o.Expanded != nil {
					label := "Show more"
					if *o.Expanded {
						label = "Show less"
					}
					if Button(c, label, ButtonOptions{Variant: Ghost, Size: XS}).Clicked() {
						*o.Expanded = !*o.Expanded
					}
				}
			})
		}
	})
}

func (s *codeBlockState) update(code string, o CodeBlockOptions, theme *ui.Theme) {
	highlighted := o.Highlight != nil
	// Theme colors can change in place; compare the value rather than its pointer.
	sameTheme := s.theme == *theme
	if len(s.lines) > 0 && s.source == code && s.key == o.HighlightKey && s.highlighted == highlighted && sameTheme {
		return
	}
	if !o.LineNumbers && len(o.Lines) == 0 && !highlighted {
		s.source = code
		s.lines = nil
		s.spans = nil
		s.ends = nil
		return
	}
	lines := strings.Split(strings.TrimSuffix(code, "\n"), "\n")
	first := 0
	if highlighted && s.highlighted && s.key == o.HighlightKey && sameTheme {
		for first < len(lines) && first < len(s.lines) && lines[first] == s.lines[first] {
			first++
		}
	}
	if first > len(s.spans) {
		first = 0
	}
	s.spans = s.spans[:first]
	s.ends = s.ends[:first]
	state := uint64(0)
	if first > 0 {
		state = s.ends[first-1]
	}
	for _, line := range lines[first:] {
		var spans []ui.Span
		if highlighted {
			spans, state = o.Highlight(line, state, theme)
		}
		s.spans = append(s.spans, spans)
		s.ends = append(s.ends, state)
	}
	s.source, s.lines, s.theme, s.key, s.highlighted = code, lines, *theme, o.HighlightKey, highlighted
}
