# MyGo v0.7.4 → v0.8.3 同步记录

同步日期：2026-11-09。作者：HyperparameterSI。

基线为 MyGo `v0.7.4`（`9bd34e1e`），目标为 `v0.8.3`（`5add3817`）。
按两个标签的最终文件差异迁移核心、平台后端、原生 UI、插件、CLI、TypeScript bridge/runtime、测试及相关文档。
不会引入 v0.8.3 之后 main 分支的提交，也不包含 v0.8.3 已撤回的 Swift 原生模块。
Golet 自身版本暂保持 0.1.15，本次未创建发布标签。

## 已同步的功能与优化

- **渲染与资源生命周期**：增长中的圆角遮罩复用、空闲裁剪遮罩预算、嵌套圆角裁剪、透视变换边界修复；Metal/D3D11 空闲纹理回收；大路径分块绘制；窗口关闭释放内容、位图像素显式释放；Linux OpenGL 批处理减少反射调用、按显示器刷新率调度。
- **原生 UI**：Sticky、滚动吸附、ZIndex、路径裁剪、多层渐变遮罩、命中测试、图片位置、文本平衡/Pretty 换行、描边与阴影、三维位移、匹配过渡与布局动画、退出顺序和减弱动态效果支持。
- **控件与菜单**：三态复选框、受控 Select、TreeBase、菜单/折叠面板基础控件、分隔线双击复位、表头自定义内容、单元格样式、列表保持位置、拖出行；绘制菜单支持副标题、尾部信息、危险操作、自定义内容、稳定 ID、SVG 图标、图标颜色、宽度及打开位置。
- **主题与玻璃**：窗口活动状态、非活动强调色、可配置焦点环、窗口控制样式；系统 Liquid Glass、WebView 上方玻璃面板、合并玻璃组、Tint、Clarity、高对比度与减少透明度适配。保留 Golet 的 Forma 配色、控件材质接口及 Vivid 预设。
- **输入与无障碍**：物理键位置 `InputEvent.Code`、独立修饰键、观察输入、指针类型/压力/取消、焦点离开、按键重复；扩充角色、LabelledBy/DescribedBy、状态、表格索引、选择和键盘提示。macOS 韩文组合输入提交及韩文按词换行修复。
- **编辑器与 Markdown**：代码片段及关联字段、Tab 跳转、长 JSON 行词法与大 diff 布局优化；GitHub alerts、mentions、issue references，链接解析/预览/上下文菜单与图标，流式显示完成状态。
- **国际化与音频**：CLDR 紧凑数字、货币、日期范围、相对时间、SignDisplay；音频保留音高的变速播放。
- **插件**：新增 `plugins/sfsymbols` 与 `plugins/permissionflow` 及独立示例；automation CLI；终端路径识别、AltGr、Alt+点击定位、列宽及同步更新修复；fetch 中止/响应 cookie 过滤、websocket 子协议校验、watch 授权范围修复。
- **平台与窗口**：GPUProcessID、活动/隐藏状态、MIME 文件过滤；Linux 鼠标穿透、GTK CSS、设置门户暗色同步、可配置 DMA-BUF；Windows 托盘菜单和最大化恢复；macOS 窗口拖动、页面滚动恢复。
- **CLI、IPC 与更新器**：资源目录点文件和排除规则、Go workspace 监控、Web 构建原子替换、开发配置重新加载、打包修复；TypeScript 内建名称及符号冲突校验、IPC 通道释放；更新时清除旧文件并恢复权限，文档编辑/恢复的数据保护。

## Golet 适配

- Go 源码仍放在 `backend/`，TypeScript 放在 `frontend/`，Go 模块留在根目录；修正新增生成器测试的模块替换路径。
- 编辑器示例继续与 Studio 共用 `examples/internal/editordemo`；Studio 增加三态复选框和受控下拉演示，并保留已有全控件覆盖检查。
- 新增示例使用 HiddenInset 标题栏与 `{X:12,Y:11}` 窗口按钮位置，预留拖动区域；权限示例关闭后停止刷新定时器。
- 新增自动化命令为 `backend/plugins/automation/cmd/golet-automation`。
- 保留 Golet 结构化日志、异步 Eval 取消测试钩子、Forma、Chromium、Scene3D、BorderBeam 与其他自有扩展。
- 保留上游版权和源码归属；没有导入上游网站 showcase、品牌发布元数据或 Rust 对比基准及结果。上游比较数据不能作为 Golet 的性能结论。

## 应用升级需要留意

1. `HandleInput` / `ObserveInput` 现在接收独立修饰键等事件；未消费的键应返回 false。只处理有逻辑名称的按键时检查 `ev.Key != ui.KeyUnknown`，按物理位置绑定时使用 `ev.Code`。
2. 非活动窗口会使用非活动强调色；需要保留强调色时设置 `theme.InactiveAccent = theme.Accent`。
3. Markdown 默认解析器包含 alerts、mentions 和 references；遍历 AST 时处理新节点。需要旧行为可移除 `parser.ExtMentions | parser.ExtReferences`。
4. 重新生成 TypeScript 客户端后，与 `Record`、`Promise`、`Channel` 等名称冲突的类型会加包名前缀，应用中的类型 import 需要相应调整。
5. `ui.Local` 保存的 `io.Closer` 会在元素移除时关闭；元素与窗口的资源生命周期应保持一致。

## 使用新增示例

```sh
go run ./backend/examples/liquid-glass
go run ./backend/examples/sf-symbols
go run ./backend/examples/permissionflow
go run ./backend/examples/code-editor
go run ./backend/cmd/golet dev backend/examples/studio
go run ./backend/plugins/automation/cmd/golet-automation --help
```

SF Symbols 使用 macOS 系统符号，其他平台提供 fallback；系统玻璃依赖支持它的 macOS，其他渲染器使用 Golet 的 shader 路径。

## 验证

- `CGO_ENABLED=0 go test ./...` 全量通过，包括 Forma、Browser、Scene3D 和玻璃插件。
- macOS、Linux、Windows 的 `go build ./...` / `go vet ./...` 通过；生产配置 `golet_noinspector` vet 通过。
- WASM 全包 vet（排除桌面 CLI），web/GL 测试及 Canvas/韩文按词换行测试通过。
- 前端 252 项测试、TypeScript 类型检查、Prettier 格式检查通过；Studio 绑定生成、生产构建与 SPA 预渲染通过。
- Intel Mac、Windows ARM64、FreeBSD 交叉编译通过；生产配置 `golet_noinspector` 的 core/UI 测试通过。
- 真实 macOS 回归：窗口活动状态、材质跟随活动状态、隐藏窗口首帧、Vibrancy、组合输入 Return、绘制编辑菜单通过。
- Linux/Windows GUI 未在本机执行，仍需 CI 的真实运行环境验证。

同一台机器、同一 Go 1.27.3、`-trimpath -tags golet_noinspector -ldflags '-s -w -X github.com/HyperparameterSI/golet/backend.production=1'`：

| 示例 | 同步前 | 同步后 | 变化 |
| --- | ---: | ---: | ---: |
| counter-native | 9,505,890 B | 9,847,346 B | +341,456 B（3.59%） |
| hello | 9,644,898 B | 9,763,618 B | +118,720 B（1.23%） |

通过 `-dumpdep` 和链接符号检查，普通 WebView 示例没有链接新增可选插件。其增长主要来自共享 darwin 后端的无障碍树与网格角色、玻璃面板/窗口方法、生成器及链接元数据；原生示例还包括新增 UI 布局、输入和渲染路径。此处是二进制体积测量，不是运行内存或 GPU 使用率，不据此宣称倍数性能提升。

## 上游来源

- [完整差异 v0.7.4...v0.8.3](https://github.com/egoist/mygo/compare/v0.7.4...v0.8.3)
- [v0.8.0](https://github.com/egoist/mygo/releases/tag/v0.8.0)
- [v0.8.1](https://github.com/egoist/mygo/releases/tag/v0.8.1)
- [v0.8.3](https://github.com/egoist/mygo/releases/tag/v0.8.3)

### 批次构建基线

Apple M5 / darwin arm64，两个版本分别执行 3 次 `go test -run '^$' -bench '^BenchmarkBuild(Sprites|Boxes)$' -benchmem -count=3 ./backend/internal/gpu`，下表取中位数：

| 场景 | 同步前 | 同步后 | 变化 | 分配 |
| --- | ---: | ---: | ---: | ---: |
| BuildSprites（1872 instances） | 66,622 ns/op | 67,155 ns/op | +0.8% | 两者均 0 allocs/op |
| BuildBoxes（1200 instances） | 86,765 ns/op | 88,409 ns/op | +1.9% | 两者均 0 allocs/op |

这两个基准没有显示通用构建速度提升，也不测实际 GPU 执行时间。纹理/遮罩回收和长文本优化的收益依赖具体场景，不能推导成整个框架更快若干倍。
