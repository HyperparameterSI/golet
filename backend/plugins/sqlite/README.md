# SQLite

SQLite for Golet, with Go connections and the `@golet-plugins/sqlite`
TypeScript client. Its SQLite is a native module written in Zig
(`internal/sqlitecore`), published prebuilt for each platform, so apps
build without cgo, a C compiler or Zig.

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

## Building the native module

`internal/sqlitecore` declares the interface the plugin calls, `Engine`,
and `internal/sqlitecore/zig` implements it, compiling SQLite **3.53.4**,
which `build.zig.zon` pins by the hash of sqlite.org's amalgamation. Its
library exports only the module's functions: SQLite's own symbols stay
hidden, so a process that loaded another SQLite globally (WebKit does, on
Linux) cannot interpose on them.

With Zig **0.16.0** on the path, from the repository's root:

```sh
# Build this machine's library, which go test and go run use:
go run ./cmd/golet native build -target host plugins/sqlite/internal/sqlitecore
# Build and publish every platform's:
go run ./cmd/golet native build -upload plugins/sqlite/internal/sqlitecore
```

Linux targets glibc 2.28 or later; macOS targets 12 or later.

SQLite is public domain; the plugin and its module use the repository's MIT
license. FTS5, RTree and JSON are enabled, with foreign keys on by default.
Loadable extensions, shared cache and deprecated APIs are omitted.
