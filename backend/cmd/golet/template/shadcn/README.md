# Golet desktop application

This project uses **Golet + Go** for the desktop backend and **TanStack Start + React + TypeScript + Vite + shadcn/ui** for the frontend.

## Start developing

```sh
pnpm install
pnpm dev
```

Run commands from this directory. Install Go and Node.js/pnpm first. Use the package manager matching the existing lockfile if this project already has one.

- `backend/`: Go services and desktop application entry point.
- `frontend/src/routes/`: TanStack Router pages.
- `frontend/src/components/ui/`: editable shadcn components.
- `frontend/src/style.css`: Tailwind and semantic light/dark theme tokens.
- `frontend/src/golet.ts`: generated Go bindings; run `pnpm generate` after API changes.
- `golet.config.ts`: development commands and desktop packaging settings.

TanStack Start runs in SPA mode. The packaged app embeds `frontend/dist/client` in the system WebView and does not need a Node.js server. `pnpm dev:web` alone cannot provide desktop IPC.

## Checks and packaging

```sh
pnpm typecheck
pnpm build:web
go test ./...
pnpm build
```

Desktop packaging requires the target platform's build tools. Output goes to `build/`.

On macOS, select the target explicitly (including when building on Apple Silicon for Intel):

```sh
pnpm build:mac:intel      # Intel Mac, build/darwin-amd64/
pnpm build:mac:arm64      # Apple Silicon, build/darwin-arm64/
pnpm build:mac:universal  # Both architectures, build/darwin-universal/
```

Each target produces an app bundle and DMG. Signing and DMG creation require macOS.
Cross-compilation does not verify execution on the target machine. Optional native
plugins must provide matching architecture binaries; the separate Chromium host
is not automatically bundled by these commands.

## Developing with AI

Project conventions live in [AGENTS.md](AGENTS.md). Ask your AI coding tool to read it before implementing a feature. Tool support varies; attach the file or paste its requirements when automatic loading is unavailable.

Example prompt:

> Read AGENTS.md and golet.config.ts first. Build this feature with the existing TanStack Start + React + Vite + shadcn/ui setup. Keep frontend code in frontend/ and Go services in backend/. Use the generated Golet client for IPC, preserve SPA packaging, and run the relevant checks. Feature: …

CLAUDE.md and .github/copilot-instructions.md point to the same conventions. Golet maintenance guidance is included in .agents/skills/. Explicit user requirements take precedence over these defaults.
