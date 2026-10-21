// Author: HyperparameterSI

// Builds dist/ of the plugin package in the working directory: an ES
// module, which imports golet-runtime, and its type declarations. dist/ is
// not committed, like golet-runtime's.
import { rm } from "node:fs/promises";
import { join } from "node:path";

const dir = process.cwd();
await rm(join(dir, "dist"), { recursive: true, force: true });
const result = await Bun.build({
  entrypoints: [join(dir, "src/index.ts")],
  outdir: join(dir, "dist"),
  format: "esm",
  target: "browser",
  external: ["golet-runtime"],
});
if (!result.success) {
  for (const log of result.logs) console.error(log);
  process.exit(1);
}
const tsc = Bun.spawnSync(["tsc", "-p", join(dir, "tsconfig.json")], {
  stdout: "inherit",
  stderr: "inherit",
});
if (tsc.exitCode !== 0) process.exit(1);
console.log(`wrote ${join(dir, "dist")}`);
