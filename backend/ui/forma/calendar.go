// Author: HyperparameterSI
package forma

import (
	"fmt"
	"time"

	"github.com/HyperparameterSI/golet/backend/i18n"
	"github.com/HyperparameterSI/golet/backend/ui"
)

// CalendarOptions controls unavailable dates and the first weekday. Locale and
// the default first weekday follow the Context's locale and region.
type CalendarOptions struct {
	Disabled     func(time.Time) bool
	FirstWeekday *time.Weekday
	// Multiple enables independent date toggles. Range takes precedence if both
	// bindings are provided. The date argument still tracks the last activation.
	Multiple *[]time.Time
	Range    *DateRange
}

// DateRange is an inclusive calendar selection. A zero End means that only
// Start has been chosen and the next activation will complete the range.
type DateRange struct{ Start, End time.Time }

func calendarSelected(d, single time.Time, o CalendarOptions) (selected, middle bool) {
	if o.Range != nil {
		r := o.Range
		if r.Start.IsZero() {
			return false, false
		}
		if sameDate(d, r.Start) || !r.End.IsZero() && sameDate(d, r.End) {
			return true, false
		}
		middle = !r.End.IsZero() && compareDate(d, r.Start) > 0 && compareDate(d, r.End) < 0
		return middle, middle
	}
	if o.Multiple != nil {
		for _, selected := range *o.Multiple {
			if sameDate(d, selected) {
				return true, false
			}
		}
		return false, false
	}
	return sameDate(d, single), false
}

func calendarActivate(d time.Time, date *time.Time, o CalendarOptions) {
	*date = d
	if o.Range != nil {
		r := o.Range
		if r.Start.IsZero() || !r.End.IsZero() || compareDate(d, r.Start) < 0 {
			*r = DateRange{Start: d}
		} else {
			r.End = d
		}
	} else if o.Multiple != nil {
		for i, selected := range *o.Multiple {
			if sameDate(d, selected) {
				*o.Multiple = append((*o.Multiple)[:i], (*o.Multiple)[i+1:]...)
				return
			}
		}
		*o.Multiple = append(*o.Multiple, d)
	}
}

// compareDate compares local calendar dates, independent of time zones or DST.
func compareDate(a, b time.Time) int {
	if a.Year() != b.Year() {
		return a.Year() - b.Year()
	}
	return a.YearDay() - b.YearDay()
}

// Calendar presents the Maia month grid with 32 DIP days and 8 DIP row gaps.
func Calendar(c *ui.Context, date *time.Time, options ...CalendarOptions) ui.Element {
	o := CalendarOptions{}
	if len(options) > 0 {
		o = options[0]
	}
	e := ui.Column(c).Padding(12).Gap(16).Width(248).Background(c.Theme().Background).Role(ui.RoleGroup)
	cursor := ui.Local(e, "calendar-month", func() time.Time {
		d := *date
		if d.IsZero() {
			d = c.Now()
		}
		return d
	})
	focusDate := ui.Local(e, "calendar-focus", func() string { return "" })
	month := time.Date(cursor.Year(), cursor.Month(), 1, 0, 0, 0, 0, cursor.Location())
	locale := c.Locale()
	title := i18n.FormatDate(locale, month, i18n.DateOptions{Skeleton: "yMMMM"})
	first := c.Region().FirstWeekday
	if o.FirstWeekday != nil {
		first = *o.FirstWeekday
	}
	start := month.AddDate(0, 0, -((int(month.Weekday()) - int(first) + 7) % 7))
	names := i18n.Weekdays(locale, i18n.Short, true)
	e.Label(title).Children(func() {
		ui.Row(c).Height(32).AlignItems(ui.Center).Children(func() {
			Button(c, "‹", ButtonOptions{Variant: Ghost, Size: IconSmall}).Label("Previous month").OnClick(func() { *cursor = month.AddDate(0, -1, 0) })
			ui.Text(c, title).FontSize(fontSize(c, 14)).FontWeight(500).Grow(1).TextAlign(ui.Center)
			Button(c, "›", ButtonOptions{Variant: Ghost, Size: IconSmall}).Label("Next month").OnClick(func() { *cursor = month.AddDate(0, 1, 0) })
		})
		ui.Column(c).Gap(8).Children(func() {
			ui.Row(c).Children(func() {
				for i := 0; i < 7; i++ {
					ui.Text(c, names[(int(first)+i)%7]).Width(32).FontSize(fontSize(c, 12.8)).TextAlign(ui.Center).TextColor(c.Theme().TextMuted)
				}
			})
			for week := 0; week < 6; week++ {
				ui.Row(c).Children(func() {
					for day := 0; day < 7; day++ {
						d := start.AddDate(0, 0, week*7+day)
						chosen, middle := calendarSelected(d, *date, o)
						v := Ghost
						if chosen {
							v = Default
						} else if sameDate(d, c.Now()) {
							v = Secondary
						}
						disabled := o.Disabled != nil && o.Disabled(d)
						b := Button(c, fmt.Sprint(d.Day()), ButtonOptions{Variant: v, Size: IconSmall, Disabled: disabled}).Label(d.Format("2006-01-02")).Checked(chosen)
						b.FontWeight(400)
						if middle {
							b.Background(c.Theme().SurfaceHover).TextColor(c.Theme().Text).Radius(radius(c, 0))
						}
						if *focusDate == d.Format("2006-01-02") {
							b.Focus()
							*focusDate = ""
						}
						if d.Month() != month.Month() && !chosen {
							b.TextColor(c.Theme().TextMuted)
						}
						b.OnClick(func() { calendarActivate(d, date, o); *cursor = d })
						b.HandleInput(func(ev ui.InputEvent) bool {
							if ev.Kind != ui.InputKeyDown {
								return false
							}
							delta := 0
							switch ev.Key {
							case ui.KeyLeft:
								delta = -1
							case ui.KeyRight:
								delta = 1
							case ui.KeyUp:
								delta = -7
							case ui.KeyDown:
								delta = 7
							case ui.KeyHome:
								delta = -((int(d.Weekday()) - int(first) + 7) % 7)
							case ui.KeyEnd:
								delta = 6 - ((int(d.Weekday()) - int(first) + 7) % 7)
							default:
								return false
							}
							next := d.AddDate(0, 0, delta)
							if o.Disabled == nil || !o.Disabled(next) {
								*cursor = next
								*focusDate = next.Format("2006-01-02")
							}
							return true
						})
					}
				})
			}
		})
	})
	return e
}
func sameDate(a, b time.Time) bool { return a.Year() == b.Year() && a.YearDay() == b.YearDay() }
