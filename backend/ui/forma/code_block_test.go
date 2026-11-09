// Author: HyperparameterSI
package forma

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"testing"
	"time"
)

func TestCodeBlockCopyAndExpand(t *testing.T) {
	const code = "func main() {\n\n    println(\"你好\")\n}\n"
	for _, dark := range []bool{false, true} {
		expanded := false
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(Theme(dark))
			CodeBlock(c, "main.go", code, CodeBlockOptions{LineNumbers: true, StartLine: 98, Expanded: &expanded, CollapsedHeight: 48, Lines: map[int]CodeLineVariant{2: CodeLineAdded}})
		}, 480, 300)
		if err := tt.Click("Copy"); err != nil {
			t.Fatal(err)
		}
		if tt.Clipboard() != code {
			t.Fatalf("copy changed source: %q", tt.Clipboard())
		}
		if _, ok := tt.Find("Copied"); !ok {
			t.Fatal("missing copy feedback")
		}
		tt.Advance(2 * time.Second)
		if _, ok := tt.Find("Copy"); !ok {
			t.Fatal("copy feedback did not expire")
		}
		if err := tt.Click("Show more"); err != nil {
			t.Fatal(err)
		}
		if !expanded || !tt.HasText("Show less") {
			t.Fatal("did not expand")
		}
		if !tt.HasText("100") {
			t.Fatal("missing start line numbering")
		}
	}
}

func TestCodeBlockIncrementalHighlight(t *testing.T) {
	calls := 0
	highlight := func(line string, state uint64, theme *ui.Theme) ([]ui.Span, uint64) {
		calls++
		return []ui.Span{{Text: line}}, state + 1
	}
	st := codeBlockState{}
	theme := Theme(false)
	o := CodeBlockOptions{Highlight: highlight, HighlightKey: "test"}
	st.update("one\ntwo", o, theme)
	st.update("one\ntwo", o, theme)
	if calls != 2 {
		t.Fatal(calls)
	}
	st.update("one\ntwo\nthree", o, theme)
	if calls != 3 || st.ends[2] != 3 {
		t.Fatal("append re-tokenized prefix or lost lexer state")
	}
	st.update("one\nchanged\nthree", o, theme)
	if calls != 5 {
		t.Fatal("did not invalidate changed suffix")
	}
	st.update("one\nchanged\nthree", o, Theme(true))
	if calls != 8 {
		t.Fatal("theme did not invalidate colors")
	}
	st.update("", CodeBlockOptions{}, theme)
	st.update("", CodeBlockOptions{LineNumbers: true}, theme)
	if len(st.lines) != 1 {
		t.Fatal("empty code lost its line")
	}
}

func TestCodeBlockWrapAndTheme(t *testing.T) {
	wrap := false
	tt := ui.NewTester(func(c *ui.Context) {
		CodeBlock(c, "", "a very long line that should wrap in a narrow panel rather than being truncated or removed", CodeBlockOptions{Wrap: wrap, HideCopy: true})
	}, 180, 300)
	_, ok := tt.Find("Code content")
	if !ok {
		t.Fatal("content is not accessible")
	}
	wrap = true
	tt.Frame()
	if len(tt.Texts()) == 0 {
		t.Fatal("wrapped code disappeared")
	}
}

func TestCodeBlockSelectionAndKeyboardScroll(t *testing.T) {
	var state *codeBlockState
	tt := ui.NewTester(func(c *ui.Context) {
		root := CodeBlock(c, "source", "first\n\nthird\nfourth\nfifth", CodeBlockOptions{LineNumbers: true, StartLine: 99, MaxHeight: 80, Lines: map[int]CodeLineVariant{2: CodeLineAdded}})
		state = ui.Local(root, "forma.code-block", func() codeBlockState { return codeBlockState{} })
	}, 320, 200)
	r, ok := tt.Find("first")
	if !ok {
		t.Fatal("missing first line")
	}
	tt.ClickAt(r.X+2, r.Y+r.H/2)
	tt.Command("selectAll")
	tt.Command("copy")
	if got := tt.Clipboard(); got != "first\n\nthird\nfourth\nfifth" {
		t.Fatalf("selection includes decoration or loses blank lines: %q", got)
	}
	tt.Click("Code content")
	tt.Key(0, ui.KeyDown)
	if state.scroll.Y == 0 {
		t.Fatal("keyboard did not scroll")
	}
}
