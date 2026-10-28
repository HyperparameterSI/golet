// Author: HyperparameterSI

//go:build linux

package e2e

import (
	"testing"
	"time"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/internal/linux"
)

// TestTrayActivation checks that a panel's activation, the double click on
// GNOME's AppIndicator icons, reaches Tray.OnClick, and that an icon whose
// clicks have nowhere to go answers no activation, so that a panel opens
// its menu on the first click.
func TestTrayActivation(t *testing.T) {
	tray, err := golet.NewTray(golet.TrayOptions{Icon: squarePNG(), ToolTip: "activation"})
	if err != nil {
		t.Skip("no tray: ", err)
	}
	defer tray.Destroy()
	var activates bool
	golet.RunOnMain(func() { activates = linux.TestTrayActivates() })
	if !activates {
		t.Skip("this libayatana-appindicator3 has no activate signal")
	}

	var activated bool
	golet.RunOnMain(func() { activated = linux.TestActivateTray() })
	if activated {
		t.Fatal("an icon with no click listener answered activation")
	}

	clicked := make(chan struct{}, 1)
	off := tray.OnClick(func() { clicked <- struct{}{} })
	golet.RunOnMain(func() { activated = linux.TestActivateTray() })
	if !activated {
		t.Fatal("an icon with a click listener answered no activation")
	}
	select {
	case <-clicked:
	case <-time.After(5 * time.Second):
		t.Fatal("OnClick did not run when the icon was activated")
	}

	off()
	golet.RunOnMain(func() { activated = linux.TestActivateTray() })
	if activated {
		t.Fatal("an icon whose last click listener left answered activation")
	}
}
