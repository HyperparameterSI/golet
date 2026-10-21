# Studio frontend

Author: HyperparameterSI

TanStack Start + React + shadcn 组件目录。包含分类筛选、API 搜索、原生演示入口、Browser 控制面板和 Go IPC 交互。

从仓库根目录运行 `pnpm dev` 启动完整桌面示例。单独运行 `npm --prefix frontend/examples/studio run dev:web` 只有前端，没有 Go 桥，不能打开原生控件。构建完整应用使用 `pnpm build:app`。

组件目录来自 Go 的 `Showcase.Catalog`，由 `golet generate` 生成类型化绑定。修改服务后运行：

```sh
go run ./backend/cmd/golet generate backend/examples/studio
```
