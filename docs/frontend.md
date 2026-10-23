# The frontend

Windows that show web pages show the pages of the app's frontend: HTML,
CSS and JavaScript built with any tools. New projects use TypeScript and
[Vite](https://vite.dev), but Golet only needs a URL to load, during
development, and a directory of built files, for builds. Windows of
[native UI](ui/README.md) need no frontend, and an app that only has those has
none.

## How pages load

Windows load pages of the frontend with URLs without a scheme, such as `/`
or `/settings?tab=general`:

```go
golet.NewWindow(golet.WindowOptions{URL: "/"})
win.Page().LoadURL("/settings")
```

They resolve against:

- during `golet dev`, the dev server at `devUrl` in golet.config.ts, which
  `devCommand` starts, e.g. `http://localhost:5173/`, so the dev server's
  hot reload works in the app;
- in builds, `golet://localhost/`, which serves the files of `frontendDist`
  that `golet build` embeds into the executable after running
  `buildCommand`.

```ts
export default defineConfig({
  devUrl: "http://localhost:5173",
  devCommand: "bun run dev:web",
  buildCommand: "bun run build:web",
  frontendDist: "dist",
});
```

Without `devUrl`, `golet dev` serves `frontendDist` from disk instead. Paths
without a file extension that match no file serve `index.html`, so
client-side routers work.

The frontend is a secure context, and its origin is `golet://localhost` on
macOS and Linux and `http://golet.localhost` on Windows, where WebView2
serves custom schemes under `http://<scheme>.localhost`. Storage such as
`localStorage` and IndexedDB belongs to an origin, so data a page stores
under `golet dev` (the dev server's origin) is not what the built app sees.
Keep data you care about in Go, for example in a file in
`App.Path(golet.PathUserData)`, and hand it to pages through bound methods.

Windows can also load other content into their page:

- `win.Page().LoadURL("https://example.com")`, a web page;
- `win.Page().LoadHTML(html, baseURL)`, an HTML string, which reloads load
  again;
- `win.Page().LoadFile("docs/index.html")`, a local file, resolved against the
  working directory, then the executable's directory (and the Resources
  directory of a macOS app).

Pages that are not the app's own cannot call Go methods: see
[who may call](bindings.md#who-may-call).

### Without the CLI

`golet build` embeds the frontend by calling `golet.SetFrontend`, and a
program built with `go build` can do the same with `embed`:

```go
//go:embed all:dist
var dist embed.FS

func main() {
	web, _ := fs.Sub(dist, "dist")
	golet.SetFrontend(web)
	// ...
}
```

## Custom protocols

`golet.Protocol.Handle` serves a URL scheme from an `http.Handler`, for
content your app produces or reads at run time: pages, images, files of
the user. Handlers run on their own goroutines and stream their responses:
writes wait while the webview is behind, so a large file is never held in
memory whole. Pages load such URLs like web URLs, with `fetch`, `<img>`,
ES modules and relative URLs:

```go
// Thumbnails generated into the cache directory.
cache, _ := golet.App.Path(golet.PathCache)
golet.Protocol.Handle("thumbs", golet.FileServer(os.DirFS(cache)))

// Anything an http.Handler does.
golet.Protocol.HandleFunc("api", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"path": r.URL.Path})
})
```

```html
<img src="thumbs://localhost/photo-42.png" />
```

Register schemes before creating the windows that use them.
`golet.FileServer` serves an `fs.FS` with content types and range requests,
falling back to `index.html` like the frontend. Handling the `golet` scheme
replaces the frontend with your handler. Pages of custom schemes are the
app's own: they may call Go methods.

Audio and video are the exception: WebKitGTK's media player does not play
from a custom scheme, so on Linux a `<video>` or `<audio>` of one fails with
`MEDIA_ERR_SRC_NOT_SUPPORTED` (`error.code` 4), whatever the handler
answers. Serve media from a loopback HTTP server instead, the same on every
platform, and give the page its URL:

```go
ln, err := net.Listen("tcp", "127.0.0.1:0") // a free port
if err != nil {
	log.Fatal(err)
}
go http.Serve(ln, golet.FileServer(os.DirFS(videos)))
mediaURL := "http://" + ln.Addr().String() + "/" // <video src="{mediaURL}clip.mp4">
```

Every program of the machine can reach such a server: serve it only what
they may read, or put a secret of the app's in its paths.

Each scheme is its own origin, so a page fetching another scheme makes a
cross-origin request: the frontend reading `api://localhost/`, or, during
`golet dev`, the dev server reading `thumbs://localhost/`. Responses to the
app's own pages, those that may call Go methods, carry the CORS headers
that let them read the response, and Golet answers their preflights.
Other sites get none, unless the window trusts them
(`PageOptions.TrustedOrigins`), and handlers may set their own.
On Windows, fetch a scheme at `http://<scheme>.localhost/`, where
WebView2 serves it.

## The golet-runtime package

Golet injects its runtime into every page as `window.golet`. The
`golet-runtime` npm package, which new projects depend on, gives it types
and imports; the generated client is built on it.

```ts
import { currentWindow, isGolet, onFileDrop, runtime } from "golet-runtime";

if (isGolet() && runtime().platform === "darwin") {
  document.documentElement.classList.add("mac");
}

await currentWindow.toggleMaximize();
```

| Export | What it does |
|---|---|
| `call(method, ...args)`, `on(event, listener)`, `once(event, listener)`, `event<T>(name)` | untyped calls and events, see [bindings](bindings.md#without-the-generated-client) |
| `Channel` | a stream of values from a Go method, see [channels](bindings.md#channels) |
| `isCallError(err)` | whether a rejection came from a Go error |
| `isGolet()` | whether the page runs in a Golet window, not in a browser tab of the dev server |
| `runtime()` | the runtime: `platform` (`"darwin"`, `"linux"` or `"win32"`), `windowId` (the Go window's `ID()`) and `version`; throws outside Golet |
| `currentWindow` | the page's window: `minimize`, `maximize`, `unmaximize`, `toggleMaximize`, `isMaximized`, `toggleFullScreen`, `close` and `setTitle` |
| `onFileDrop(listener)` | files dropped on the window, with their paths |
| `preferences()` | the accent color and the accessibility settings of the desktop; throws outside Golet |
| `onPreferences(listener)` | the changes of `preferences()` |
| `locale()`, `direction()`, `onLocaleChange(listener)` | the language the app shows (`App.Locale`), such as `"zh-Hans"`, its direction, `"ltr"` or `"rtl"`, and its changes; the page has them from document start, as `golet.locale` and `golet.dir` |

`golet-runtime/i18n` formats messages of ICU MessageFormat, as package
`i18n` does in Go. The module `golet generate` writes for the app's
[catalogs](i18n.md#in-pages), `locales.ts`, is built on it, and holds the
messages of every language, which the frontend's build bundles: a page
says its first words in the app's language without waiting for a request.

Opened in a regular browser, which is handy for working on the layout,
pages have no runtime: calls reject and `isGolet()` is false.

## Custom title bars

A frameless window has no title bar or borders, and the page draws its own.
Mark the elements that drag the window with the CSS property
`--app-region: drag`:

```go
golet.NewWindow(golet.WindowOptions{URL: "/", Frameless: true})
```

```css
.titlebar {
  --app-region: drag;
  height: 40px;
}
.titlebar .search {
  --app-region: no-drag;
}
```

```ts
import { currentWindow } from "golet-runtime";

minimizeButton.onclick = () => currentWindow.minimize();
maximizeButton.onclick = () => currentWindow.toggleMaximize();
closeButton.onclick = () => currentWindow.close();
```

The property is inherited: `no-drag` excludes parts of a region again, and
buttons, links, inputs and other interactive elements inside a region stay
clickable. Double-clicking a region does what double-clicking a title bar
does.

Frameless windows still resize from their edges. On Linux, where GTK gives
windows without decorations no borders, the outermost 5 px of the page act
as the borders: they show resize cursors, and the page gets no clicks or
pointer moves there.

On Linux a frameless window also keeps the menu bar that `App.SetMenu` or
`Window.SetMenu` gives it, above the page and so above a custom title bar.
Give such a window an empty menu of its own,
`win.SetMenu(golet.NewMenu(nil))`, and show menus from the page instead
(`Menu.PopupAt`, from a bound method the page calls).

### Hidden title bars

A window can keep its window controls and only hide the title bar, so the
page extends under them and draws the rest of the title bar: set
`TitleBarStyle` to `golet.TitleBarHidden`. The controls stay at the top of
the window:

- On macOS, the traffic lights. `golet.TitleBarHiddenInset` leaves more
  room around them, and `TrafficLightPosition` moves them: the position of
  the close button's top-left corner in the window.
- On Windows, minimize, maximize or restore, and close at the top-right
  corner, drawn as Windows 11 draws its own. They fill the height of the
  title bar, `TitleBarHeight` (32 by default), and hovering maximize opens
  the snap layouts. The window resizes from its top edge.
- On Linux, GTK's own title buttons, where the desktop's button layout puts
  them, centered in `TitleBarHeight` (by default the height of the
  desktop's header bars). The layout decides which buttons show, possibly
  none, as tiling window managers are often set up. A Wayland compositor
  that decorates windows itself, as KWin, Hyprland and Sway do, has none:
  its own title bar holds the buttons, and a tiling one shows none. The
  outer pixels of the page resize the window, as a frameless one's.

```go
golet.NewWindow(golet.WindowOptions{
	URL:            "/",
	TitleBarStyle:  golet.TitleBarHidden,
	TitleBarHeight: 40,
})
```

From its first paint, the page learns the room the controls take from
CSS variables on `:root`:

| variable | |
|---|---|
| `--golet-titlebar-height` | the height of the title bar |
| `--golet-titlebar-inset-left`, `--golet-titlebar-inset-right` | the room the controls take from the left and from the right edge |

They are `0px` in full screen, where the controls hide (on macOS they come
back with the menu bar). The page's title bar keeps clear of the controls
and, marked with `--app-region: drag`, moves the window:

```css
.titlebar {
  --app-region: drag;
  height: var(--golet-titlebar-height, 40px);
  padding-left: var(--golet-titlebar-inset-left, 0px);
  padding-right: var(--golet-titlebar-inset-right, 0px);
}
```

`Transparent: true` lets a page with a transparent background show the
desktop through, and `Vibrancy` puts a blurred material behind it on macOS
and Windows 11 22H2 and later:

```go
golet.NewWindow(golet.WindowOptions{
	URL:           "/",
	TitleBarStyle: golet.TitleBarHidden,
	Transparent:   true,
	Vibrancy:      golet.VibrancySidebar, // Mica on Windows 11
})
```

See `examples/frameless` for a complete custom title bar. Native UI shows
materials the same way ([Vibrancy](ui/windows.md#vibrancy)), as
`examples/vibrancy` does in a translucent sidebar under an inset title bar.

## Dropped files

Files dragged from Finder or Explorer and dropped on the page reach the
page with their paths, which DOM drop events do not tell:

```ts
import { onFileDrop } from "golet-runtime";

onFileDrop(({ paths, x, y }) => {
  const target = document.elementFromPoint(x, y);
  // ...
});
```

and Go, with `win.OnFileDrop`. The page's own drag and drop keeps working:
its drop listeners still get the `File` objects, and dropping files where
the page does not handle them no longer replaces the page with the file.

## Dark mode

Pages follow the system appearance through the `prefers-color-scheme` media
query. Declare `color-scheme: light dark` in CSS so that form controls and
scrollbars follow too, and give windows a `BackgroundColor` matching the
page in both appearances, which the window shows until the page paints:

```css
:root {
  color-scheme: light dark;
}
```

```go
golet.NewWindow(golet.WindowOptions{URL: "/", BackgroundColor: "light-dark(#f5f5f7, #1e1e1e)"})
```

`golet.Theme.SetSource(golet.ThemeDark)` forces an appearance, which pages see
through the same media query.

## Accent color and accessibility

`prefers-reduced-motion` and `prefers-contrast` follow the desktop on their
own. The accent color the desktop uses and the factor it asks text to grow
by are what the media queries do not tell a page: they are CSS variables on
`:root` from the first paint.

| variable | |
|---|---|
| `--golet-accent` | the accent as `#rrggbb`; unset where the desktop has none, so `var(--golet-accent, AccentColor)` falls back to the desktop's own color |
| `--golet-text-scale` | the factor text grows by; `1` is the usual size, so `calc(1rem * var(--golet-text-scale, 1))` |

```css
:root {
  --accent: var(--golet-accent, AccentColor);
  font-size: calc(1rem * var(--golet-text-scale, 1));
}
```

`preferences()` of `golet-runtime` returns the settings themselves, current
from document start on, and `onPreferences` tells their changes, which only
the desktop makes; a component reads the first and subscribes to the
second, as with `matchMedia`:

```ts
import { onPreferences, preferences, type ThemePreferences } from "golet-runtime";

const apply = ({ reduceTransparency }: ThemePreferences) =>
  document.documentElement.classList.toggle("opaque", reduceTransparency);
apply(preferences());
const off = onPreferences(apply);
```

| field | |
|---|---|
| `accent` | the accent as `"#rrggbb"`, or `null` where the desktop has none |
| `reduceMotion` | the desktop asks for less motion, as `prefers-reduced-motion` |
| `reduceTransparency` | the desktop asks for opaque surfaces: macOS's Reduce Transparency, Windows's transparency effects turned off; false on Linux, which has no such setting |
| `highContrast` | the desktop asks for more contrast, as `prefers-contrast: more` |
| `textScale` | how many times larger than usual text should be; `1` is usual |

A page that navigates starts with the values of the moment: the script
that tells them at document start is replaced when they change.
`golet.Theme.Preferences()` reads them on the Go side, see
[native](native.md#accent-color-and-accessibility).

## Links and new windows

Links with `target="_blank"` and `window.open()` call the window's open
handler. Without one, `http(s)` URLs open in the default browser and
anything else is denied; `win.Page().SetWindowOpenHandler` decides otherwise:

```go
win.Page().SetWindowOpenHandler(func(req golet.WindowOpenRequest) *golet.WindowOptions {
	if strings.HasPrefix(req.URL, "https://docs.example.com/") {
		return &golet.WindowOptions{Width: 900, Height: 700} // a window of the app
	}
	go golet.Shell.OpenExternal(req.URL) // the default browser
	return nil
})
```

`win.Page().OnWillNavigate` can cancel the page's own navigations, e.g. to keep the
app's window on the app: see [Windows](windows.md#pages).

## Preload scripts

`PageOptions.PreloadScript` runs JavaScript in every page of the window
before the page's own scripts, once `window.golet` exists.

## The web inspector

Development builds have the web inspector: right-click a page and choose
Inspect Element (Inspect on Windows), or call `win.Page().OpenDevTools()`.
Production builds do not, unless built with `golet build -debug` or a window's
page sets `DevTools: golet.DevToolsEnabled` in `WindowOptions.Page`.

On Windows, WebView2's own keys come with the inspector: F5 and Ctrl+R reload
the page (a reload [`OnWillNavigate`](windows.md#pages) is not asked about),
Ctrl+P prints, Ctrl+F finds and Ctrl+Plus zooms only where the inspector is
on. Without it, as in production builds, a page has none of them, as on macOS
and Linux, and only the app's menu items take these keys, such as the View
menu's (`RoleViewMenu`).
