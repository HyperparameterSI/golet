// Author: HyperparameterSI

import { defineConfig } from "../../../frontend/packages/cli/index.js";

export default defineConfig({
  name: "Todo",
  identifier: "dev.golet.todo",
  version: "1.0.0",
  devUrl: "http://localhost:5173",
  devCommand: "npm --prefix ../../../frontend/packages/runtime run build && npm --prefix ../../../frontend/examples/todo run dev:web",
  buildCommand: "npm --prefix ../../../frontend/packages/runtime run build && npm --prefix ../../../frontend/examples/todo run build:web",
  frontendDist: "../../../frontend/examples/todo/dist",
  bindings: "../../../frontend/examples/todo/src/golet.ts",
  out: "build",
});
