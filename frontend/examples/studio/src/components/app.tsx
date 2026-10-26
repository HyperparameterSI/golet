// Author: HyperparameterSI
import { useEffect, useMemo, useState } from "react";
import {
  ArrowUpRight,
  Box,
  Check,
  Code2,
  Cpu,
  Globe,
  Layers,
  Loader2,
  Moon,
  Search,
  Sun,
} from "lucide-react";
import { BrowserDemo, Greeter, Showcase, type Demo, type BrowserStatus } from "@/golet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import goletMark from "@/assets/golet-mark.svg";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";

export function App() {
  const [dark, setDark] = useState(false);
  const [catalog, setCatalog] = useState<Demo[]>([]);
  const [category, setCategory] = useState("全部组件");
  const [search, setSearch] = useState("");
  const [selected, setSelected] = useState<Demo | null>(null);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [busy, setBusy] = useState(false);
  const [info, setInfo] = useState<Awaited<ReturnType<typeof Greeter.info>> | null>(null);
  const [name, setName] = useState("HyperparameterSI");
  const [reply, setReply] = useState("");
  useEffect(() => {
    setDark(matchMedia("(prefers-color-scheme: dark)").matches);
    let active = true;
    Promise.all([Showcase.catalog(), Greeter.info()])
      .then(([items, system]) => {
        if (active) {
          setCatalog(items);
          setInfo(system);
          setSelected(items[0] ?? null);
        }
      })
      .catch((e) => {
        if (active) setError(`请通过 pnpm dev 启动桌面后端：${String(e)}`);
      });
    return () => {
      active = false;
    };
  }, []);
  useEffect(() => {
    document.documentElement.classList.toggle("dark", dark);
  }, [dark]);
  const visible = useMemo(
    () =>
      catalog.filter(
        (d) =>
          (category === "全部组件" || d.category === category) &&
          `${d.title} ${d.description} ${d.controls.join(" ")}`
            .toLowerCase()
            .includes(search.toLowerCase()),
      ),
    [catalog, category, search],
  );
  async function open(demo: Demo) {
    setSelected(demo);
    setError("");
    setNotice("");
    if (demo.id === "Browser") return;
    setBusy(true);
    try {
      await Showcase.open(demo.id);
      setNotice(`已打开 ${demo.title}，可以在原生窗口中操作。`);
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }
  return (
    <div className="studio-shell">
      <aside className="studio-sidebar">
        <div className="flex items-center gap-3">
          <img className="studio-mark" src={goletMark} alt="Golet" width={42} height={42} />
          <div>
            <strong className="text-xl tracking-tight">
              golet<span className="text-brand-mark">.</span>
            </strong>
            <p className="text-xs text-muted-foreground">STUDIO / 0.1.0</p>
          </div>
        </div>
        <p className="mt-12 mb-4 text-[10px] tracking-[.24em] text-muted-foreground">
          COMPONENT EXPLORER
        </p>
        <nav aria-label="组件分类">
          {["全部组件", "原生控件", "高级组件", "插件"].map((item) => (
            <button
              key={item}
              className={`studio-nav ${category === item ? "active" : ""}`}
              onClick={() => setCategory(item)}
            >
              <Box size={17} />
              {item}
              <span>
                {item === "全部组件"
                  ? catalog.length
                  : catalog.filter((d) => d.category === item).length}
              </span>
            </button>
          ))}
        </nav>
        <div className="mt-auto space-y-3 pt-12 text-xs text-muted-foreground">
          <div className="flex items-center gap-2">
            <span className={`size-2 rounded-full ${info ? "bg-emerald-500" : "bg-amber-500"}`} />
            {info ? "Go 后端已连接" : "等待桌面连接"}
          </div>
          <p>{info ? `${info.os} / ${info.arch}` : "运行 pnpm dev"}</p>
          <p>{info?.goVersion}</p>
          <p>Author · HyperparameterSI</p>
        </div>
      </aside>
      <main className="studio-main">
        <header className="flex items-center justify-between gap-6 border-b pb-6">
          <span className="text-sm text-muted-foreground">
            Workspace <span className="mx-3">/</span> Components
          </span>
          <Button
            variant="ghost"
            size="icon"
            aria-label="切换明暗主题"
            onClick={() => setDark(!dark)}
          >
            {dark ? <Sun /> : <Moon />}
          </Button>
        </header>
        <section className="studio-hero">
          <div>
            <p className="text-xs font-medium tracking-[.25em] text-brand">
              BUILD SOMETHING NATIVE
            </p>
            <h1>
              每一个控件，
              <br />
              都能亲手试一试。
            </h1>
            <p className="max-w-xl text-sm leading-7 text-muted-foreground">
              从基础表单到代码编辑器，从系统 WebView 到 Chromium 浏览器。统一浏览、搜索并打开真实的
              Golet 组件演示。
            </p>
          </div>
          <div className="studio-hero-stat">
            <Layers size={24} />
            <strong>{catalog.length || "—"}</strong>
            <span>交互演示分组</span>
            <p>Go native · React · Plugins</p>
          </div>
        </section>
        <section className="mb-6 flex flex-wrap items-center justify-between gap-4 rounded-2xl border bg-card p-6">
          <div>
            <p className="text-xs font-medium tracking-widest text-muted-foreground">
              NATIVE / MATERIAL LAB
            </p>
            <h2 className="mt-2 text-xl font-semibold">原生主题实验室</h2>
            <p className="mt-2 text-sm text-muted-foreground">
              Maia 中性 · Liquid Glass · 明暗切换 · 真实 Go 控件
            </p>
          </div>
          <Button
            disabled={busy || !catalog.some((d) => d.id === "Themes")}
            onClick={() => {
              const demo = catalog.find((d) => d.id === "Themes");
              if (demo) void open(demo);
            }}
          >
            打开原生演示 <ArrowUpRight size={16} />
          </Button>
        </section>
        <div className="relative mb-6">
          <Search className="absolute left-3 top-3 text-muted-foreground" size={18} />
          <Input
            aria-label="搜索控件"
            className="pl-10"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="搜索组件，例如 Button、Table、Browser、Loading…"
          />
        </div>
        {error && (
          <p
            role="alert"
            className="mb-5 rounded-lg border border-destructive/30 p-4 text-sm text-destructive"
          >
            {error}
          </p>
        )}
        {notice && (
          <p role="status" className="mb-5 flex items-center gap-2 text-sm text-brand">
            <Check size={16} />
            {notice}
          </p>
        )}
        <div className="studio-explorer">
          <div className="studio-catalog">
            {visible.map((d) => (
              <button
                key={d.id}
                className={`studio-demo ${selected?.id === d.id ? "selected" : ""}`}
                onClick={() => setSelected(d)}
              >
                <span className="studio-demo-icon">
                  {d.id === "Browser" ? (
                    <Globe size={20} />
                  ) : d.category === "高级组件" ? (
                    <Code2 size={20} />
                  ) : (
                    <Box size={20} />
                  )}
                </span>
                <div>
                  <span className="text-[10px] text-muted-foreground">{d.category}</span>
                  <h2>{d.title}</h2>
                  <p>{d.description}</p>
                </div>
                <ArrowUpRight size={16} className="shrink-0 text-muted-foreground" />
              </button>
            ))}
            {!visible.length && (
              <p className="p-8 text-sm text-muted-foreground">
                {info ? "没有找到匹配组件，试试其他关键词。" : "等待连接 Go 后端…"}
              </p>
            )}
          </div>
          <div className="studio-detail">
            {selected && (
              <Card>
                <CardHeader>
                  <p className="mb-2 text-xs text-brand">{selected.category} / LIVE DEMO</p>
                  <CardTitle className="text-2xl">{selected.title}</CardTitle>
                  <CardDescription className="leading-6">{selected.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="mb-7 flex flex-wrap gap-2">
                    {selected.controls.map((control) => (
                      <span
                        key={control}
                        className="rounded-md border bg-muted/40 px-2 py-1 font-mono text-[11px]"
                      >
                        {control}
                      </span>
                    ))}
                  </div>
                  {selected.id === "Browser" ? (
                    <BrowserPanel />
                  ) : (
                    <>
                      <Button
                        className="w-full"
                        disabled={busy || !info}
                        onClick={() => open(selected)}
                      >
                        {busy ? <Loader2 className="animate-spin" /> : <ArrowUpRight />}打开原生演示
                      </Button>
                      <p className="mt-3 text-xs leading-6 text-muted-foreground">
                        在 Studio
                        进程内打开原生窗口。重复点击会聚焦已打开的演示；关闭窗口会释放对应资源。
                      </p>
                    </>
                  )}
                  <div className="mt-7 border-t pt-5">
                    <p className="mb-2 text-xs text-muted-foreground">实现位置</p>
                    <code className="break-all text-xs leading-6">{selected.source}</code>
                  </div>
                </CardContent>
              </Card>
            )}
            <Card className="mt-5">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-base">
                  <Cpu size={17} /> Web 组件与 Go 互通
                </CardTitle>
                <CardDescription>shadcn Button / Input / Card · typed IPC</CardDescription>
              </CardHeader>
              <CardContent>
                <form
                  className="space-y-3"
                  onSubmit={async (e) => {
                    e.preventDefault();
                    try {
                      setReply(await Greeter.greet(name));
                    } catch (e) {
                      setError(String(e));
                    }
                  }}
                >
                  <Input
                    aria-label="问候姓名"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    maxLength={200}
                  />
                  <Button variant="outline" disabled={!info} type="submit">
                    调用 Go 方法
                  </Button>
                  <p aria-live="polite" className="text-sm leading-6">
                    {reply}
                  </p>
                </form>
              </CardContent>
            </Card>
          </div>
        </div>
        <footer className="mt-10 border-t py-6 text-xs text-muted-foreground">
          Golet Studio · 原生控件使用真实 Go 实现 · Browser 与 Terminal 按需启动
        </footer>
      </main>
    </div>
  );
}

function BrowserPanel() {
  const [status, setStatus] = useState<BrowserStatus | null>(null);
  const [style, setStyle] = useState("alloy");
  const [preset, setPreset] = useState("Blue");
  const [dark, setDark] = useState(false);
  const [radius, setRadius] = useState(16);
  const [fontSize, setFontSize] = useState(14);
  const theme = { preset, dark, radius, fontSize };
  const canTheme = status?.themeSupported && style === "alloy";
  const [url, setURL] = useState("https://go.dev");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  async function refresh() {
    try {
      setStatus(await BrowserDemo.status());
    } catch (e) {
      setError(String(e));
    }
  }
  useEffect(() => {
    void refresh();
  }, []);
  async function run(action: string) {
    setBusy(true);
    setError("");
    try {
      if (action === "restart") {
        await BrowserDemo.command("close", "");
        await BrowserDemo.startWithTheme(style, url, theme);
      } else if (action === "start") await BrowserDemo.startWithTheme(style, url, theme);
      else if (action === "theme") await BrowserDemo.applyTheme(theme);
      else await BrowserDemo.command(action, url);
    } catch (e) {
      setError(String(e));
    } finally {
      await refresh();
      setBusy(false);
    }
  }
  return (
    <div className="space-y-4">
      <p className="rounded-lg bg-muted p-3 text-xs leading-6">
        Browser 是独立 CEF 窗口，不是嵌入当前页面的 iframe。原生宿主支持 macOS，Windows
        为实验性支持。Studio 使用临时演示会话，不访问钥匙串、不保留登录状态；请勿用于真实账号。
      </p>
      <label className="block text-sm">
        窗口风格
        <select
          aria-label="Browser 风格"
          className="mt-2 w-full rounded-md border bg-background p-2"
          value={style}
          disabled={busy}
          onChange={(e) => setStyle(e.target.value)}
        >
          <option value="alloy">Alloy · Golet 自定义界面</option>
          <option value="chrome">Chrome · 插件兼容界面</option>
        </select>
      </label>
      <p className="text-xs leading-6 text-muted-foreground">
        {status?.running
          ? "选择风格后，点击“重启并应用风格”。重启将关闭当前演示的所有标签，并打开下方网址；临时会话数据不会保留。"
          : "选择 Alloy 自定义界面或 Chrome 插件兼容界面，然后启动 Browser。"}
      </p>
      <fieldset
        disabled={busy || !canTheme}
        className="space-y-3 rounded-lg border p-4 disabled:opacity-50"
      >
        <legend className="px-2 text-sm font-medium">Forma Browser 主题</legend>
        <div className="grid grid-cols-2 gap-3">
          <label className="text-sm">
            配色
            <select
              className="mt-2 w-full rounded-md border bg-background p-2"
              value={preset}
              onChange={(e) => setPreset(e.target.value)}
            >
              {(status?.presets ?? ["Blue"]).map((name) => (
                <option key={name} value={name}>
                  {name}
                </option>
              ))}
            </select>
          </label>
          <label className="text-sm">
            外观
            <select
              className="mt-2 w-full rounded-md border bg-background p-2"
              value={dark ? "dark" : "light"}
              onChange={(e) => setDark(e.target.value === "dark")}
            >
              <option value="light">浅色</option>
              <option value="dark">深色</option>
            </select>
          </label>
          <label className="text-sm">
            圆角
            <select
              className="mt-2 w-full rounded-md border bg-background p-2"
              value={radius}
              onChange={(e) => setRadius(Number(e.target.value))}
            >
              <option value={16}>Maia · 16</option>
              <option value={8}>Soft · 8</option>
              <option value={0}>Sharp · 0</option>
            </select>
          </label>
          <label className="text-sm">
            字号
            <select
              className="mt-2 w-full rounded-md border bg-background p-2"
              value={fontSize}
              onChange={(e) => setFontSize(Number(e.target.value))}
            >
              <option value={14}>14</option>
              <option value={16}>16</option>
            </select>
          </label>
        </div>
        <Button
          variant="outline"
          disabled={!status?.running || status.style !== "alloy"}
          onClick={() => run("theme")}
        >
          应用主题 · 保留标签页
        </Button>
      </fieldset>
      <p className="text-xs leading-6 text-muted-foreground">
        {canTheme
          ? "启动时应用所选 Forma 主题；运行中可点击应用主题。需要重新构建最新 Chromium 宿主。主题仅影响原生浏览器外壳，网页和新标签页不随之改色。"
          : "Forma Browser 主题仅支持 macOS Alloy。Windows 和 Chrome 风格继续使用各自的默认外观。"}
      </p>
      <Input
        aria-label="浏览器网址"
        value={url}
        onChange={(e) => setURL(e.target.value)}
        placeholder="https://go.dev"
      />
      {!status?.installed && (
        <p className="text-xs leading-6 text-amber-700">
          尚未安装宿主，请先执行 <code>pnpm setup:chromium</code>。
        </p>
      )}
      <div className="flex flex-wrap gap-2">
        <Button
          disabled={busy || !status?.installed || !status?.supported}
          onClick={() => run(status?.running ? "restart" : "start")}
        >
          {busy ? <Loader2 className="animate-spin" /> : <Globe />}
          {status?.running ? "重启并应用风格" : "启动 Browser"}
        </Button>
        <Button variant="outline" disabled={busy} onClick={refresh}>
          检查状态
        </Button>
      </div>
      <div className="flex flex-wrap gap-2">
        {(
          [
            ["navigate", "打开网址"],
            ["newtab", "新建标签"],
            ["back", "后退"],
            ["forward", "前进"],
            ["reload", "刷新"],
            ["stop", "停止加载"],
            ["close", "关闭"],
          ] as const
        ).map(([action, label]) => (
          <Button
            key={action}
            variant="outline"
            size="sm"
            disabled={
              busy ||
              !status?.running ||
              (action === "back" && !status.state.canGoBack) ||
              (action === "forward" && !status.state.canGoForward)
            }
            onClick={() => run(action)}
          >
            {label}
          </Button>
        ))}
      </div>
      {error && (
        <p role="alert" className="text-sm text-destructive">
          {error}
        </p>
      )}
      <div className="break-all rounded-lg border p-3 text-xs leading-6">
        <p>
          {status?.running ? "● 运行中" : "○ 未运行"} · {status?.state.tabs?.length ?? 0} 个标签
        </p>
        {status?.running && (
          <>
            <p>{status.state.title}</p>
            <p>{status.state.url}</p>
          </>
        )}
      </div>
      <pre className="overflow-x-auto rounded-lg bg-muted p-3 text-[11px] leading-5">{`chromium.Launch(ctx, chromium.Options{\n  Style: chromium.${style === "chrome" ? "StyleChrome" : "StyleAlloy"},\n  URL: "https://go.dev",\n})`}</pre>
    </div>
  );
}
