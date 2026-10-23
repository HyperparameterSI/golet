# MyGo v0.6.3 integration

<!-- Author: HyperparameterSI -->

Golet integrates the framework changes through MyGo **v0.6.3**, commit
`15242900157b77228d35ca071f43e94a11bc90a0`. This is the latest stable tag
verified on 2026-10-23. Upstream sources retain their MIT license; see the
repository license and [upstream release](https://github.com/egoist/mygo/releases/tag/v0.6.3).

The comparison includes changes since upstream `5fb6009`, the source baseline
matching Golet's existing terminal and UI implementation, rather than only the
v0.6.2 → v0.6.3 patch. Existing Golet changes are merged, not replaced.
Golet's version, release configuration, package names, frontend/backend layout,
Chromium plugin, structured diagnostics and native themes remain independent.

## Integrated features

- **Content rules:** `CompileContentRules`, `PageOptions.ContentRules` and
  `Page.SetContentRules`; system WebViews can block resources, hide matching
  elements, block cookies and upgrade HTTP requests. Rules use WebKit's JSON
  format, are compiled/reused, and remain opt-in for linking. macOS and Linux
  use WebKit; Windows uses WebView2 request interception. Browser/wasm builds
  report unsupported. Windows cannot intercept WebSockets through this API.
  This is not an automatic implementation for the separate CEF plugin.
- **WebView PDF text:** correct Unicode text mapping in page PDF output,
  restoring Chinese characters that macOS WebKit represents as radicals so
  copying and searching use the expected characters.
- **Input zoom:** touchpad pinch events reach native UI through `InputZoom`,
  with desktop and browser platform implementations and test hooks.
- **Window controls:** `HideWindowButtons`, custom `MaximizeButton` elements
  and Windows 11 snap-layout support, linked only when used.
- **Terminal:** scrollback search, Kitty graphics images, image mirroring,
  program colors and light/dark reporting, minimum text contrast, hovered-link
  feedback, remote-session modes and clipboard forwarding, safe shell-specific
  file-drop quoting, clear-screen actions and Windows ConPTY synchronization.
- **Native menus:** Windows 11 menu drawing, padding and shadows, and native
  UI menu shortcut handling on Windows/Linux.
- **Packaging:** `linux.webview: false` avoids WebKitGTK installation
  dependencies for native-only applications; NSIS scripts use UTF-8 BOM for
  non-ASCII application names and paths.

## Integrated fixes and performance work

- Direct3D paces unchanged animation frames; Metal presentation handles
  drawable/layout transitions; OpenGL swap interval and probe handling improve.
- Large text layouts may be cached alongside recently used layouts within the
  existing 64 MiB live budget (at most 32 MiB per large layout); inline transforms
  are computed once per paint call. This trades bounded cache memory for less
  repeated shaping, rather than promising lower memory in every workload.
- Metal atlas sampling explicitly selects mip level zero; matching compiled
  renderer and glass shader assets are included.
- Scene damage, folding, path flattening, list selection, diff layout, hover
  after layout, scroll chaining and scroll cursor handling are updated.
- Hidden scrollbars no longer intercept content clicks; a flexible final table
  column has no ineffective resize handle; Windows text fields consume handled
  shortcuts so menu Paste cannot paste twice.
- Document equality callbacks run outside document locks; saving, watching and
  file replacement preserve correct state and permissions.
- WebView2 registration tokens and modal task-dialog outputs have safe pointer
  lifetimes; maximized material windows avoid duplicate caption buttons.
- Single-instance activation, notifications, Chinese/Japanese font fallback,
  tray/taskbar behavior, browser accelerators and accessibility fixes are included.

## Golet-specific integration

The upstream toolbar regression check allows a subpixel floating-point rounding
error because Golet's theme spacing is fractional. It still rejects actual
extra layout space. New Go files include the requested author annotation.
Upstream issue references remain links to upstream, not to nonexistent Golet issues.

The existing theme/material code and Chromium host implementation are retained.
Their APIs are independent of the system-WebView content rules above.

## Validation

Run from the repository root with `CGO_ENABLED=0`:

```sh
go test ./...
go build ./...
GOOS=linux go build ./...
GOOS=windows go build ./...
go vet ./...
GOOS=linux go vet ./...
GOOS=windows go vet ./...
go vet -tags golet_noinspector ./...
GOLET_E2E=1 go test ./backend/internal/e2e -run 'Test(ContentRules|PrintToPDFText|PinchZooms)$'
```

Windows/Linux GUI tests require their respective desktop environments; a
cross-build does not establish native runtime behavior. The upstream regression
tests are included for those platforms' CI. Frontend tests and type checking
remain part of the workspace checks. No Golet version/tag is changed by this
integration alone.

### Local verification results

On macOS arm64, using production flags (`-trimpath -tags golet_noinspector`,
`-ldflags '-s -w -X github.com/HyperparameterSI/golet/backend.production=1'`):

| Example | Before | After | Difference |
| --- | ---: | ---: | ---: |
| counter-native | 9,232,690 B | 9,266,226 B | +33,536 B (0.36%) |
| hello | 9,354,274 B | 9,422,258 B | +67,984 B (0.73%) |

Linker `-dumpdep` inspection finds no content-rule compiler/matcher implementation
in either example and no PDF parser in counter-native. The hello example already
uses `Bind`, whose reflective method discovery retains exported Page methods,
including `PrintToPDF`; its new PDF text-correction parser consequently contributes
to hello's increase even without calling that method. The integration preserves
that existing API behavior; it does not claim zero binary growth or measured
reductions in runtime CPU/RSS. Core input/rendering fixes also change ordinary
native applications.

The integration also fixes test infrastructure inherited from the directory
split: plugin GUI tests now read built assets under `frontend/`, so building the
workspace actually enables them. macOS menu test actions focus their requested
window; a two-document regression test verifies the right file is saved. The
language screenshot test waits for an actual new frame and reads its screenshot
and generation on the main thread, avoiding fixed delays and stale captures.

Final checks passed locally: all Go tests; production-tag core/UI tests;
macOS/Linux/Windows builds and vet; amd64 cross-builds for all three platforms;
Windows/Linux 386 and FreeBSD builds; WebAssembly vet; 204 frontend tests;
frontend type checking and package builds; and the full macOS GUI suite
(40.262 seconds, platform-specific/unsupported checks skipped as declared by
the tests). The GUI suite includes the built Fetch/WebSocket, SQLite and watch
plugin clients. A further drag-and-drop test race was fixed by recording the
frame count before invalidation, not after the requested frame may have run.
Windows and Linux GUI execution remains for their native CI runners.
