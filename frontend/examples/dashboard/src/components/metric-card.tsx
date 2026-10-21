// Author: HyperparameterSI

import type { ReactNode } from "react";
import { Card } from "./ui/card";

export function MetricCard({
  icon,
  title,
  value,
  detail,
  data,
  color,
}: {
  icon: ReactNode;
  title: string;
  value: string;
  detail: string;
  data: number[];
  color: string;
}) {
  const max = Math.max(1, ...data);
  return (
    <Card className="metric">
      <div className="metric-top">
        <span>{title}</span>
        <span style={{ color }}>{icon}</span>
      </div>
      <strong>{value}</strong>
      <div className="metric-bottom">
        <span>{detail}</span>
        <svg viewBox="0 0 100 28" aria-hidden="true">
          <polyline
            points={data
              .map((v, i) => `${(i / Math.max(1, data.length - 1)) * 100},${26 - (v / max) * 22}`)
              .join(" ")}
            fill="none"
            stroke={color}
            strokeWidth="2"
          />
        </svg>
      </div>
    </Card>
  );
}
