// Author: HyperparameterSI

import { useState } from "react";
import type { ProcessSample } from "@/golet";
import { time } from "@/lib/format-metrics";

export function PerformanceChart({
  samples,
  metric,
  color = "#5369df",
  format,
  ceiling = 1,
  label,
}: {
  samples: ProcessSample[];
  metric: (s: ProcessSample) => number;
  color?: string;
  format: (n: number) => string;
  ceiling?: number;
  label: string;
}) {
  const [hover, setHover] = useState<number | null>(null);
  const values = samples.map(metric);
  const max = Math.max(ceiling, ...values) * 1.15;
  const points = values.map(
    (v, i) => `${40 + (i / Math.max(1, values.length - 1)) * 720},${155 - (v / max) * 130}`,
  );
  const selected = hover === null ? undefined : samples[Math.min(hover, samples.length - 1)];
  return (
    <div className="chart-wrap">
      <div className="chart-readout">
        {selected
          ? `${time(selected.timestamp)} · ${format(metric(selected))}`
          : `${label} · 移入图表查看采样值`}
      </div>
      <svg
        viewBox="0 0 780 186"
        role="img"
        aria-label={`${label}趋势，${samples.length}个真实采样点`}
        onMouseLeave={() => setHover(null)}
        onMouseMove={(e) => {
          const r = e.currentTarget.getBoundingClientRect();
          setHover(
            Math.max(
              0,
              Math.min(
                samples.length - 1,
                Math.round(
                  ((((e.clientX - r.left) / r.width) * 780 - 40) / 720) * (samples.length - 1),
                ),
              ),
            ),
          );
        }}
      >
        {[0, 1, 2, 3].map((i) => (
          <g key={i}>
            <line
              x1="40"
              x2="760"
              y1={25 + (i * 130) / 3}
              y2={25 + (i * 130) / 3}
              stroke="#edf0f5"
              strokeDasharray="3 5"
            />
            <text x="0" y={29 + (i * 130) / 3} fill="#929aab" fontSize="9">
              {format(max * (1 - i / 3))}
            </text>
          </g>
        ))}
        {points.length > 1 && (
          <>
            <path d={`M40,155 L${points.join(" L")} L760,155 Z`} fill={color} opacity=".07" />
            <polyline
              points={points.join(" ")}
              fill="none"
              stroke={color}
              strokeWidth="2.3"
              strokeLinejoin="round"
            />
          </>
        )}
        {selected && (
          <circle
            cx={40 + (Math.min(hover!, samples.length - 1) / Math.max(1, samples.length - 1)) * 720}
            cy={155 - (metric(selected) / max) * 130}
            r="4"
            fill={color}
            stroke="white"
            strokeWidth="2"
          />
        )}
        <text x="40" y="180" fill="#929aab" fontSize="10">
          {samples[0] ? time(samples[0].timestamp) : "等待采样"}
        </text>
        <text x="760" y="180" textAnchor="end" fill="#929aab" fontSize="10">
          {samples.at(-1) ? time(samples.at(-1)!.timestamp) : "—"}
        </text>
      </svg>
    </div>
  );
}
