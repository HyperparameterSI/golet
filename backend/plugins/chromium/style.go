// Author: HyperparameterSI
package chromium

import "fmt"

// BrowserStyle selects the window implementation for a host process.
// Choose it at launch; changing style requires closing and relaunching the host.
type BrowserStyle string

const (
	// StyleAlloy uses Golet's custom window, tabs, and address bar (the default).
	// It supports the assistant and Golet user scripts, but not Chrome extensions.
	StyleAlloy BrowserStyle = "alloy"
	// StyleChrome uses the standard Chrome window and supports Chrome extensions.
	StyleChrome BrowserStyle = "chrome"
)

func (opts Options) browserStyle() (BrowserStyle, error) {
	switch opts.Style {
	case "":
		if opts.ChromeStyle {
			return StyleChrome, nil
		}
		return StyleAlloy, nil
	case StyleAlloy:
		if opts.ChromeStyle {
			return "", fmt.Errorf("chromium: StyleAlloy conflicts with ChromeStyle=true; use Style only")
		}
		return StyleAlloy, nil
	case StyleChrome:
		return StyleChrome, nil
	default:
		return "", fmt.Errorf("chromium: unsupported browser style %q; use StyleAlloy or StyleChrome", opts.Style)
	}
}
