// Author: HyperparameterSI

import { useMemo, useState } from "react";
import {
  Activity,
  ArrowDownToLine,
  ArrowUpRight,
  Check,
  ChevronRight,
  CircleHelp,
  Cpu,
  Database,
  Gauge,
  MemoryStick,
  Pause,
  Play,
  Radio,
  ShieldCheck,
  Terminal,
  Timer,
  Waves,
} from "lucide-react";
import { useProcessMonitor } from "@/hooks/use-process-monitor";
import { PerformanceChart as Chart } from "./performance-chart";
import { MetricCard as Metric } from "./metric-card";
import { bytes, number, time, duration, average } from "@/lib/format-metrics";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import "./dashboard.css";
import brandLogo from "@/assets/golet-logo.svg";

export function App() {
  const [paused, setPaused] = useState(false);
  const [windowSize, setWindowSize] = useState(60);
  const [refresh, setRefresh] = useState(1000);
  const [section, setSection] = useState("overview");
  const [exportError, setExportError] = useState("");
  const { snapshot, error, latency } = useProcessMonitor(paused, refresh);
  const samples = useMemo(() => snapshot?.samples.slice(-windowSize) ?? [], [snapshot, windowSize]);
  const latest = snapshot?.samples.at(-1);
  const cpuSamples = samples.filter((s) => s.cpuAvailable);
  const gcEvents = (snapshot?.samples ?? [])
    .filter((s) => s.gcDelta > 0)
    .slice(-6)
    .reverse();
  const checks = latest
    ? [
        {
          name: "CPU 负载",
          ok: !latest.cpuAvailable ? null : latest.cpu < 80,
          detail: !latest.cpuAvailable
            ? "等待 CPU 差分采样或平台不支持"
            : `${latest.cpu.toFixed(1)}% / 单核 100%，观察阈值 80%`,
        },
        {
          name: "GC 暂停",
          ok: latest.gcCycles === 0 ? null : latest.gcPauseP95MS < 10,
          detail:
            latest.gcCycles === 0
              ? "进程尚未发生 GC"
              : `最近至多 256 次 P95 ${latest.gcPauseP95MS.toFixed(3)} ms，阈值 10 ms`,
        },
        {
          name: "Goroutine 数量",
          ok: latest.goroutines < 1000,
          detail: `${number(latest.goroutines)} 个活跃协程，观察阈值 1,000`,
        },
        {
          name: "采样开销",
          ok: latest.sampleDurationMS < 10,
          detail: `${latest.sampleDurationMS.toFixed(3)} ms / 次，阈值 10 ms`,
        },
      ]
    : [];
  const warnings = checks.filter((c) => c.ok === false).length;
  const jump = (id: string) => {
    setSection(id);
    document.getElementById(id)?.scrollIntoView({
      behavior: matchMedia("(prefers-reduced-motion: reduce)").matches ? "auto" : "smooth",
      block: "start",
    });
  };
  function exportSnapshot() {
    if (!snapshot) return;
    try {
      const url = URL.createObjectURL(
        new Blob(
          [
            JSON.stringify(
              {
                exportedAt: new Date().toISOString(),
                scope: "Go process only; excludes WebView; runtime memory is not RSS",
                ...snapshot,
              },
              null,
              2,
            ),
          ],
          { type: "application/json" },
        ),
      );
      const a = document.createElement("a");
      a.href = url;
      a.download = `golet-metrics-${Date.now()}.json`;
      a.click();
      setTimeout(() => URL.revokeObjectURL(url), 10000);
      setExportError("");
    } catch {
      setExportError("导出失败，请重试。");
    }
  }
  return (
    <div className="observatory">
      <aside className="sidebar">
        <a className="brand" href="#overview" aria-label="Golet 运行概览">
          <img className="brand-logo" src={brandLogo} alt="Golet" width="176" height="44" />
        </a>
        <div className="workspace">
          <span className="workspace-icon">
            <Activity size={18} />
          </span>
          <div>
            Observatory<small>运行性能观测台</small>
          </div>
          <ChevronRight size={14} />
        </div>
        <div className="nav-label">WORKSPACE</div>
        <nav>
          {[
            ["overview", "运行概览", Gauge],
            ["resources", "资源趋势", Activity],
            ["runtime", "Go 运行时", Cpu],
            ["diagnostics", "健康检查", ShieldCheck],
          ].map(([id, title, Icon]) => {
            const I = Icon as typeof Gauge;
            return (
              <button
                key={String(id)}
                className={section === id ? "selected" : ""}
                onClick={() => jump(String(id))}
              >
                <I size={17} />
                {String(title)}
                {id === "diagnostics" && warnings > 0 && <b>{warnings}</b>}
              </button>
            );
          })}
        </nav>
        <div className="sidebar-note">
          <span className="live-dot" /> LOCAL PROCESS<code>PID {snapshot?.pid ?? "—"}</code>
          <small>
            数据来自本机 Go 进程
            <br />
            每秒采样 · 最多保留 300 点
          </small>
        </div>
        <div className="sidebar-footer">
          <Terminal size={15} /> Golet Framework <span>v0.1.0</span>
        </div>
      </aside>
      <div className="dashboard-main">
        <header className="topbar">
          <div>
            Workspace <ChevronRight size={13} />
            <span>Observatory</span>
          </div>
          <span className="local-pill">
            <span className="live-dot" /> 本地监控
          </span>
        </header>
        <main className="dashboard-body" id="overview">
          <div className="page-heading">
            <div>
              <div className="eyebrow">PERFORMANCE MONITOR</div>
              <h1>
                每一刻运行，都清晰可见<span>.</span>
              </h1>
              <p>实时洞察 Go 应用的资源消耗、运行状态与性能表现。</p>
            </div>
            <Button variant="outline" size="sm" disabled={!snapshot} onClick={exportSnapshot}>
              <ArrowDownToLine />
              导出快照
            </Button>
          </div>
          {(error || exportError) && (
            <div className="error-banner" role="alert">
              <strong>
                {snapshot ? "连接中断，当前显示最后一次采样。" : "尚未连接 Go 采样服务。"}
              </strong>{" "}
              请通过 <code>pnpm dev:dashboard</code> 启动桌面示例。
              <small>{error || exportError}</small>
            </div>
          )}
          <div className="status-strip">
            <div>
              <span className={`live-dot ${paused || error ? "muted-dot" : ""}`} />
              <strong>
                {error ? "连接异常" : paused ? "显示已暂停" : snapshot ? "实时监测中" : "正在连接"}
              </strong>
              <span className="status-separator" />
              <span>{snapshot ? `${snapshot.os} / ${snapshot.arch}` : "等待进程信息"}</span>
              <span className="status-separator" />
              <Timer size={14} />
              <span>{snapshot ? duration(snapshot.uptime) : "—"}</span>
            </div>
            <div>
              <label>
                刷新{" "}
                <select
                  aria-label="刷新频率"
                  value={refresh}
                  onChange={(e) => setRefresh(Number(e.target.value))}
                >
                  <option value={1000}>1 秒</option>
                  <option value={2000}>2 秒</option>
                  <option value={5000}>5 秒</option>
                </select>
              </label>
              <button className="pause-control" onClick={() => setPaused(!paused)}>
                {paused ? <Play size={13} /> : <Pause size={13} />}
                {paused ? "继续" : "暂停"}
              </button>
            </div>
          </div>
          <div className="metric-grid">
            <Metric
              icon={<Cpu size={18} />}
              title="进程 CPU"
              value={latest?.cpuAvailable ? `${latest.cpu.toFixed(1)}%` : "—"}
              detail={`窗口均值 ${cpuSamples.length ? average(cpuSamples.map((s) => s.cpu)).toFixed(1) + "%" : "—"}`}
              data={cpuSamples.map((s) => s.cpu)}
              color="#5369df"
            />
            <Metric
              icon={<MemoryStick size={18} />}
              title="Go 活跃堆内存"
              value={latest ? bytes(latest.heap) : "—"}
              detail={`窗口峰值 ${latest ? bytes(Math.max(...samples.map((s) => s.heap))) : "—"}`}
              data={samples.map((s) => s.heap)}
              color="#269a8b"
            />
            <Metric
              icon={<Waves size={18} />}
              title="Goroutines"
              value={latest ? number(latest.goroutines) : "—"}
              detail={`${latest?.threadRecords ?? "—"} 个运行时线程记录`}
              data={samples.map((s) => s.goroutines)}
              color="#a17ade"
            />
            <Metric
              icon={<Timer size={18} />}
              title="最近 GC 暂停"
              value={latest?.gcCycles ? `${latest.gcPauseMS.toFixed(3)} ms` : "—"}
              detail={`进程累计 ${latest?.gcCycles ?? 0} 次 GC`}
              data={samples.map((s) => s.gcPauseMS)}
              color="#d39a39"
            />
          </div>
          <section id="resources">
            <div className="section-heading">
              <h2>
                资源趋势 <span>RESOURCE ACTIVITY</span>
              </h2>
              <div className="range-controls">
                {[60, 180, 300].map((n) => (
                  <button
                    className={windowSize === n ? "active" : ""}
                    onClick={() => setWindowSize(n)}
                    key={n}
                  >
                    {n / 60} 分钟
                  </button>
                ))}
              </div>
            </div>
            <div className="charts-grid">
              <Card className="panel">
                <div className="panel-heading">
                  <h3>
                    <span className="legend-dot" />
                    CPU 使用率
                  </h3>
                  <span>单核 100% · 可超过 100%</span>
                </div>
                <Chart
                  samples={cpuSamples}
                  metric={(s) => s.cpu}
                  format={(v) => `${v.toFixed(0)}%`}
                  ceiling={10}
                  label="进程 CPU"
                />
                <div className="chart-footer">
                  <span>
                    均值{" "}
                    <b>
                      {cpuSamples.length ? average(cpuSamples.map((s) => s.cpu)).toFixed(2) : "—"}%
                    </b>
                  </span>
                  <span>
                    峰值{" "}
                    <b>
                      {cpuSamples.length
                        ? Math.max(...cpuSamples.map((s) => s.cpu)).toFixed(2)
                        : "—"}
                      %
                    </b>
                  </span>
                  <span>{snapshot?.cpus ?? "—"} 逻辑核心</span>
                </div>
              </Card>
              <Card className="panel">
                <div className="panel-heading">
                  <h3>
                    <span className="legend-dot teal" />
                    堆内存分配
                  </h3>
                  <span>HeapAlloc · 实际活跃对象</span>
                </div>
                <Chart
                  samples={samples}
                  metric={(s) => s.heap / 2 ** 20}
                  color="#269a8b"
                  format={(v) => `${v.toFixed(0)}M`}
                  ceiling={8}
                  label="堆内存 MB"
                />
                <div className="chart-footer">
                  <span>
                    分配速率 <b>{latest ? bytes(latest.allocRate) : "—"}/s</b>
                  </span>
                  <span>
                    活跃对象 <b>{latest ? number(latest.objects) : "—"}</b>
                  </span>
                </div>
              </Card>
            </div>
          </section>
          <section id="runtime">
            <div className="section-heading">
              <h2>
                运行时透视 <span>RUNTIME INSIGHTS</span>
              </h2>
              <span className="subtle">runtime.ReadMemStats · OS process CPU time</span>
            </div>
            <div className="runtime-grid">
              <Card className="panel memory-panel">
                <div className="panel-heading">
                  <h3>
                    <Database size={16} />
                    内存分布
                  </h3>
                  <span>Go runtime</span>
                </div>
                <div className="memory-total">
                  {latest ? bytes(latest.runtimeMemory) : "—"}
                  <small>运行时保留内存 · Sys − HeapReleased，非 RSS</small>
                </div>
                <div className="memory-bar">
                  {latest &&
                    [
                      latest.heap,
                      Math.max(0, latest.heapInuse - latest.heap),
                      latest.stack,
                      Math.max(0, latest.runtimeMemory - latest.heapInuse - latest.stack),
                    ].map((v, i) => (
                      <span
                        key={i}
                        style={{
                          flex: v || 0.001,
                          background: ["#5369df", "#a4aff3", "#43b7a6", "#e3e7ef"][i],
                        }}
                      />
                    ))}
                </div>
                <div className="memory-legend">
                  {[
                    ["活跃堆", latest?.heap, "#5369df"],
                    [
                      "堆页预留",
                      latest ? Math.max(0, latest.heapInuse - latest.heap) : undefined,
                      "#a4aff3",
                    ],
                    ["运行时栈", latest?.stack, "#43b7a6"],
                    [
                      "其他 / 未释放空闲",
                      latest
                        ? Math.max(0, latest.runtimeMemory - latest.heapInuse - latest.stack)
                        : undefined,
                      "#e3e7ef",
                    ],
                  ].map(([name, v, color]) => (
                    <div key={String(name)}>
                      <span>
                        <i style={{ background: String(color) }} />
                        {name}
                      </span>
                      <b>{v === undefined ? "—" : bytes(Number(v))}</b>
                    </div>
                  ))}
                </div>
                <div className="memory-note">
                  已归还操作系统 <b>{latest ? bytes(latest.heapReleased) : "—"}</b>
                </div>
              </Card>
              <Card className="panel">
                <div className="panel-heading">
                  <h3>
                    <Waves size={16} />
                    垃圾回收
                  </h3>
                  <span>Garbage collection</span>
                </div>
                <div className="kv-list">
                  {[
                    ["累计 GC 次数", latest ? number(latest.gcCycles) : "—"],
                    [
                      "暂停 P95（最近至多 256 次）",
                      latest?.gcCycles ? `${latest.gcPauseP95MS.toFixed(3)} ms` : "—",
                    ],
                    ["累计 STW 暂停", latest ? `${latest.gcTotalPauseMS.toFixed(2)} ms` : "—"],
                    [
                      "GC CPU 占比（进程生命周期）",
                      latest ? `${(latest.gcCPUFraction * 100).toFixed(3)}%` : "—",
                    ],
                    ["下次 GC 目标堆大小", latest ? bytes(latest.nextGC) : "—"],
                    ["累计分配", latest ? bytes(latest.totalAlloc) : "—"],
                  ].map(([k, v]) => (
                    <div key={k}>
                      <span>{k}</span>
                      <strong>{v}</strong>
                    </div>
                  ))}
                </div>
              </Card>
              <Card className="panel">
                <div className="panel-heading">
                  <h3>
                    <Terminal size={16} />
                    进程环境
                  </h3>
                  <span className="tag">GO</span>
                </div>
                <div className="kv-list">
                  {[
                    ["进程 ID", snapshot?.pid],
                    ["Go 版本", snapshot?.goVersion],
                    ["操作系统 / 架构", snapshot ? `${snapshot.os} / ${snapshot.arch}` : null],
                    [
                      "逻辑 CPU / GOMAXPROCS",
                      snapshot ? `${snapshot.cpus} / ${snapshot.gomaxprocs}` : null,
                    ],
                    ["IPC 往返耗时", snapshot ? `${latency.toFixed(1)} ms` : null],
                    ["本次采样耗时", latest ? `${latest.sampleDurationMS.toFixed(3)} ms` : null],
                  ].map(([k, v]) => (
                    <div key={String(k)}>
                      <span>{k}</span>
                      <strong>{v ?? "—"}</strong>
                    </div>
                  ))}
                </div>
              </Card>
            </div>
          </section>
          <section id="diagnostics">
            <div className="section-heading">
              <h2>
                检查与事件 <span>DIAGNOSTICS</span>
              </h2>
              <span className="subtle">观察阈值仅供排查参考，不代表故障诊断</span>
            </div>
            <div className="diagnostic-grid">
              <Card className="panel">
                <div className="panel-heading">
                  <h3>
                    <ShieldCheck size={16} />
                    运行健康检查
                  </h3>
                  <span className={warnings ? "warning-tag" : "healthy-tag"}>
                    {!latest ? "等待采样" : warnings ? `${warnings} 项需关注` : "未触发观察阈值"}
                  </span>
                </div>
                <div className="checks">
                  {checks.map((c) => (
                    <div key={c.name}>
                      <span className={`check-icon ${c.ok !== true ? "check-warning" : ""}`}>
                        {c.ok === true ? <Check size={14} /> : <CircleHelp size={14} />}
                      </span>
                      <div>
                        <strong>{c.name}</strong>
                        <p>{c.detail}</p>
                      </div>
                      <span>{c.ok === null ? "待观察" : c.ok ? "正常" : "关注"}</span>
                    </div>
                  ))}
                  {!latest && <p className="empty">连接后自动检查真实指标</p>}
                </div>
              </Card>
              <Card className="panel">
                <div className="panel-heading">
                  <h3>
                    <Radio size={16} />
                    近期 GC 事件
                  </h3>
                  <span>每秒聚合 · 最近 6 条</span>
                </div>
                <div className="events">
                  {gcEvents.map((s) => (
                    <div key={s.timestamp}>
                      <span className="event-mark" />
                      <div>
                        <strong>完成 {s.gcDelta} 次垃圾回收</strong>
                        <p>
                          最后一次暂停 {s.gcPauseMS.toFixed(3)} ms · 堆 {bytes(s.heap)}
                        </p>
                      </div>
                      <time>{time(s.timestamp)}</time>
                    </div>
                  ))}
                  {gcEvents.length === 0 && (
                    <div className="empty">
                      <Activity size={26} />
                      <p>采样窗口内暂无 GC 事件</p>
                      <small>发生回收后会自动显示，无需手动刷新</small>
                    </div>
                  )}
                </div>
              </Card>
            </div>
          </section>
          <footer className="dashboard-footer">
            <span>
              <CircleHelp size={13} />
              仅监控当前 Go 进程；不包含 WebView 子进程与整机资源。暂停仅冻结显示。
            </span>
            <span>
              {samples.length} 个样本 · {latest ? time(latest.timestamp) : "等待采样"}
              <ArrowUpRight size={12} />
            </span>
          </footer>
        </main>
      </div>
    </div>
  );
}
