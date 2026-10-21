// Author: HyperparameterSI

import { defineConfig } from "../../../frontend/packages/cli/index.js";

export default defineConfig({
  name: "Golet Studio",
  identifier: "com.example.goletshadcndemo",
  version: "0.1.0",
  // `golet dev` runs devCommand and loads devUrl; `golet build` runs
  // buildCommand and embeds frontendDist into the app.
  devUrl: "http://localhost:0",
  devCommand: "npm --prefix ../../../frontend/packages/runtime run build && npm --prefix ../../../frontend/examples/studio run dev:web",
  buildCommand: "npm --prefix ../../../frontend/packages/runtime run build && npm --prefix ../../../frontend/examples/studio run build:web",
  frontendDist: "../../../frontend/examples/studio/dist/client",
  bindings: "../../../frontend/examples/studio/src/golet.ts",
  out: "build",
});
