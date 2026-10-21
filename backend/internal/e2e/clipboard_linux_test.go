// Author: HyperparameterSI

//go:build linux && (amd64 || arm64)

package e2e

import (
	"bytes"
	"os/exec"
	"slices"
	"strings"
	"testing"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/transfer"
)

func TestUnifiedClipboardGTKForeignHTML(t *testing.T) {
	xclip, err := exec.LookPath("xclip")
	if err != nil {
		t.Skip("xclip is needed to exercise a non-Golet selection owner")
	}
	preserveClipboard(t)
	const html = "<b>foreign HTML 日本語</b>"
	cmd := exec.Command(xclip, "-selection", "clipboard", "-target", "text/html")
	cmd.Stdin = strings.NewReader(html)
	// xclip forks a selection owner after publishing. Its child exits when
	// cleanup restores the clipboard and GTK takes ownership back.
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	eventually(t, "foreign HTML selection", func() bool { return slices.Contains(golet.Clipboard.Formats(), transfer.HTML) })
	b, err := golet.Clipboard.ReadFormat(transfer.HTML)
	if err != nil || !bytes.Equal(b, []byte(html)) {
		t.Fatalf("foreign HTML bytes %q: %v", b, err)
	}
	if got := golet.Clipboard.ReadHTML(); got != html {
		t.Fatalf("native HTML convenience %q", got)
	}
}
