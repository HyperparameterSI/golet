# Native AI Elements

Author: HyperparameterSI

Forma exposes native Go compositions for all 21 entries in the current
[shadcn AI Elements directory](https://ui.shadcn.com/docs/ai-elements).
The registry is `AIComponents()`. Studio → Forma includes a selectable demo for
 each entry, with the same theme picker as other Forma controls.

| Official entry | Go API | Behavior |
| --- | --- | --- |
| Message | AIMessage | User bubble, Markdown response, copy/edit/feedback/retry |
| Model Selector | ModelSelector, ModelComposer | Models, providers and supported effort levels |
| Prompt Input | PromptInput | Multiline input, attachments, submit/stop, model selection |
| Chain of Thought | ChainOfThought | Expandable application-provided step summaries |
| Chat | Chat | Message scrolling, follow live edge, jump to latest |
| Checkpoint | Checkpoint | Explicit restore callback |
| Confirmation | Confirmation | Single-use approve/deny state and callback |
| Context | AIContext | Usage percentage, token breakdown and caller-supplied cost |
| File Tree | FileTree | Folder expansion, selection, changed-file labels, open callback |
| Inline Citation | InlineCitation | Focus/hover source card and explicit link callback |
| Plan | Plan | Description, step states and completed count |
| Queue | Queue | Independent pending/completed groups |
| Reasoning | Reasoning | Supplied Markdown, streaming flag, measured duration |
| Sources | Sources | Deduplicated source links and descriptions |
| Stack Trace | StackTrace | Error, source frames, external-frame styling, navigation callback |
| Subagent | Subagent | Task, model, duration, steps and result |
| Suggestion | Suggestion, Suggestions | Prompt selection callbacks and horizontal overflow |
| Task | Task | Expandable steps and aggregate status |
| Terminal | Terminal | Read-only command output, running/exit status, copy |
| Test Results | TestResults | Pass/fail summary, suites, failures initially expanded |
| Tool | Tool | Input, output, errors and lifecycle status |

## Usage

Keep mutable state outside the view callback. Change it on the UI thread;
marshal asynchronous model events through the application's UI update path.

```go
// Author: HyperparameterSI
open := true
approval := forma.ConfirmationPending
steps := []forma.AIStep{
    {ID: "read", Title: "Read project", Status: forma.AIComplete},
    {ID: "verify", Title: "Verify changes", Status: forma.AIRunning},
}
view := func(c *ui.Context) {
    forma.Plan(c, "Implement settings", "Track the requested work", &open, steps)
    forma.Confirmation(c, "Allow reading main.go?", "read_file", &approval,
        func(allow bool) { /* Forward this decision to your pending tool call. */ })
}
_ = view
```

Callbacks do not implicitly execute tools, restore files, open URLs, upload
attachments, or spawn agents. Connect them to your application or `plugins/dsh-go`.
Use stable unique IDs for file nodes and test suites. Set a bounded height on
`Chat`; its scroll state belongs to the caller. Supply real context limits and
costs; unknown limits are displayed as unknown. Source callbacks accept HTTP(S)
links only. Reasoning and ChainOfThought display only text explicitly supplied
by the application, not private model reasoning.

## Scope and differences

These are native Go adaptations using Forma tokens and Golet input/layout APIs,
not embedded React or pixel-identical web replicas. Expanded state is explicitly
controlled by the caller; automatic streaming collapse and shimmer animations
are not implemented. Chat follows the live edge but does not implement the web
example's user-turn anchor. StackTrace marks external frames but does not group
and collapse them separately. Terminal is plain text, not an ANSI emulator or
interactive PTY; use the terminal plugin for that. Code highlighting uses Forma's
native lexer rather than Shiki. Older upstream elements outside the current 21
entries (including audio/canvas integrations) are not covered by this mapping.
