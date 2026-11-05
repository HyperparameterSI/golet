# MyGo v0.7.4 同步记录

<!-- Author: HyperparameterSI -->

2026-11-05，以 Golet 0.1.14 已同步的 `85fd4bbad893199869f975156d3cfe704fc497a8`
为基线，移植至 MyGo v0.7.4：`9bd34e1eb9525d9a1fe698aadc964c49c8ff8e7e`。
本次目标是指定的稳定标签，不包含该标签之后的 main 分支开发提交。
[上游发布说明](https://github.com/egoist/mygo/releases/tag/v0.7.4) ·
[源码差异](https://github.com/egoist/mygo/compare/85fd4bbad893199869f975156d3cfe704fc497a8...v0.7.4)。
保留 MIT 来源和许可证；这些更新随 Golet 0.1.15 发布。

## 更新内容

| 范围 | 已同步功能与修复 |
| --- | --- |
| 窗口首帧 | Windows 原生窗口首帧完成后再显示，避免白闪；三平台 Hidden 原生窗口可先绘制，通过 OnReadyToShow 再显示 |
| macOS 空闲调度 | 仅在帧后统计必要的事件循环轮次，避免空闲窗口持续唤醒 Go；未将上游功耗数字视为 Golet 实测 |
| 局部主题 | `ui.Themed` 为子树提供独立颜色、字体、间距和绘制主题，弹层继承对应主题 |
| 输入 | `DecimalInput` 用字符串精确表示小数并验证范围/精度；`PlaceholderColor`、`CaretWidth`；辅助技术可编辑自定义文本客户端 |
| 绘图 | `RadialGradient`、`Painter.FillRadialGradient`；同步 CPU、Metal、D3D11、OpenGL 着色器实现及预编译资源 |
| 控件 | `IgnoreInput`、`Calendar.Browse`、`Scrollbars(false)`、`DottedUnderline`、`Context.Instant`；修复网格子元素百分比相对单元格计算 |
| 滚动条 | `ui.Scrollbar`/`Painter.Scrollbar`；编辑器、富文本、终端统一主题宽度、悬停与高对比度外观 |
| 编辑菜单 | Undo/Redo/Cut/Copy/Paste/Delete/Select All 跟随当前焦点的操作能力；密码与只读输入遵循各自限制 |
| WebView | `ScrollToContainer` 将网页上的滚轮/触控板事件转交外层原生滚动容器；macOS 非活动窗口首击可到达页面 |
| Windows WebView2 | 支持附加浏览器参数与用户数据目录环境变量；启动时正确应用强制深浅色主题 |
| 对话框/通知 | `DestructiveButtons`；macOS Escape/Command+. 取消按钮行为修复；`Notification.ShowContext` 可取消授权等待 |
| 更新器 | `AvailableVersion()` 查询已发现更新；限制大文件差分匹配工作量，避免生成 delta 时卡住 |
| 编辑器 | 选区超出左右边缘时自动横向滚动；Diff 分隔区域滚轮、水平滚动条、填充条纹及长短悬殊区段滚动修复 |
| Markdown | `Options.Anchor` 定位标题；`Theme.LinkWeight` 设置链接字重 |
| 终端 | 绘制 block sextants，修复 Expo 等工具的二维码显示 |
| 平台修复 | macOS 多显示器托盘弹窗位置；Windows 10 无边框窗口边缘、绘制菜单文本裁剪、剪贴板格式回退 |
| Linux 分发 | 生成 Arch Linux `PKGBUILD`；`linux.command` 支持单个名称或名称数组 |
| 构建工具 | `-skip-installer`；Windows 安装器桌面快捷方式文案本地化；扫描器占用文件时重试最终重命名；macOS 27 磁盘镜像转换回退 |
| 诊断工具 | idlemem 清理 Windows 应用进程并显示退出原因；文件监听测试清理修复 |

## Golet 适配

- 保留 backend/frontend 目录结构、根模块路径和现有 Golet 发布版本。
- 保留 Forma、Browser、Scene3D、BorderBeam、结构化日志及增强玻璃主题。
- `Painter.Theme()` 使用局部主题，同时保留 Golet 的 `themeMaterial` 适配器，继续支持插件材质构建。
- Studio 的“表单与选择”新增精确金额输入；“样式与布局”新增局部深浅色主题预览和径向渐变。
- 保留 Gallery 作为共享模块的 Open 入口，合并新增日历 Browse 演示。
- 局部主题尺寸测试显式指定外层圆角，避免把 MyGo 默认 6px 当作 Golet 的默认样式。
- 更新 CLI TypeScript 配置声明、文档演示与翻译；不复制上游品牌、第三方 showcase 和 npm 版本号。

## 验证和体积

环境：macOS arm64 / Apple M5 / Go 1.27.3，`CGO_ENABLED=0`。

- 全量 Go 测试，包括 Forma、Scene3D、玻璃插件、新主题和输入/滚动行为。
- 三平台 build/vet、macOS 生产模式 vet、js/wasm vet。
- Intel macOS、Windows ARM64、Windows/Linux 386、FreeBSD 交叉编译。
- Bun 工作区构建、226 项测试、TypeScript 类型检查。
- macOS 实机通过隐藏原生窗口首帧、WebView 滚动转交、WebView/原生编辑菜单状态、首次点击、原生按钮对话框测试。

Windows/Linux 的实际运行表现仍需对应系统或 CI；交叉编译不等于实机验证。

基础应用以同一工具链和 `-trimpath -tags golet_noinspector`、
`-ldflags '-s -w -X github.com/HyperparameterSI/golet/backend.production=1'` 对比：

| 应用 | 同步前（字节） | 同步后（字节） | 增量 |
| --- | ---: | ---: | ---: |
| counter-native | 9,472,386 | 9,505,890 | 33,504（约 0.35%） |
| hello | 9,644,674 | 9,644,898 | 224 |

通过 `-dumpdep` 检查共享 UI 输入、主题、滚动条及渲染路径的依赖变化；此次增长来自这些公共行为和元数据。
新增 DecimalInput 仍通过构造函数按需引用，没有引入新的第三方 Go 依赖。
本次没有测量应用整体 CPU/GPU 或功耗提升百分比，不承诺所有场景加速。
