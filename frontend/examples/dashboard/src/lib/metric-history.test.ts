// Author: HyperparameterSI

import { describe, expect, test } from "bun:test";
import type { ProcessSample } from "../golet";
import { mergeSamples } from "./metric-history";

const sample = (sequence: number) => ({ sequence }) as ProcessSample;

describe("incremental metric history", () => {
  test("keeps history immutable and bounded", () => {
    const previous = Array.from({ length: 300 }, (_, i) => sample(i + 1));
    const merged = mergeSamples(previous, [sample(301), sample(302)]);
    expect(merged.length).toBe(300);
    expect(merged[0]!.sequence).toBe(3);
    expect(previous[0]!.sequence).toBe(1);
    expect(merged.at(-1)!.sequence).toBe(302);
  });
  test("retains history on empty poll and resets on restart", () => {
    const previous = [sample(100)];
    expect(mergeSamples(previous, [])).toBe(previous);
    expect(mergeSamples(previous, [sample(1)])).toEqual([sample(1)]);
  });
  test("bounds a complete catch-up response", () => {
    expect(
      mergeSamples(
        [sample(1)],
        Array.from({ length: 350 }, (_, i) => sample(i + 2)),
      ).length,
    ).toBe(300);
  });
});
