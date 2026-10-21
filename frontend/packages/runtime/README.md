# golet-runtime

Typed access to the runtime of [Golet](https://github.com/HyperparameterSI/golet) apps,
for their web frontends.

Golet injects its runtime into every page of an app as `window.golet`. This
package gives it types and imports. The client that `golet generate` writes
for your bound Go services is built on it, so you mostly use that:

```ts
import { Greeter, events } from "./golet"; // generated

const greeting = await Greeter.greet("Ada");
events.tick.on((time) => console.log(time));
```

Use the package directly for the window and the platform:

```ts
import { currentWindow, isCallError, isGolet, runtime } from "golet-runtime";

if (isGolet() && runtime().platform === "darwin") {
  document.body.classList.add("mac");
}
await currentWindow.toggleMaximize();
```

A Go method that streams values through a `*golet.Channel[T]` parameter takes
a `Channel` in its place:

```ts
import { Channel } from "golet-runtime";

const lines = new Channel<string>();
const done = Shell.tail("ls -R", lines);
for await (const line of lines) console.log(line);
await done;
```

A Go method returning an error rejects with a `CallError` (see
`isCallError`). Outside a Golet window, e.g. when the dev server is opened in a
browser, calls reject and `isGolet()` is false.

`golet-runtime/i18n` formats the messages of the app's catalogs in its pages,
as package `i18n` does in Go: messages of ICU MessageFormat, with the plural
rules, numbers and dates of `Intl`. `golet generate` writes a typed module of
the catalogs next to the client (`locales.ts`) on top of it:

```ts
import { messages } from "./locales"; // generated

status.textContent = messages.filesSelected(3); // "3 files selected"
```

`formatMessage(locale, message, args)` formats a message of your own.
