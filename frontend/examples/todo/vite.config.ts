// Author: HyperparameterSI

import { defineConfig } from "vite";

export default defineConfig({
  server: {
    // Golet selects the port; standalone Vite uses a free port too.
    port: Number(process.env.GOLET_DEV_PORT ?? 0),
    strictPort: true,
    // The development app and the packaged builds are not the page's.
    watch: { ignored: ["**/.golet/**", "**/build/**"] },
  },
});
