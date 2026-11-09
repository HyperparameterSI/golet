// Author: HyperparameterSI
package forma

import "github.com/HyperparameterSI/golet/backend/ui"

// ModelEffort is a provider-supported reasoning setting. IDs are passed to the
// application unchanged; the UI never assumes that all providers support them.
type ModelEffort struct {
	ID, Name string
	Disabled bool
}

// ModelOption describes one selectable model. Empty Efforts hides its effort
// menu. Provider groups are displayed in first-occurrence order.
type ModelOption struct {
	ID, Name, Provider, Description string
	Disabled                        bool
	Efforts                         []ModelEffort
}
type ModelSelectorOptions struct {
	Disabled    bool
	Placeholder string
	OnChange    func(model, effort string)
}

// ModelSelector binds model and effort IDs to two keyboard-accessible menus.
// Changing models retains a supported effort or picks the first enabled one.
// It does not contact providers or persist settings. Pointers must be non-nil.
func ModelSelector(c *ui.Context, model, effort *string, models []ModelOption, options ...ModelSelectorOptions) ui.Element {
	o := ModelSelectorOptions{}
	if len(options) > 0 {
		o = options[0]
	}
	if o.Placeholder == "" {
		o.Placeholder = "Select model"
	}
	beforeModel, beforeEffort := *model, *effort
	root := ui.Row(c).Gap(8).AlignItems(ui.Center)
	root.Children(func() {
		name := o.Placeholder
		for _, m := range models {
			if m.ID == *model {
				name = m.Name
				if name == "" {
					name = m.ID
				}
				break
			}
		}
		s := ui.SelectBase(c, model)
		modelPill(c, s.Trigger).Disabled(o.Disabled || len(models) == 0).Label("Model")
		if s.Trigger.Hovered() {
			s.Trigger.Background(c.Theme().SurfaceHover)
		}
		s.Trigger.Children(func() {
			ui.Text(c, name).FontSize(fontSize(c, 14))
			ui.Icon(c, modelControlIcons().chevron).Size(14, 14)
		})
		s.Popup(func(p ui.Element) {
			panel(c, p).MinWidth(240).MaxHeight(320)
			ui.Scroll(c).Children(func() {
				groups := []string{}
				seen := map[string]bool{}
				for _, m := range models {
					if !seen[m.Provider] {
						seen[m.Provider] = true
						groups = append(groups, m.Provider)
					}
				}
				for _, group := range groups {
					if group != "" {
						ui.Text(c, group).Padding(8, 10).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted).Unselectable()
					}
					for _, m := range models {
						if m.Provider != group {
							continue
						}
						label := m.Name
						if label == "" {
							label = m.ID
						}
						item := s.Item(m.ID).Disabled(m.Disabled).Padding(8, 10).Gap(10).Radius(radius(c, 8)).Label(label)
						if item.Highlighted() {
							item.Background(c.Theme().SurfaceHover)
						}
						item.Children(func() {
							ui.Column(c).Grow(1).Gap(3).Children(func() {
								ui.Text(c, label)
								if m.Description != "" {
									ui.Text(c, m.Description).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted)
								}
							})
							if *model == m.ID {
								ui.Text(c, "✓").Unselectable()
							}
						})
					}
				}
			})
		})
		var current *ModelOption
		for i := range models {
			if models[i].ID == *model {
				current = &models[i]
				break
			}
		}
		if beforeModel != *model {
			valid := false
			if current != nil {
				for _, e := range current.Efforts {
					if e.ID == *effort && !e.Disabled {
						valid = true
					}
				}
			}
			if !valid {
				*effort = ""
				if current != nil {
					for _, e := range current.Efforts {
						if !e.Disabled {
							*effort = e.ID
							break
						}
					}
				}
			}
		}
		if current != nil && len(current.Efforts) > 0 {
			selected := "Select effort"
			for _, e := range current.Efforts {
				if e.ID == *effort {
					selected = e.Name
					if selected == "" {
						selected = e.ID
					}
				}
			}
			es := ui.SelectBase(c, effort)
			modelPill(c, es.Trigger).Disabled(o.Disabled || current.Disabled).Label("Reasoning effort")
			if es.Trigger.Hovered() {
				es.Trigger.Background(c.Theme().SurfaceHover)
			}
			es.Trigger.Children(func() {
				ui.Text(c, "Effort").TextColor(c.Theme().TextMuted)
				ui.Text(c, selected)
				ui.Icon(c, modelControlIcons().chevron).Size(14, 14)
			})
			es.Popup(func(p ui.Element) {
				panel(c, p).MinWidth(160).MaxHeight(280)
				ui.Scroll(c).Children(func() {
					for _, e := range current.Efforts {
						label := e.Name
						if label == "" {
							label = e.ID
						}
						item := es.Item(e.ID).Disabled(e.Disabled).Padding(8, 10).Gap(10).Radius(radius(c, 8)).Label(label)
						if item.Highlighted() {
							item.Background(c.Theme().SurfaceHover)
						}
						item.Children(func() {
							ui.Text(c, label).Grow(1)
							if *effort == e.ID {
								ui.Text(c, "✓")
							}
						})
					}
				})
			})
		}
	})
	if o.OnChange != nil && (beforeModel != *model || beforeEffort != *effort) {
		o.OnChange(*model, *effort)
	}
	return root
}
