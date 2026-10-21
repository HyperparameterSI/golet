# Go backend and native UI

The repository root owns `go.mod`. All Go implementation lives here.

```go
import "github.com/HyperparameterSI/golet/backend"
import "github.com/HyperparameterSI/golet/backend/ui"
```

The first import declares package `golet`, so existing `golet.NewWindow` calls remain unchanged after updating import paths.

From the repository root:

```sh
CGO_ENABLED=0 go test ./backend/...
go run ./backend/examples/counter-native
go run ./backend/cmd/golet dev backend/examples/studio
```

`internal` contains platform backends; `ui` contains the native Go toolkit; `plugins` contains Go plugin implementations; `examples` contains Go entrypoints. TypeScript implementations are in the sibling `frontend` directory. CLI scaffold payloads and generated embedded assets are intentional exceptions: they must ship inside the CLI/framework binaries.
