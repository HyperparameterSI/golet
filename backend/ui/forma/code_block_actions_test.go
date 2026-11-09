// Author: HyperparameterSI
package forma

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"testing"
)

func TestCodeBlockHeaderActions(t *testing.T) {
	wrap := true
	download := ""
	source := "$ pnpm lint\napps/web/components/message-list.tsx\n  18:22  error  'messages' is possibly 'undefined'. Add a guard before calling map  @typescript-eslint/no-unsafe-call\n  42:7  warning  React Hook React.useEffect has a missing dependency: 'scrollToBottom'  react-hooks/exhaustive-deps\n\n✖ 2 problems (1 error, 1 warning)"
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(Theme(false))
		ui.Column(c).Fill().Padding(24).Background(c.Theme().Background).Children(func() {
			CodeBlock(c, "Terminal", source, CodeBlockOptions{WrapToggle: &wrap, OnDownload: func(s string) { download = s }, MaxHeight: -1})
		})
	}, 600, 410)
	if err := tt.Click("Download code"); err != nil {
		t.Fatal(err)
	}
	if download != source {
		t.Fatal("download changed source")
	}
	tt.Click("Wrap lines")
	if wrap {
		t.Fatal("wrap did not toggle")
	}
}
