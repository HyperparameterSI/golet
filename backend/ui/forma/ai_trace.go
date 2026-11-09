// Author: HyperparameterSI
package forma

import (
	"fmt"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
)

// AIStatus describes displayed work, never starts that work.
type AIStatus string

const (
	AIPending   AIStatus = "pending"
	AIRunning   AIStatus = "running"
	AIComplete  AIStatus = "complete"
	AIError     AIStatus = "error"
	AICancelled AIStatus = "cancelled"
)

type AIStep struct {
	ID, Title, Description string
	Status                 AIStatus
}

func aiStatus(c *ui.Context, s AIStatus) {
	label := "Pending"
	color := c.Theme().TextMuted
	switch s {
	case AIRunning:
		label = "Running"
	case AIComplete:
		label = "Done"
	case AIError:
		label = "Failed"
		color = c.Theme().Danger
	case AICancelled:
		label = "Cancelled"
	}
	ui.Text(c, label).FontSize(fontSize(c, 12)).TextColor(color).Unselectable()
}
func aiTrace(c *ui.Context, title string, status AIStatus, open *bool, body func()) ui.Element {
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(4)
	return root.Children(func() {
		parts := ui.CollapsibleBase(c, open)
		parts.Trigger.FillWidth().Padding(8, 0).Gap(8).AlignItems(ui.Center).Label(title).Children(func() {
			chevron := "›"
			if *open {
				chevron = "⌄"
			}
			ui.Text(c, chevron).TextColor(c.Theme().TextMuted).Unselectable()
			ui.Text(c, title).Grow(1).MaxLines(1).FontSize(fontSize(c, 14))
			aiStatus(c, status)
		})
		parts.Panel(func() {
			ui.Column(c).FillWidth().MinWidth(0).MarginStart(8).Padding(8, 0, 8, 16).Gap(10).Children(body)
		})
	})
}
func stepRows(c *ui.Context, steps []AIStep) {
	for _, s := range steps {
		ui.Row(c).FillWidth().Gap(12).AlignItems(ui.Start).Children(func() {
			aiStatus(c, s.Status)
			ui.Column(c).Grow(1).Gap(4).Children(func() {
				ui.Text(c, s.Title)
				if s.Description != "" {
					ui.Text(c, s.Description).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted)
				}
			})
		})
	}
}
func stepsStatus(steps []AIStep) AIStatus {
	status := AIComplete
	if len(steps) == 0 {
		return AIPending
	}
	for _, s := range steps {
		if s.Status == AIError {
			return AIError
		}
		if s.Status == AIRunning {
			status = AIRunning
		} else if s.Status == AICancelled && status == AIComplete {
			status = AICancelled
		} else if s.Status != AIComplete && status != AIRunning {
			if s.Status != AICancelled {
				status = AIPending
			}
		}
	}
	return status
}

// ChainOfThought displays application-provided step summaries. It does not
// request, infer, or reveal a model's private reasoning.
func ChainOfThought(c *ui.Context, title string, open *bool, steps []AIStep) ui.Element {
	return aiTrace(c, title, stepsStatus(steps), open, func() { stepRows(c, steps) })
}

// Reasoning displays only the text the application explicitly supplies.
func Reasoning(c *ui.Context, text string, running bool, duration time.Duration, open *bool) ui.Element {
	title, status := "Reasoning", AIComplete
	if running {
		title, status = "Thinking…", AIRunning
	} else if duration > 0 {
		title = fmt.Sprintf("Thought for %.1fs", duration.Seconds())
	}
	return aiTrace(c, title, status, open, func() { MessageResponse(c, text, running) })
}
func Plan(c *ui.Context, title, description string, open *bool, steps []AIStep) ui.Element {
	done := 0
	for _, s := range steps {
		if s.Status == AIComplete {
			done++
		}
	}
	return aiTrace(c, fmt.Sprintf("%s · %d/%d", title, done, len(steps)), stepsStatus(steps), open, func() {
		if description != "" {
			ui.Text(c, description).TextColor(c.Theme().TextMuted)
		}
		stepRows(c, steps)
	})
}
func Task(c *ui.Context, title string, open *bool, steps []AIStep) ui.Element {
	return aiTrace(c, title, stepsStatus(steps), open, func() { stepRows(c, steps) })
}

// QueueState preserves the independent open state of the pending/done groups.
type QueueState struct{ PendingOpen, CompletedOpen bool }

func Queue(c *ui.Context, state *QueueState, items []AIStep) ui.Element {
	pending, done := []AIStep{}, []AIStep{}
	for _, s := range items {
		if s.Status == AIComplete {
			done = append(done, s)
		} else {
			pending = append(pending, s)
		}
	}
	return ui.Column(c).FillWidth().Children(func() {
		aiTrace(c, fmt.Sprintf("To do · %d", len(pending)), stepsStatus(pending), &state.PendingOpen, func() { stepRows(c, pending) })
		aiTrace(c, fmt.Sprintf("Completed · %d", len(done)), AIComplete, &state.CompletedOpen, func() { stepRows(c, done) })
	})
}

type ToolRun struct {
	Name, Input, Output, Error string
	Status                     AIStatus
}

// Tool renders a tool's lifecycle. Input/output are text supplied by the caller,
// allowing JSON, partial JSON or logs without reflection or implicit execution.
func Tool(c *ui.Context, run ToolRun, open *bool) ui.Element {
	status := run.Status
	if run.Error != "" {
		status = AIError
	}
	return aiTrace(c, run.Name, status, open, func() {
		if run.Input != "" {
			CodeBlock(c, "Input", run.Input, CodeBlockOptions{Size: Small, Wrap: true})
		}
		if run.Output != "" {
			CodeBlock(c, "Output", run.Output, CodeBlockOptions{Size: Small, Wrap: true})
		}
		if run.Error != "" {
			CodeBlock(c, "Error", run.Error, CodeBlockOptions{Size: Small, Wrap: true, Variant: Destructive})
		}
	})
}

type TerminalOutput struct {
	Command, Output string
	Running         bool
	ExitCode        *int
}

// Terminal is a read-only agent output view, not a PTY or shell executor.
func Terminal(c *ui.Context, output TerminalOutput, open *bool) ui.Element {
	status := AIPending
	title := output.Command
	if output.Running {
		status = AIRunning
	} else if output.ExitCode != nil {
		status = AIComplete
		if *output.ExitCode != 0 {
			status = AIError
		}
		title += fmt.Sprintf(" · exit %d", *output.ExitCode)
	}
	return aiTrace(c, title, status, open, func() {
		v := Default
		if status == AIError {
			v = Destructive
		}
		CodeBlock(c, "Terminal", output.Output, CodeBlockOptions{Size: Small, Variant: v})
	})
}

type SubagentRun struct {
	Name, Task, Model, Result string
	Status                    AIStatus
	Steps                     []AIStep
	Duration                  time.Duration
}

func Subagent(c *ui.Context, run SubagentRun, open *bool) ui.Element {
	return aiTrace(c, run.Name, run.Status, open, func() {
		if run.Task != "" {
			ui.Text(c, run.Task)
		}
		ui.Row(c).Gap(8).Children(func() {
			if run.Model != "" {
				Badge(c, run.Model, Secondary)
			}
			if run.Duration > 0 {
				ui.Text(c, run.Duration.String()).TextColor(c.Theme().TextMuted)
			}
		})
		stepRows(c, run.Steps)
		if run.Result != "" {
			MessageResponse(c, run.Result)
		}
	})
}

type StackFrame struct {
	Function, File, Source string
	Line                   int
	External               bool
}

func StackTrace(c *ui.Context, message string, frames []StackFrame, open *bool, onOpen func(StackFrame)) ui.Element {
	return aiTrace(c, "Failed: "+message, AIError, open, func() {
		ui.Text(c, message).TextColor(c.Theme().Danger).Selectable()
		for _, f := range frames {
			label := fmt.Sprintf("%s · %s:%d", f.Function, f.File, f.Line)
			b := Button(c, label, ButtonOptions{Variant: Ghost, Size: Small, Disabled: onOpen == nil})
			if f.External {
				b.TextColor(c.Theme().TextMuted)
			}
			if b.Clicked() && onOpen != nil {
				onOpen(f)
			}
			if f.Source != "" {
				CodeBlock(c, "", f.Source, CodeBlockOptions{Size: Small, HideCopy: true})
			}
		}
	})
}

type AITestCase struct {
	Name, Error string
	Status      AIStatus
}
type AITestSuite struct {
	ID, Name string
	Tests    []AITestCase
}
type TestResultsState struct {
	Open   bool
	Suites map[string]bool
}

func TestResults(c *ui.Context, s *TestResultsState, suites []AITestSuite) ui.Element {
	passed, failed := 0, 0
	var steps []AIStep
	for _, suite := range suites {
		for _, test := range suite.Tests {
			status := test.Status
			if test.Error != "" {
				status = AIError
			}
			steps = append(steps, AIStep{Status: status})
			if status == AIComplete {
				passed++
			}
			if status == AIError {
				failed++
			}
		}
	}
	return aiTrace(c, fmt.Sprintf("Tests · %d passed · %d failed", passed, failed), stepsStatus(steps), &s.Open, func() {
		if s.Suites == nil {
			s.Suites = map[string]bool{}
		}
		for _, suite := range suites {
			rows := make([]AIStep, 0, len(suite.Tests))
			hasError := false
			for _, test := range suite.Tests {
				status := test.Status
				if test.Error != "" {
					status = AIError
				}
				hasError = hasError || status == AIError
				rows = append(rows, AIStep{Title: test.Name, Description: test.Error, Status: status})
			}
			expanded, known := s.Suites[suite.ID]
			if !known {
				expanded = hasError
			}
			aiTrace(c, suite.Name, stepsStatus(rows), &expanded, func() { stepRows(c, rows) })
			s.Suites[suite.ID] = expanded
		}
	})
}
