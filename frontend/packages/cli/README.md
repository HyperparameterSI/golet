# golet-cli

The `golet` command line tool of [Golet](https://github.com/HyperparameterSI/golet), the
desktop app framework for Go whose windows show web pages, in the system
webview, or native UI written in Go, as prebuilt binaries for macOS, Linux
and Windows (arm64 and x64).

Create a project, which depends on golet-cli and runs it from its scripts:

```sh
bunx golet-cli init my-app    # or: npx golet-cli init my-app
cd my-app
bun run dev                  # golet dev: the app with live reload
bun run build                # golet build: the packaged app
```

In a project that depends on golet-cli, `bunx golet <command>` runs its
`golet`. Or install it globally with `npm install -g golet-cli`. The commands
are those of the CLI installed with Go
(`go install github.com/HyperparameterSI/golet/backend/cmd/golet@latest`): `init`, `dev`,
`build`, `generate`, `install-skills`, `keygen` and `doctor`. Building apps
still needs [Go](https://go.dev/dl/), which compiles them.

Outside such a project, run `bunx golet-cli` or `npx golet-cli` rather than
`bunx golet`: npm has an unrelated package named golet.

## How it works

Each platform's binary is in a package of its own, such as
`@golet/cli-darwin-arm64`, which golet-cli lists as an optional dependency:
package managers install only the one that matches the machine. The `golet`
command then runs that binary. Set `GOLET_CLI_BINARY` to the path of another
build of the CLI to use it instead.

## golet.config.ts

`defineConfig` types the configuration of a project, `golet.config.ts`:

```ts
import { defineConfig } from "golet-cli";

export default defineConfig({
  name: "My App",
  identifier: "com.example.myapp",
  devUrl: "http://localhost:5173",
});
```

## The binary

From JavaScript, `binaryPath()` returns the path of the binary:

```ts
import { binaryPath } from "golet-cli";

Bun.spawnSync([binaryPath(), "generate"]);
```
