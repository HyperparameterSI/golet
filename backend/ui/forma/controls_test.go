// Author: HyperparameterSI
package forma

import (
	"testing"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
)

func TestButtonSizesAndDisabled(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(Theme(false))
		ui.Column(c).Gap(8).Children(func() {
			for i, size := range []Size{XS, Small, Normal, Large} {
				Button(c, []string{"xs", "sm", "normal", "lg"}[i], ButtonOptions{Size: size}).OnClick(func() { clicks++ })
			}
			Button(c, "disabled", ButtonOptions{Disabled: true}).OnClick(func() { clicks++ })
		})
	}, 300, 250)
	for label, height := range map[string]float32{"xs": 24, "sm": 32, "normal": 36, "lg": 40} {
		r, ok := tt.Find(label)
		if !ok || r.H != height {
			t.Fatalf("%s height %v, want %v", label, r.H, height)
		}
	}
	tt.Click("disabled")
	if clicks != 0 {
		t.Fatal("disabled button activated")
	}
	tt.Click("normal")
	if clicks != 1 {
		t.Fatal("button did not activate")
	}
}
func TestSelectionKeyboard(t *testing.T) {
	value := "Alpha"
	tt := ui.NewTester(func(c *ui.Context) { Select(c, &value, []string{"Alpha", "Beta", "Gamma"}).Label("Choose") }, 320, 240)
	tt.Click("Choose")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if value != "Beta" {
		t.Fatalf("selection %q", value)
	}
}
func TestContextMenuIsNativeDrawn(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		anchor := ui.Box(c).Size(200, 100).Label("Target")
		ContextMenu(c, anchor, []MenuItem{{Label: "Inspect", Action: func() { clicks++ }}})
	}, 400, 300)
	if err := tt.RightClick("Target"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Inspect") {
		t.Fatal("native menu not visible")
	}
	tt.Click("Inspect")
	if clicks != 1 {
		t.Fatal("menu action did not run once")
	}
}
func TestDialogDismissesAndRestoresInput(t *testing.T) {
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		Button(c, "Open").OnClick(func() { open = true })
		Dialog(c, &open, "Settings", func() { Button(c, "Inside") })
	}, 600, 400)
	tt.Click("Open")
	if !open || !tt.HasText("Inside") {
		t.Fatal("dialog did not open")
	}
	tt.Key(0, ui.KeyEscape)
	if open {
		t.Fatal("escape did not close dialog")
	}
	tt.Click("Open")
	tt.Click("Close dialog")
	if open {
		t.Fatal("close button failed")
	}
}
func TestOTPFiltersPaste(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { InputOTP(c, &value, 6) }, 300, 100)
	tt.Click("One-time password")
	tt.SetClipboard("12a34-5678")
	tt.Command("paste")
	if value != "123456" {
		t.Fatalf("OTP = %q", value)
	}
}
func TestQuestionnaireRequiresAnswer(t *testing.T) {
	index := 0
	answer := ""
	submitted := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Questionnaire(c, &index, []Question{{Title: "Language", Choices: []string{"Go", "Rust"}, Answer: &answer, Required: true}}, func() { submitted++ })
	}, 400, 300)
	tt.Click("Submit")
	if submitted != 0 {
		t.Fatal("required question submitted empty")
	}
	tt.Click("Go")
	tt.Click("Submit")
	if submitted != 1 || answer != "Go" {
		t.Fatal("answer not submitted")
	}
}
func TestSkeletonReducedMotionSettles(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Skeleton(c).Size(180, 20) }, 200, 50)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Fatal("reduced motion skeleton requests frames")
	}
}

func TestCalendarArrowNavigation(t *testing.T) {
	date := time.Date(2026, 10, 24, 0, 0, 0, 0, time.UTC)
	tt := ui.NewTester(func(c *ui.Context) { Calendar(c, &date) }, 300, 380)
	tt.Click("2026-10-24")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if date.Day() != 24 {
		t.Fatal("arrow navigation must not commit a date before activation")
	}
	tt.Key(0, ui.KeyEnter)
	if date.Day() != 26 {
		t.Fatalf("date after two right arrows and Enter = %v", date)
	}
	tt.Click("Next month")
	if !tt.HasText("November 2026") {
		t.Fatal("month navigation failed")
	}
}

func TestCalendarSelectionModes(t *testing.T) {
	date := time.Date(2026, 10, 24, 0, 0, 0, 0, time.UTC)
	var dates []time.Time
	rangeValue := DateRange{}
	rangeMode := false
	tt := ui.NewTester(func(c *ui.Context) {
		o := CalendarOptions{Multiple: &dates, Disabled: func(d time.Time) bool { return d.Day() == 28 }}
		if rangeMode {
			o.Multiple = nil
			o.Range = &rangeValue
		}
		Calendar(c, &date, o)
		Button(c, "Range mode").OnClick(func() { rangeMode = true })
	}, 300, 420)
	tt.Click("2026-10-24")
	tt.Click("2026-10-26")
	tt.Click("2026-10-24")
	if len(dates) != 1 || dates[0].Day() != 26 {
		t.Fatalf("multiple selection: %v", dates)
	}
	tt.Click("2026-10-28")
	if len(dates) != 1 {
		t.Fatal("disabled day selected")
	}
	tt.Click("Range mode")
	tt.Click("2026-10-24")
	tt.Click("2026-10-27")
	if rangeValue.Start.Day() != 24 || rangeValue.End.Day() != 27 {
		t.Fatalf("range: %+v", rangeValue)
	}
	tt.Click("2026-10-23")
	if rangeValue.Start.Day() != 23 || !rangeValue.End.IsZero() {
		t.Fatalf("range did not restart: %+v", rangeValue)
	}
}

func TestBreadcrumbCurrentPageIsNotInteractive(t *testing.T) {
	selected := -1
	tt := ui.NewTester(func(c *ui.Context) { Breadcrumb(c, []string{"Home", "Library", "Current"}, &selected) }, 400, 80)
	tt.Click("Library")
	if selected != 1 {
		t.Fatalf("activated ancestor = %d", selected)
	}
	tt.Click("Current")
	if selected != 1 {
		t.Fatal("current page activated as an ancestor")
	}
}
func TestDropdownKeyboard(t *testing.T) {
	chosen := 0
	tt := ui.NewTester(func(c *ui.Context) {
		DropdownMenu(c, "Menu", []MenuItem{{Label: "First", Action: func() { chosen = 1 }}, {Label: "Second", Action: func() { chosen = 2 }}})
	}, 320, 250)
	tt.Click("Menu")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if chosen != 2 {
		t.Fatalf("chosen %d, expected second item", chosen)
	}
}
func TestCarouselKeyboard(t *testing.T) {
	index := 0
	tt := ui.NewTester(func(c *ui.Context) { Carousel(c, &index, 3, func(i int) { ui.Text(c, "Slide") }).Height(150) }, 400, 200)
	tt.Click("Slide 1 of 3")
	tt.Key(0, ui.KeyRight)
	if index != 1 {
		t.Fatalf("slide %d", index)
	}
	tt.Click("Previous slide")
	if index != 0 {
		t.Fatal("previous failed")
	}
}

func TestAlertBackdropCannotConfirmOrDismiss(t *testing.T) {
	open := true
	confirmed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		AlertDialog(c, &open, "Delete?", "This action needs confirmation.", func() { confirmed++ })
	}, 600, 400)
	tt.ClickAt(5, 5)
	if !open || confirmed != 0 {
		t.Fatal("backdrop changed alert state")
	}
	tt.Click("Cancel")
	if open || confirmed != 0 {
		t.Fatal("cancel confirmed action")
	}
}
func TestCommandRunsOnce(t *testing.T) {
	query := ""
	count := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Command(c, &query, []CommandItem{{ID: "settings", Label: "Settings", Run: func() { count++ }}})
	}, 320, 240)
	tt.Click("settings")
	tt.Frame()
	if count != 1 {
		t.Fatalf("command ran %d times", count)
	}
}
