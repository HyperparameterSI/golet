# The golet CLI

`golet` creates, develops, packages and publishes Golet apps.

## Install

Projects with a web frontend depend on the CLI as the `golet-cli` npm
package, which holds a prebuilt binary for macOS, Linux and Windows on
arm64 and x64, and run it from their scripts:

```sh
bun add -d golet-cli    # or: npm install -D golet-cli
bun run dev            # "dev": "golet dev" in package.json
bunx golet doctor       # any command, in the project
```

Outside a project, run `bunx golet-cli` or `npx golet-cli`, e.g.
`bunx golet-cli init my-app`: npm has an unrelated package named `golet`, so
`bunx golet` would run it where golet-cli is not installed. `npm install -g
golet-cli` puts `golet` on your `PATH`. `GOLET_CLI_BINARY` makes the package
run another build of the CLI.

Projects of [native UI](ui/README.md) have no npm packages: their `go.mod` has
the CLI as a [tool](https://go.dev/doc/modules/managing-dependencies#tools),
which `go tool golet` runs at the version it pins.

The CLI is a Go program, which Go runs too, without installing it:

```sh
go run github.com/HyperparameterSI/golet/backend/cmd/golet@latest init my-app
```

`go install` puts `golet` on your `PATH`:

```sh
go install github.com/HyperparameterSI/golet/backend/cmd/golet@latest
```

Either way, building apps needs Go, which compiles them.

The commands that take a project directory, `[dir]`, default to the
current directory.

## golet init

```
golet init [flags] <dir>
```

Creates a project in a new or empty directory: a Go module and a
TypeScript frontend built with Vite, side by side, a default icon in
`resources/icon.png`, and a package.json with the scripts `dev`, `build`
and `generate`. It installs the dependencies with Bun and generates the
TypeScript client. See [a web frontend](getting-started.md#a-web-frontend).

With `-template native` the project is a Go module alone, whose window
shows [native UI](ui/README.md): `main.go`, a test of its view, `golet.json` and the
icon. The module has the CLI as a [tool](https://go.dev/doc/modules/managing-dependencies#tools),
so `go tool golet dev` and `go tool golet build` run the version it pins,
without Bun. See [native UI](getting-started.md#native-ui).

Both templates include the `golet-maintenance` agent skill in
`.agents/skills/golet-maintenance/`, with guidance on element lifetimes,
stable identity, threading, and testing when maintaining the app.

| Flag | |
|---|---|
| `-template` | `web`, a TypeScript frontend (the default), or `native`, native UI in Go |
| `-name` | the app's name (default: the directory's name) |
| `-module` | the Go module path (default: the directory's name) |
| `-golet` | a checkout of Golet to use, through a `replace` directive, instead of the released module; the scripts then run the checkout's CLI with `go run`; run `bun install && bun run build` in the checkout first, for the web template |

## golet install-skills

```sh
golet install-skills [dir]
```

Installs or updates the agent skills bundled with the CLI version you run
in an existing project's `.agents/skills/`. The directory defaults to the
current directory; no project configuration, Go build, or frontend tools
are needed. New projects receive these skills through `golet init` too.

The command overwrites bundled files such as
`golet-maintenance/SKILL.md` and `golet-maintenance/agents/openai.yaml`,
including local edits to those files. Other skills and extra custom files
are preserved. Upgrade the CLI, then run the command again to refresh its
guidance.

## golet dev

```
golet dev [flags] [dir]
```

Develops the app with live reload. It writes the TypeScript client, runs
`devCommand` from the configuration, such as a Vite dev server, waits for
`devUrl` to answer, then builds a development app, which loads `devUrl` in
place of its built frontend, and starts it. Without `devUrl` the app serves
`frontendDist` from disk.

Changes to the Go code, the configuration, the icon, the resources or the
[catalogs](i18n.md) rebuild the app, regenerate the TypeScript client and
the functions of the catalogs, and restart the app: the running
build quits, then the new one starts, so the two never hold the same files,
locks or web view profile at once. A build that fails to compile keeps the
running one. Frontend changes are the dev server's to handle. Quitting the app, or
Ctrl+C, ends golet dev (press Ctrl+C again not to wait for the app to quit),
and `App.Relaunch` restarts the app.

On a terminal, each step shows a spinner with what it does, such as the Go
packages being compiled, then how long it took; what the dev server and the
app print is labeled `web │` and `app │`, and a rebuild names the files that
changed. Type a key, then Enter, as in Vite: `r` rebuilds and restarts the
app, `c` clears the console, `q` quits and `h` lists them. `NO_COLOR` turns
colors off, and output to a file or a pipe has a line per step.

The development app is named `<name> Dev`, with the identifier
`<identifier>.dev`, so that its data, preferences and single instance lock
stay apart from the installed app's. On macOS it is a real app bundle, in
`.golet/dev`; on Windows its executable carries the icon, manifest and
version information that `golet build` embeds.

| Flag | |
|---|---|
With `-web`, golet dev develops an app of [native UI in browsers](ui/browsers.md)
instead: it compiles the app to WebAssembly, serves it at
`http://localhost:5174` (or the next free port), opens it, and rebuilds it as
its Go code changes, which reloads the page. A build that fails to compile
leaves the page running the last one, and shows why over it until a fix
compiles. `o` opens the page again.

| Flag | |
|---|---|
| `-skip-dev-command` | does not run `devCommand`, e.g. when the dev server already runs |
| `-sign` | the identity that signs the development app on macOS (default: `-`, ad hoc) |
| `-web` | develops the app in browsers |
| `-port` | with `-web`, the port to serve the app at, or the next free one (default: 5174) |
| `-open` | with `-web`, opens the app in the browser (default: true) |
| `-lang` | runs the app in a language, such as `de` or `zh-Hant`, or a pseudo locale, `en-XA` or `ar-XB`, without changing the language it keeps; see [trying a language](i18n.md#trying-a-language) |

## golet build

```
golet build [flags] [dir]
```

Builds production apps: runs `buildCommand`, compiles the app with
`frontendDist` embedded, and packages it for each platform in
`<out>/<os>-<arch>`. See [Building and distributing](distribution.md).

The platform `web` (or `js/wasm`) builds an app of native UI for browsers
into `<out>/web`, a static site to serve from anywhere: `index.html`, which
loads the app, Go's `wasm_exec.js`, the app's WebAssembly, named by its
content so that it can be cached for good, and the icon. See
[Browsers](ui/browsers.md).

| Flag | |
|---|---|
| `-platform` | comma separated `GOOS/GOARCH` targets, e.g. `darwin/universal,windows/amd64,linux/amd64`, or `web` (default: this machine's) |
| `-debug` | keeps development features, such as the web inspector and the [inspector of native UI](ui/inspector.md) |
| `-o` | the output directory (default: `out` of the configuration) |
| `-sign` | the macOS signing identity (default: `macos.signingIdentity`) |
| `-skip-build-command` | does not run `buildCommand` |
| `-skip-dmg` | does not make macOS disk images |
| `-skip-notarize` | does not notarize, even with `macos.notarize` set |
| `-upload` | uploads the installers and updates to a draft GitHub release of the version (`updates.github`), with `gh`, or to the bucket of `updates.s3` |

## golet generate

```
golet generate [-o file] [dir]
```

Writes the typed TypeScript client for the services bound with `golet.Bind`
and the events declared with `golet.NewEvent`, to `bindings` of the
configuration unless `-o` names another file. `golet dev` and `golet build` run
it for you; apps without a frontend get no client. See
[the generated client](bindings.md#the-generated-client).

With `i18n` in the configuration, it first checks the app's
[catalogs](i18n.md) and writes their functions: `catalog_gen.go` in the
catalogs' directory, and a TypeScript module next to the client, named
after the directory (`src/locales.ts`). Errors in the catalogs stop it.

## golet i18n check

```
golet i18n check [dir]
```

Checks the app's [catalogs](i18n.md#checking-the-catalogs): it lists their
errors, such as messages that do not parse or arguments a translation adds,
and their warnings, such as missing translations, and exits with status 1
when there are errors.

## golet keygen

```
golet keygen [flags]
```

Creates the key pair that signs [updates](updates.md): `golet-update.key`,
the secret, and `golet-update.pub`, and prints the `updates` section to add
to the configuration.

| Flag | |
|---|---|
| `-o` | the directory to write the keys to (default: `golet/update-keys` in the user's configuration directory) |
| `-force` | replaces existing keys |

## golet doctor

Checks the development machine: Go, Bun and the webview of the platform,
which projects with a web frontend need (on Linux GTK 3 too, and not
WebKitGTK in a project whose `linux.webview` is `false`), and the tools of optional
features: NSIS for Windows installers (which `golet build` downloads on
Windows), `signtool` or `osslsigncode` for signing Windows apps, `gh` for
`-upload`, and Developer ID identities for signing macOS apps.

## golet version

Prints the version of the CLI.

## Environment variables

| Variable | |
|---|---|
| `CGO_ENABLED` | `0` by default when unset or empty; set to `1` for app dependencies that need cgo, including in `golet dev`, `golet build` and `golet generate` |
| `GOLET_INSPECTOR` | `1` keeps the [inspector of native UI](ui/inspector.md) in production builds of `golet build`, which leave it out |
| `GOLET_UPDATER_PRIVATE_KEY` | the secret key that signs updates, for `golet build` |
| `GOLET_WINDOWS_CERTIFICATE_PASSWORD` | the password of `windows.certificate`, for `golet build` |
| `GOLET_CLI_BINARY` | a build of the CLI for the `golet-cli` package to run |
| `GOLET_ENV` | `production` makes an app behave like a production build, e.g. without the web inspector |
| `GOLET_LOCALE` | the language the app runs in, over the one it keeps, which `golet dev -lang` sets |

For an app with dependencies that need cgo:

```sh
CGO_ENABLED=1 go tool golet dev
CGO_ENABLED=1 go tool golet build
```
