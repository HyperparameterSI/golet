# Maia 组件映射与验收清单

Author: HyperparameterSI

固定来源：2026-10-24 官方 Base UI Maia registry，65 项。每项均有原生入口与 Studio 页面；这份清单不把“页面可构建”当作“与原版完全一致”。高级能力差异见 [README.md](README.md#移植范围与验收)。

| 官方目录项 | Go 入口 | Studio 明暗构建测试 | 全状态视觉验收 |
| --- | --- | --- | --- |
| [accordion](reference/accordion.tsx) | `Accordion` | 已覆盖 | 待完成 |
| [alert-dialog](reference/alert-dialog.tsx) | `AlertDialog` | 已覆盖 | 待完成 |
| [alert](reference/alert.tsx) | `Alert` | 已覆盖 | 待完成 |
| [aspect-ratio](reference/aspect-ratio.tsx) | `AspectRatio` | 已覆盖 | 待完成 |
| [attachment](reference/attachment.tsx) | `Attachment` | 已覆盖 | 待完成 |
| [avatar](reference/avatar.tsx) | `Avatar` | 已覆盖 | 待完成 |
| [badge](reference/badge.tsx) | `Badge` | 已覆盖 | 待完成 |
| [breadcrumb](reference/breadcrumb.tsx) | `Breadcrumb` | 已覆盖 | 待完成 |
| [bubble](reference/bubble.tsx) | `Bubble` | 已覆盖 | 待完成 |
| [button-group](reference/button-group.tsx) | `ButtonGroup` | 已覆盖 | 待完成 |
| [button](reference/button.tsx) | `Button` | 已覆盖 | 待完成 |
| [calendar](reference/calendar.tsx) | `Calendar` | 已覆盖 | 待完成 |
| [card](reference/card.tsx) | `Card` | 已覆盖 | 待完成 |
| [carousel](reference/carousel.tsx) | `Carousel` | 已覆盖 | 待完成 |
| [chart](reference/chart.tsx) | `ChartContainer / ChartLegend / ChartTooltip` | 已覆盖 | 待完成 |
| [checkbox](reference/checkbox.tsx) | `Checkbox` | 已覆盖 | 待完成 |
| [code-block](reference/code-block.tsx) | `CodeBlock` | 已覆盖 | 待完成 |
| [collapsible](reference/collapsible.tsx) | `Collapsible` | 已覆盖 | 待完成 |
| [combobox](reference/combobox.tsx) | `Combobox` | 已覆盖 | 待完成 |
| [command](reference/command.tsx) | `Command` | 已覆盖 | 待完成 |
| [context-menu](reference/context-menu.tsx) | `ContextMenu` | 已覆盖 | 待完成 |
| [dialog](reference/dialog.tsx) | `Dialog` | 已覆盖 | 待完成 |
| [direction](reference/direction.tsx) | `Direction` | 已覆盖 | 待完成 |
| [drawer](reference/drawer.tsx) | `Drawer` | 已覆盖 | 待完成 |
| [dropdown-menu](reference/dropdown-menu.tsx) | `DropdownMenu` | 已覆盖 | 待完成 |
| [empty](reference/empty.tsx) | `Empty` | 已覆盖 | 待完成 |
| [field](reference/field.tsx) | `Field` | 已覆盖 | 待完成 |
| [form（上游占位项）](https://ui.shadcn.com/r/styles/base-maia/form.json) | `Form` | 已覆盖 | 无参考源码，提供原生组合 |
| [hover-card](reference/hover-card.tsx) | `HoverCard` | 已覆盖 | 待完成 |
| [input-group](reference/input-group.tsx) | `InputGroup` | 已覆盖 | 待完成 |
| [input-otp](reference/input-otp.tsx) | `InputOTP` | 已覆盖 | 待完成 |
| [input](reference/input.tsx) | `Input` | 已覆盖 | 待完成 |
| [item](reference/item.tsx) | `Item` | 已覆盖 | 待完成 |
| [kbd](reference/kbd.tsx) | `Kbd` | 已覆盖 | 待完成 |
| [label](reference/label.tsx) | `Label` | 已覆盖 | 待完成 |
| [marker](reference/marker.tsx) | `Marker` | 已覆盖 | 待完成 |
| [menubar](reference/menubar.tsx) | `Menubar` | 已覆盖 | 待完成 |
| [message-scroller](reference/message-scroller.tsx) | `MessageScroller` | 已覆盖 | 待完成 |
| [message](reference/message.tsx) | `Message` | 已覆盖 | 待完成 |
| [native-select](reference/native-select.tsx) | `NativeSelect` | 已覆盖 | 待完成 |
| [navigation-menu](reference/navigation-menu.tsx) | `NavigationMenu` | 已覆盖 | 待完成 |
| [pagination](reference/pagination.tsx) | `Pagination` | 已覆盖 | 待完成 |
| [popover](reference/popover.tsx) | `Popover` | 已覆盖 | 待完成 |
| [progress](reference/progress.tsx) | `Progress` | 已覆盖 | 待完成 |
| [questionnaire](reference/questionnaire.tsx) | `Questionnaire` | 已覆盖 | 待完成 |
| [radio-group](reference/radio-group.tsx) | `RadioGroup / Radio` | 已覆盖 | 待完成 |
| [resizable](reference/resizable.tsx) | `Resizable` | 已覆盖 | 待完成 |
| [scroll-area](reference/scroll-area.tsx) | `ScrollArea` | 已覆盖 | 待完成 |
| [select](reference/select.tsx) | `Select` | 已覆盖 | 待完成 |
| [separator](reference/separator.tsx) | `Separator` | 已覆盖 | 待完成 |
| [sheet](reference/sheet.tsx) | `Sheet` | 已覆盖 | 待完成 |
| [sidebar](reference/sidebar.tsx) | `Sidebar` | 已覆盖 | 待完成 |
| [skeleton](reference/skeleton.tsx) | `Skeleton` | 已覆盖 | 待完成 |
| [slider](reference/slider.tsx) | `Slider` | 已覆盖 | 待完成 |
| [sonner](reference/sonner.tsx) | `Sonner` | 已覆盖 | 待完成 |
| [spinner](reference/spinner.tsx) | `Spinner` | 已覆盖 | 待完成 |
| [switch](reference/switch.tsx) | `Switch` | 已覆盖 | 待完成 |
| [table](reference/table.tsx) | `Table` | 已覆盖 | 待完成 |
| [tabs](reference/tabs.tsx) | `Tabs` | 已覆盖 | 待完成 |
| [textarea](reference/textarea.tsx) | `Textarea` | 已覆盖 | 待完成 |
| [toast](reference/toast.tsx) | `Toast` | 已覆盖 | 待完成 |
| [toggle-group](reference/toggle-group.tsx) | `ToggleGroup` | 已覆盖 | 待完成 |
| [toggle](reference/toggle.tsx) | `Toggle` | 已覆盖 | 待完成 |
| [tooltip](reference/tooltip.tsx) | `Tooltip` | 已覆盖 | 待完成 |
| [trace](reference/trace.tsx) | `Trace` | 已覆盖 | 待完成 |
