// Author: HyperparameterSI

import { defineConfig } from "../../../frontend/packages/cli/index.js";

export default defineConfig({
  name: "Golet Observatory",
  identifier: "dev.golet.dashboard",
  version: "0.1.0",
  // `golet dev` runs devCommand and loads devUrl; `golet build` runs
  // buildCommand and embeds frontendDist into the app.
  devUrl: "http://localhost:5174",
  devCommand: "npm --prefix ../../../frontend/packages/runtime run build && npm --prefix ../../../frontend/examples/dashboard run dev:web",
  buildCommand: "npm --prefix ../../../frontend/packages/runtime run build && npm --prefix ../../../frontend/examples/dashboard run build:web",
  frontendDist: "../../../frontend/examples/dashboard/dist/client",
  bindings: "../../../frontend/examples/dashboard/src/golet.ts",
  out: "build",
});
