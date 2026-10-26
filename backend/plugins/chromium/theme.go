// Author: HyperparameterSI
package chromium

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"runtime"
)

// Theme styles the macOS Alloy shell. Colors are eight-digit RRGGBBAA hex
// strings. Page content and Chrome-style windows retain their own styling.
type Theme struct {
	Dark            bool    `json:"dark"`
	Background      string  `json:"background"`
	Surface         string  `json:"surface"`
	Muted           string  `json:"muted"`
	Foreground      string  `json:"foreground"`
	MutedForeground string  `json:"mutedForeground"`
	Border          string  `json:"border"`
	Primary         string  `json:"primary"`
	Radius          float32 `json:"radius"`
	FontSize        float32 `json:"fontSize"`
	Font            string  `json:"font"`
}

func (t Theme) encode() (string, error) {
	for _, color := range []string{t.Background, t.Surface, t.Muted, t.Foreground, t.MutedForeground, t.Border, t.Primary} {
		if len(color) != 8 {
			return "", fmt.Errorf("chromium: theme colors must be RRGGBBAA")
		}
		for _, ch := range color {
			if !(ch >= '0' && ch <= '9' || ch >= 'a' && ch <= 'f' || ch >= 'A' && ch <= 'F') {
				return "", fmt.Errorf("chromium: invalid theme color %q", color)
			}
		}
	}
	if math.IsNaN(float64(t.Radius)) || math.IsInf(float64(t.Radius), 0) || t.Radius < 0 || t.Radius > 64 || math.IsNaN(float64(t.FontSize)) || math.IsInf(float64(t.FontSize), 0) || t.FontSize < 10 || t.FontSize > 20 {
		return "", fmt.Errorf("chromium: theme radius must be 0–64 and font size 10–20")
	}
	data, err := json.Marshal(t)
	return string(data), err
}

// SetTheme updates an existing macOS Alloy shell on its UI thread. It returns
// an error for unsupported styles, platforms, or hosts predating theme support.
func (b *Browser) SetTheme(ctx context.Context, theme Theme) error {
	data, err := theme.encode()
	if err != nil {
		return err
	}
	return b.call(ctx, "set-theme", data)
}

func themeArgument(opts Options) (string, error) {
	if opts.Theme == nil {
		return "", nil
	}
	style, err := opts.browserStyle()
	if err != nil {
		return "", err
	}
	if runtime.GOOS != "darwin" || style != StyleAlloy {
		return "", fmt.Errorf("chromium: custom themes require macOS Alloy")
	}
	return opts.Theme.encode()
}
