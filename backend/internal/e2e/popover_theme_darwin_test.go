// Author: HyperparameterSI

//go:build darwin

package e2e

import (
	"testing"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/internal/darwin"
	"github.com/HyperparameterSI/golet/backend/ui"
)

func TestPopoverAppearanceFollowsAppInsteadOfAnchor(t *testing.T) {
	defer golet.Theme.SetSource(golet.Theme.Source())
	golet.Theme.SetSource(golet.ThemeLight)
	parent := newWindow(t, golet.WindowOptions{Width: 300, Height: 200})
	eventually(t, "anchor window visible", parent.IsVisible)
	fl := golet.NewFlyout(golet.FlyoutOptions{Parent: parent, Anchor: golet.Rectangle{X: 100, Y: 20, Width: 80, Height: 24}, Popover: true, Width: 220, Height: 120, Content: ui.View(func(c *ui.Context) { ui.Text(c, "Popover theme test") })})
	defer fl.Destroy()
	eventually(t, "popover window adopted", func() bool { return fl.NativeHandle() != 0 })
	appearance := func() (chrome, window, content string) {
		golet.RunOnMain(func() { chrome, window, content = darwin.TestPopoverAppearance(fl.NativeHandle()) })
		return
	}
	if chrome, window, content := appearance(); chrome == "" {
		t.Fatalf("native popover inherits the status-item appearance instead of following the app (handle=%x visible=%v window=%s content=%s)", fl.NativeHandle(), fl.IsVisible(), window, content)
	}
	golet.RunOnMain(func() { darwin.TestSetPopoverAnchorAppearance(fl.NativeHandle(), false) })
	check := func() {
		want := "NSAppearanceNameAqua"
		if golet.Theme.IsDark() {
			want = "NSAppearanceNameDarkAqua"
		}
		eventually(t, "popover chrome, window and content match "+want, func() bool {
			chrome, window, content := appearance()
			return chrome == want && window == want && content == want
		})
	}
	for _, source := range []golet.ThemeSource{golet.ThemeDark, golet.ThemeLight, golet.ThemeDark, golet.ThemeSystem} {
		golet.Theme.SetSource(source)
		check()
	}
	fl.Hide()
	golet.Theme.SetSource(golet.ThemeDark)
	fl.Show()
	check()
}
