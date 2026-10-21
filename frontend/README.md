# Web frontend

- `examples/studio`: TanStack Start SPA + React + shadcn-style components + Vite.
- `examples/todo`: lightweight TypeScript example.
- `packages`: shared runtime, injected bridge source, npm CLI packages.
- `plugins`: TypeScript clients for Go plugins.
- `website`: TanStack Start documentation website.

Install dependencies at the repository root with `bun install`, then run `bun run build` to build shared packages. `bun run dev` launches Studio with its real Go backend. `bun run build:web` builds only the static frontend. `bun run dev:website` starts the documentation site.

Studio uses file routes in `examples/studio/src/routes`, reusable UI in `src/components/ui`, and application content in `src/components/app.tsx`. Go calls use generated `src/golet.ts`; edit the service under `backend/examples/studio` and run `go run ./backend/cmd/golet generate backend/examples/studio` from the root.

Start prerenders `dist/client/index.html`. The Go desktop application embeds `dist/client` only. Do not add Start server functions to desktop apps without separately provisioning a server. Plain browser preview has no desktop IPC and shows a disconnected state.
