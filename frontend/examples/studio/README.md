# Golet Studio frontend

TanStack Start SPA + React + shadcn-style components + Vite.

- Routes: `src/routes/`
- Application: `src/components/app.tsx`
- UI components: `src/components/ui/`
- Go client: `src/golet.ts` (generated; do not edit)
- Go service: `../../../backend/examples/studio/main.go`

From the repository root: `bun run dev` starts both Go and the frontend, `bun run build:web` checks and builds the frontend, and `bun run build:app` creates the desktop app. From this directory, `bun run dev:web` opens a browser-only preview; Go calls require the desktop WebView.

The SPA shell is emitted at `dist/client/index.html`. The Go build embeds that directory; `dist/server` is build-time output only. To add components, run `bunx shadcn@latest add dialog` here.
