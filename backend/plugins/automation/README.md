# automation

Drive a running [Golet](https://github.com/HyperparameterSI/golet/backend) app from another
process, as end-to-end test runners and agents do: list its windows, read
the elements of their native UI with the app's keys, act on them, send
pointer, wheel, key and text input, choose menu items, evaluate scripts in
pages, capture windows, and wait for them to settle. One protocol on
macOS, Linux and Windows: newline-delimited JSON over a Unix domain
socket that only the user can open.

```go
//go:build e2e

package main

import _ "github.com/HyperparameterSI/golet/backend/plugins/automation"
```

```sh
GOLET_AUTOMATION=/tmp/app.sock ./app
echo '{"id": 1, "method": "windows"}' | socat - UNIX-CONNECT:/tmp/app.sock
```

It listens only when `GOLET_AUTOMATION` is set, and builds without the
inspector (`-tags golet_noinspector`, the release builds of `golet build`)
leave it out. See [the documentation](../../docs/plugins/automation.md).
