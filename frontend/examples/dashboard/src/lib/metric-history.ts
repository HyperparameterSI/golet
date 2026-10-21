// Author: HyperparameterSI

import type { ProcessSample } from "@/golet";

export const MAX_HISTORY_SAMPLES = 300;

/** Merge monotonic cursors, resetting after a process restart; never mutate inputs. */
export function mergeSamples(
  previous: ProcessSample[],
  incoming: ProcessSample[],
): ProcessSample[] {
  if (incoming.length === 0) return previous;
  const lastSequence = previous.at(-1)?.sequence ?? 0;
  if (incoming[0]!.sequence <= lastSequence) return incoming.slice(-MAX_HISTORY_SAMPLES);
  const retainedCount = Math.max(0, MAX_HISTORY_SAMPLES - incoming.length);
  const retained = retainedCount > 0 ? previous.slice(-retainedCount) : [];
  return [...retained, ...incoming].slice(-MAX_HISTORY_SAMPLES);
}
