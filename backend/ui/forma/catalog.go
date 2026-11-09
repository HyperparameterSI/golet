// Author: HyperparameterSI
package forma

// AIComponents lists the native AI compositions separately from the pinned
// upstream base registry.
func AIComponents() []string {
	return []string{"ai-message", "model-selector", "prompt-input", "chain-of-thought", "chat", "checkpoint", "confirmation", "context", "file-tree", "inline-citation", "plan", "queue", "reasoning", "sources", "stack-trace", "subagent", "suggestion", "task", "terminal", "test-results", "tool"}
}

// Components returns the pinned Base UI Maia registry names in stable order.
// Presence in the catalog is not a claim of pixel-level parity; see README.md.
func Components() []string {
	return []string{"accordion", "alert-dialog", "alert", "aspect-ratio", "attachment", "avatar", "badge", "breadcrumb", "bubble", "button-group", "button", "calendar", "card", "carousel", "chart", "checkbox", "code-block", "collapsible", "combobox", "command", "context-menu", "dialog", "direction", "drawer", "dropdown-menu", "empty", "field", "form", "hover-card", "input-group", "input-otp", "input", "item", "kbd", "label", "marker", "menubar", "message-scroller", "message", "native-select", "navigation-menu", "pagination", "popover", "progress", "questionnaire", "radio-group", "resizable", "scroll-area", "select", "separator", "sheet", "sidebar", "skeleton", "slider", "sonner", "spinner", "switch", "table", "tabs", "textarea", "toast", "toggle-group", "toggle", "tooltip", "trace"}
}
