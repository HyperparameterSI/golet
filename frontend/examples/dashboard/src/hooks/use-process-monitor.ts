// Author: HyperparameterSI

import { useEffect, useRef, useState } from "react";
import { ProcessMonitor, type ProcessSnapshot } from "@/golet";
import { mergeSamples } from "@/lib/metric-history";

/** Poll incrementally. Hidden and paused windows have no polling timer. */
export function useProcessMonitor(paused: boolean, refreshIntervalMS: number) {
  const [snapshot, setSnapshot] = useState<ProcessSnapshot | null>(null);
  const [error, setError] = useState("");
  const [latency, setLatency] = useState(0);
  const sequence = useRef(0);
  const inFlight = useRef(false);

  useEffect(() => {
    if (paused) return;
    let cancelled = false;
    let timer: ReturnType<typeof setTimeout> | undefined;
    const schedule = () => {
      clearTimeout(timer);
      if (!cancelled && !document.hidden) timer = setTimeout(() => void poll(), refreshIntervalMS);
    };
    const poll = async () => {
      if (cancelled || document.hidden) return;
      if (inFlight.current) {
        schedule();
        return;
      }
      inFlight.current = true;
      const startedAt = performance.now();
      try {
        const update = await ProcessMonitor.updates(sequence.current);
        if (!cancelled) {
          sequence.current = update.samples.at(-1)?.sequence ?? sequence.current;
          setSnapshot((previous) => ({
            ...update,
            samples: mergeSamples(previous?.samples ?? [], update.samples),
          }));
          setLatency(performance.now() - startedAt);
          setError("");
        }
      } catch (error) {
        if (!cancelled) setError(error instanceof Error ? error.message : String(error));
      } finally {
        inFlight.current = false;
        schedule();
      }
    };
    const onVisibilityChange = () => {
      clearTimeout(timer);
      if (!document.hidden) void poll();
    };
    document.addEventListener("visibilitychange", onVisibilityChange);
    void poll();
    return () => {
      cancelled = true;
      clearTimeout(timer);
      document.removeEventListener("visibilitychange", onVisibilityChange);
    };
  }, [paused, refreshIntervalMS]);

  return { snapshot, error, latency };
}
