// Author: HyperparameterSI
package forma

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/editor"
)

// HighlightCode adapts Golet's pure Go language tokenizer. Construct it once
// outside View and pass it as CodeBlockOptions.Highlight; plain blocks never
// call this function or initialize a language registry. Unknown languages
// return nil, retaining the plain-text path.
func HighlightCode(language string) CodeHighlighter {
	lang := editor.LanguageNamed(language)
	if lang == nil || lang.Tokenizer == nil {
		return nil
	}
	return func(line string, state uint64, theme *ui.Theme) ([]ui.Span, uint64) {
		tokens, next := lang.Tokenizer.Tokenize(line, state, nil)
		colors := editor.DefaultTheme(theme)
		spans := make([]ui.Span, 0, len(tokens)+1)
		start := 0
		for _, token := range tokens {
			end := max(start, min(token.End, len(line)))
			if end > start {
				spans = append(spans, ui.Span{Text: line[start:end], Color: colors.TokenColor(token.Kind)})
			}
			start = end
		}
		if start < len(line) {
			spans = append(spans, ui.Span{Text: line[start:]})
		}
		return spans, next
	}
}
