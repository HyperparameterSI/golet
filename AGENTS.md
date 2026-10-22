# AGENTS.md

Golet is a desktop application framework in Go. Windows show web pages in the
system webview (WKWebView on macOS, WebKitGTK on Linux, WebView2 on Windows),
with typed Go ↔ TypeScript IPC, or native UI that Golet draws itself (package
`ui`).

**Read [docs/architecture.md](docs/architecture.md) before changing code.** It
describes the layers, the threading model, native interop without cgo, the IPC
wire protocol, TypeScript generation, and the checklist for adding features.

## Hard rules

- **No cgo.** Everything must build with `CGO_ENABLED=0`. Call native code
  through purego (`internal/darwin`, `internal/linux`) or `syscall`
  (`internal/windows`), never `import "C"`.
- **Bun is dev tooling only.** The repository is a Bun workspace
  (`frontend/packages/bridge`, `frontend/packages/runtime`, `frontend/packages/cli` with its platform
  packages, the plugins' packages in `frontend/plugins/`, `frontend/examples/*`, and `frontend/website`); Bun also installs and runs the template's scripts
  (Vite, the golet-cli package). Nothing in an app may need Bun at run time.
- **Scope:** system webview on macOS, Linux and Windows (WebView2), and
  native UI that Golet draws itself (`WindowOptions.Content`, package `ui`;
  see "Native UI" in the architecture guide), which also runs in browsers
  (`GOOS=js GOARCH=wasm`, `internal/web`, drawn with WebGL 2 by the GL
  renderer; see "Browsers"). Optional CEF is isolated in `backend/plugins/chromium` and currently
  provides a macOS standalone browser host and an experimental Windows host, not a default engine or embedded
  `ui` view. Keep its runtime and build tooling out of ordinary applications. Other platforms must keep compiling through
  `internal/unsupported`.
- **Great DX over Electron parity.** Keep the familiar feel (app lifecycle,
  windows, menus) but design Go-first APIs: typed IPC via `golet.Bind` and
  `golet.NewEvent[T]` with the generated client, not string channels.
- **Threading.** Platform backend methods and event callbacks run on the main thread.
  Bound IPC service methods run on their own goroutines.
  Public methods must be safe from any goroutine (`onMain`/`postMain`/`await`
  in `loop.go`). Never block the main thread waiting on work that needs it.
- **purego callbacks are a scarce, never-freed resource.** Create them once
  per signature at startup and route by user data, never per window or call.
- **Apps pay only for what they use.** A feature must not grow apps that
  don't use it: `examples/counter-native` (native UI) and `examples/hello`
  (webview) are measured on every push. What every app reaches (`App.Run`,
  windows, the default menus, `init` functions, the backends' methods)
  must not call into opt-in features: reach them through functions their
  constructor registers (as `linkedDocuments`), and keep heavy
  dependencies and data out of those paths. Never call
  `reflect.Type.Method`, `Value.Method` or `MethodByName` there: it makes
  the linker keep every exported method of the program (2.5 MB). Measure
  both apps before and after a change to `golet`, `ui`, a backend or the
  dependencies; explain any growth of an app that doesn't use the feature
  in the PR, and find its cause with `-dumpdep` ("Binary size" in the
  architecture guide).

## Commands

```sh
go build ./... && GOOS=linux go build ./... && GOOS=windows go build ./...  # all must pass
go vet ./... && GOOS=linux go vet ./... && GOOS=windows go vet ./...
go vet -tags golet_noinspector ./...        # production builds, without the inspector of native UI
GOOS=js GOARCH=wasm go vet $(GOOS=js GOARCH=wasm go list ./... | grep -v /cmd/golet)  # browsers, all but the CLI
go test ./...                              # core (fake backend), tsgen, CLI, accelerator
go build -trimpath -tags golet_noinspector -ldflags "-s -w -X github.com/HyperparameterSI/golet/backend.production=1" -o /tmp/counter-native ./backend/examples/counter-native  # an app's release size, as the benchmarks build it
GOLET_E2E=1 go test ./backend/internal/e2e          # real GUI tests (macOS desktop session)
bun install && bun run test && bun run typecheck   # the workspace, from the root
bun run build                              # build internal/bridge/bridge.js and the npm packages' dist/
bun run --cwd frontend/packages/cli binaries        # cross-compile the CLI into the golet-cli platform packages
go run ./backend/cmd/golet generate backend/examples/todo   # regenerate an example's TypeScript client
go run ./backend/cmd/golet dev backend/examples/todo        # live reload (dev bundle in examples/todo/.golet)
go run ./backend/cmd/golet dev -web backend/examples/counter-native   # native UI in the browser, with live reload
go run ./backend/cmd/golet build backend/examples/todo      # .app + .dmg in examples/todo/build
go run ./backend/cmd/golet build -platform web backend/examples/counter-native   # a static site in examples/counter-native/dist/web
go run ./backend/examples/gallery                  # the native UI toolkit's tour
go run ./backend/examples/code-editor              # the code editor and the diff viewer (ui/editor)
go run ./backend/examples/markdown-chat            # Markdown replies streaming in (ui/markdown)
go run ./backend/examples/notes                    # the rich text editor (ui/richtext), in a notes app
go run ./backend/examples/terminal                 # a terminal (downloads libghostty-vt once)
go run ./backend/examples/loading                  # loading.dev's 29 indicators (plugins/loading)
go generate ./backend/plugins/terminal             # build libghostty-vt for every platform (with Zig), write Ghostty's themes
go generate ./backend/internal/gpu/d3d11           # recompile the Direct3D shader (on Windows)
go generate ./backend/internal/gpu/metal           # recompile the Metal shader (on macOS, with Xcode)
go generate ./backend/internal/l10n                # Golet's own strings: their functions, and the packages of translations in locales/
go generate ./backend/plugins/glass                # recompile the glass plugin's shaders (on macOS, and on Windows)
bun run --cwd frontend/website dev                  # the website, with docs/ at /docs (see website/README.md)
bun scripts/bench.ts run [--e2e]           # benchmarks and app sizes, medians of 6 runs (CI: bench.yml, /benchmarks)
go run ./backend/internal/idlemem app...           # apps' memory once idle, with their webviews' processes
go run ./backend/internal/inputlatency             # native UI's input-to-screen latency on macOS (opens windows)
```

Linux GUI tests cross-compile and run in a container with WebKitGTK and Xvfb;
Windows GUI tests need Windows with the WebView2 Runtime (a GitHub Actions
`windows-latest` runner has it). See "Testing" in the architecture guide.

## Conventions

- Never write `[skip ci]` (or `[ci skip]`, `[no ci]`) in a commit message,
  not even quoted in its body: GitHub skips the push's workflows, CI and
  benchmarks included, and Cloudflare the website's deploy. Only the
  benchmarks branch's commits say it.
- New behavior goes into `package golet` first; backends only translate
  `internal/platform` calls to native ones. Implement every platform method in
  `darwin`, `linux`, `windows`, `web`, `fake` and `unsupported`.
- Objective-C: `alloc`/`init` objects are owned and must be released; wrap
  temporary objects in `withPool`; copy blocks you call later.
- Official plugins have Go implementations in `backend/plugins/<name>` and npm
  clients (`@golet-plugins/<name>`) in `frontend/plugins/<name>`, released with
  golet-runtime's version. The updater plugin is Go only; its `native`
  package draws the update window in native UI. The terminal plugin is Go
  only too, for native UI, with libghostty-vt loaded through purego: its
  `golet-plugin.json` pins the builds the CLI puts into apps, published as
  release assets (`libghostty-vt-<commit>`); change it with `go generate`. The glass
  plugin is Go only too, an effect of the renderers (`scene.Effect`): its
  shader in Metal, HLSL and GLSL with a CPU twin that must draw the same
  pixels, compiled ahead of time by `go generate` on macOS and Windows.
  The game plugin is Go only too, what games need beyond what apps do,
  opt-in, kept out of core `ui`: a game loop and a camera for native UI,
  and `plugins/game/gamepad`, which reads controllers. The loading plugin
  is Go only too, loading.dev's indicators drawn with `ui.Painter` on
  ui's public API: keep them as loading-dev draws them in a browser.
- The `dist/` of npm packages is not committed: run `bun run build` after
  `bun install` (CI and releases do). Commit `internal/bridge/bridge.js`
  whenever its sources change, and keep the generated clients of examples
  (`frontend/examples/*/src/golet.ts`) up to date. Built frontends are not
  committed: `golet build` embeds them.
- `golet-runtime` (`frontend/packages/runtime`) and `golet-cli` (`frontend/packages/cli`, with
  its per-platform packages in `frontend/packages/cli/npm`) are released to npm with
  the same version as the Go module (`golet.Version`, also the CLI's
  `version`), which `bun scripts/version.ts <version>` sets everywhere; the
  template depends on both. Pushing a `v<version>` tag releases everything
  (`.github/workflows/release.yml`, see "Releasing" in the architecture
  guide); never commit the CLI's binaries.
- The configuration types of `frontend/packages/cli/index.d.ts` (for `golet.config.ts`)
  mirror `Config` in `backend/cmd/golet/config.go`; `TestConfigTypes` checks them.
- Tests: unit tests through `internal/fake`; GUI behavior in `internal/e2e`.
- Tests never pass or fail on wall-clock time ("Waiting in tests" in the
  architecture guide): wait for the condition or event with
  `internal/wait`, whose one-minute guard only turns a hang into a
  failure and is never tuned; advance a `ui.Tester`'s clock instead of
  sleeping; measure work by counts, not durations. A test that can time
  out while the code is right has a race to fix, in the test or the
  product, not a timeout to raise or a rerun to hide behind.

## Repository layout

All Go implementation lives in `backend/`; all editable web applications and TypeScript packages live in `frontend/`. `go.mod` stays at the repository root. The public Go package is `github.com/HyperparameterSI/golet/backend`. Embedded CLI template payloads and generated bridge.js remain with Go so a published CLI needs no source checkout. The default shadcn template generates separate backend/ and frontend/ source directories and uses TanStack Start SPA mode; package only frontend/dist/client, never a Node server.
