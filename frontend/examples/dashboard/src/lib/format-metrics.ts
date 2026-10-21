// Author: HyperparameterSI

export const bytes = (n: number) =>
  n >= 2 ** 30
    ? `${(n / 2 ** 30).toFixed(2)} GB`
    : n >= 2 ** 20
      ? `${(n / 2 ** 20).toFixed(1)} MB`
      : `${(n / 1024).toFixed(1)} KB`;
export const number = (n: number) => n.toLocaleString("en-US");
export const time = (n: number) => new Date(n).toLocaleTimeString("zh-CN", { hour12: false });
export const duration = (s: number) =>
  `${Math.floor(s / 3600)
    .toString()
    .padStart(2, "0")}:${Math.floor((s / 60) % 60)
    .toString()
    .padStart(2, "0")}:${Math.floor(s % 60)
    .toString()
    .padStart(2, "0")}`;
export const average = (values: number[]) =>
  values.length ? values.reduce((a, b) => a + b, 0) / values.length : 0;
