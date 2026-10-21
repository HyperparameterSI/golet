# SQLite

SQLite for Golet, with Go connections and the `@golet-plugins/sqlite`
TypeScript client. The plugin ships prebuilt SQLite libraries and loads
them through purego, so apps build without cgo or a C compiler.

Register `golet.Use(sqlite.Plugin)` only when the web frontend uses
`@golet-plugins/sqlite`. Go and native UI apps call `sqlite.Open` directly;
they need no plugin registration and close their own connections.

```go
import (
    "github.com/HyperparameterSI/golet/backend"
    "github.com/HyperparameterSI/golet/backend/plugins/sqlite"
)

func main() {
    golet.Use(sqlite.Plugin)
    // Create windows, then call golet.App.Run().
}
```

```ts
import { open } from "@golet-plugins/sqlite";

const db = await open("app.db");
await db.execute("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, title TEXT, done INTEGER)");
const inserted = await db.execute("INSERT INTO tasks(title, done) VALUES (?, ?)", ["Learn Golet", false]);
const { columns, rows } = await db.query("SELECT id, title, done FROM tasks WHERE done = ?", [false]);
await db.transaction([
  { sql: "UPDATE tasks SET done = ? WHERE id = ?", params: [true, inserted.lastInsertId] },
  { sql: "INSERT INTO tasks(title, done) VALUES (?, ?)", params: ["Ship the app", false] },
]);
await db.close();
```

Use `bun add @golet-plugins/sqlite` once the npm package is published. In a
checkout, `bun install && bun run build` installs and builds the workspace
package. See [the full guide](../../docs/plugins/sqlite.md) for the Go API,
types, cancellation, options and packaging.

## Building the native libraries

Prebuilt libraries are published in
[sqlite-3.53.4-1](https://github.com/HyperparameterSI/golet/releases/tag/sqlite-3.53.4-1).
CI and unpackaged programs download the matching library and verify its
SHA-256. Testing and building apps require no Zig installation.

To update SQLite or its C shim, with Zig **0.16.0** on the path:

```sh
go generate ./plugins/sqlite
# Or build just this machine, from plugins/sqlite:
go run ./internal/libbuild -target host
```

The generator downloads SQLite **3.53.4** from sqlite.org and checks its
published SHA3-256 before extracting `sqlite3.c` and `sqlite3.h`. It builds
macOS, Linux and Windows for amd64 and arm64, writes SHA-256 and release
URLs to `golet-plugin.json`, and populates the CLI's native-library cache.
`golet dev` and `golet build` discover the manifest and bundle the matching
library; macOS universal builds combine both architectures. Linux targets
glibc 2.28 or later; macOS targets 12 or later. Binaries in `build/` are
ignored by git.

New binaries must be published before a fresh checkout or app build can
download them. The default asset release is `sqlite-3.53.4-1`:

```sh
gh release create sqlite-3.53.4-1 plugins/sqlite/build/* --latest=false \
  --title "SQLite 3.53.4 for Golet" --notes "SQLite compiled with Zig 0.16.0 and the Golet C shim."
```

Change the asset tag (or pass `-release <URL>`) when the shim, compiler or
compile options change; regenerate and publish the manifest and binaries
together. `-src <directory>` uses a local amalgamation instead of the
verified download, intended for development.

`GOLET_SQLITE_LIBRARY` overrides the library path. Otherwise the plugin
looks in app resources, next to the executable, then the CLI's cache.
Unpackaged programs download a missing library and verify its SHA-256;
packaged apps never download native code at run time.

SQLite is public domain; the plugin and C shim use the repository's MIT
license. FTS5, RTree and JSON are enabled, with foreign keys on by default.
Loadable extensions, shared cache and deprecated APIs are omitted.
