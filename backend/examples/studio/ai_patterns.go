// Author: HyperparameterSI
package main

import (
	"fmt"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
)

// newAIPatternDemo keeps example state across frames without starting agent work.
func newAIPatternDemo() func(*ui.Context, string) {
	open := true
	queue := forma.QueueState{PendingOpen: true, CompletedOpen: true}
	approval := forma.ConfirmationPending
	tree := forma.FileTreeState{Open: map[string]bool{"src": true}}
	tests := forma.TestResultsState{Open: true}
	scroll := ui.ScrollState{}
	notice := "演示数据 · 不执行命令或调用模型"
	steps := []forma.AIStep{{ID: "read", Title: "读取项目结构", Description: "backend/ 与 frontend/", Status: forma.AIComplete}, {ID: "build", Title: "实现原生界面", Description: "复用 Forma 主题和交互", Status: forma.AIRunning}, {ID: "check", Title: "验证结果", Status: forma.AIPending}}
	sources := []forma.Source{{Title: "Golet", URL: "https://github.com/HyperparameterSI/golet", Description: "Go 原生桌面应用框架"}, {Title: "AI Elements", URL: "https://ui.shadcn.com/docs/ai-elements", Description: "组件交互参考"}}
	onSource := func(s forma.Source) { notice = "来源回调：" + s.URL }
	return func(c *ui.Context, name string) {
		ui.Column(c).FillWidth().Gap(16).Children(func() {
			switch name {
			case "chain-of-thought":
				forma.ChainOfThought(c, "实现过程摘要", &open, steps)
			case "plan":
				forma.Plan(c, "构建 AI 助手", "三个可追踪的实现步骤", &open, steps)
			case "task":
				forma.Task(c, "当前任务", &open, steps)
			case "queue":
				forma.Queue(c, &queue, steps)
			case "reasoning":
				forma.Reasoning(c, "优先复用 **AIMessage** 和 **PromptInput**，保持统一主题。", false, 2*time.Second, &open)
			case "tool":
				forma.Tool(c, forma.ToolRun{Name: "read_file", Input: `{"path":"backend/main.go"}`, Output: "package main\n\nfunc main() {}", Status: forma.AIComplete}, &open)
			case "terminal":
				code := 0
				forma.Terminal(c, forma.TerminalOutput{Command: "go test ./...", Output: "ok  example/backend\nok  example/ui", ExitCode: &code}, &open)
			case "subagent":
				forma.Subagent(c, forma.SubagentRun{Name: "界面检查助手", Task: "检查主题和键盘操作", Model: "Demo model", Status: forma.AIComplete, Steps: steps[:1], Result: "检查完成，已保留 **原生键盘操作**。", Duration: 3 * time.Second}, &open)
			case "stack-trace":
				forma.StackTrace(c, "configuration is missing", []forma.StackFrame{{Function: "main.loadConfig", File: "backend/main.go", Line: 24, Source: "return errors.New(\"configuration is missing\")"}, {Function: "runtime.main", File: "runtime/proc.go", Line: 285, External: true}}, &open, func(f forma.StackFrame) { notice = fmt.Sprintf("打开源码回调：%s:%d", f.File, f.Line) })
			case "test-results":
				forma.TestResults(c, &tests, []forma.AITestSuite{{ID: "ui", Name: "UI tests", Tests: []forma.AITestCase{{Name: "Theme switch", Status: forma.AIComplete}, {Name: "Required configuration", Error: "expected API key configuration", Status: forma.AIError}}}})
			case "confirmation":
				forma.Confirmation(c, "允许读取当前项目？", "工具：read_file · backend/main.go", &approval, func(yes bool) { notice = fmt.Sprintf("授权结果：%t", yes) })
				forma.Button(c, "新建授权请求", forma.ButtonOptions{Variant: forma.Outline}).OnClick(func() { approval = forma.ConfirmationPending })
			case "checkpoint":
				forma.Checkpoint(c, "添加设置窗口之前", func() { notice = "已触发恢复回调；示例不修改文件" })
			case "context":
				forma.AIContext(c, forma.ContextUsage{Used: 24576, Limit: 65536, Input: 22000, Output: 2576, Cost: "示例费用：$0.01"})
				ui.Text(c, "悬停或聚焦查看上下文用量")
			case "sources":
				forma.Sources(c, sources, &open, onSource)
			case "inline-citation":
				ui.Text(c, "Golet 支持原生 UI 和系统 WebView。")
				forma.InlineCitation(c, "Golet · 1", sources[0], onSource)
			case "suggestion":
				forma.Suggestions(c, []string{"解释项目结构", "创建设置窗口", "检查主题一致性"}, func(s string) { notice = "已选择提示词：" + s })
			case "file-tree":
				forma.FileTree(c, &tree, []forma.FileNode{{ID: "src", Name: "backend", Directory: true, Children: []forma.FileNode{{ID: "main", Name: "main.go", Path: "backend/main.go", Change: "M"}, {ID: "config", Name: "config.go", Path: "backend/config.go", Change: "A"}}}}, func(f forma.FileNode) { notice = "文件回调：" + f.Path })
			case "chat":
				forma.Chat(c, &scroll, func() {
					for i := 0; i < 8; i++ {
						forma.AIMessage(c, forma.MessageAssistant, fmt.Sprintf("第 %d 条消息：原生滚动区域支持返回最新消息。", i+1))
					}
				}).Height(300)
			}
			if name == "plan" || name == "task" || name == "queue" || name == "chain-of-thought" {
				forma.Button(c, "推进一步", forma.ButtonOptions{Variant: forma.Outline}).OnClick(func() {
					for i := range steps {
						if steps[i].Status != forma.AIComplete {
							steps[i].Status = forma.AIComplete
							return
						}
					}
					steps[1].Status = forma.AIRunning
					steps[2].Status = forma.AIPending
				})
			}
			ui.Text(c, notice).TextColor(c.Theme().TextMuted)
		})
	}
}
