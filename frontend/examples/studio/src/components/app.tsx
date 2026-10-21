// Author: HyperparameterSI

import { useEffect, useState, type FormEvent } from "react";
import { ArrowUpRight, Check, Code2, Cpu, Layers, Loader2, Moon, Sun } from "lucide-react";
import { Greeter } from "@/golet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";

export function App() {
  const [dark, setDark] = useState(false);
  const [name, setName] = useState("");
  const [reply, setReply] = useState("");
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);
  const [info, setInfo] = useState<Awaited<ReturnType<typeof Greeter.info>> | null>(null);

  useEffect(() => {
    setDark(matchMedia("(prefers-color-scheme: dark)").matches);
  }, []);
  useEffect(() => {
    document.documentElement.classList.toggle("dark", dark);
  }, [dark]);
  useEffect(() => {
    let active = true;
    Greeter.info()
      .then((v) => {
        if (active) setInfo(v);
      })
      .catch(() => {
        if (active) setError("Open this app with golet dev to connect to the Go backend.");
      });
    return () => {
      active = false;
    };
  }, []);

  async function greet(event: FormEvent) {
    event.preventDefault();
    if (pending) return;
    setPending(true);
    setError("");
    setReply("");
    try {
      setReply(await Greeter.greet(name.trim()));
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setPending(false);
    }
  }

  return (
    <div className="mx-auto min-h-screen max-w-5xl px-6 md:px-12">
      <header className="flex h-20 items-center justify-between border-b">
        <div className="flex items-center gap-3">
          <span className="flex size-9 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <Layers size={19} />
          </span>
          <span className="text-lg font-semibold tracking-tight">
            golet
            <span className="ml-3 rounded-md border px-2 py-0.5 text-xs font-normal text-muted-foreground">
              workspace
            </span>
          </span>
        </div>
        <Button
          variant="ghost"
          size="icon"
          aria-label={dark ? "Use light theme" : "Use dark theme"}
          onClick={() => setDark(!dark)}
        >
          {dark ? <Sun /> : <Moon />}
        </Button>
      </header>
      <main className="py-12 md:py-16">
        <div className="mb-10">
          <p className="mb-3 text-xs font-medium uppercase tracking-[0.2em] text-muted-foreground">
            Your next desktop app
          </p>
          <h1 className="text-4xl font-semibold tracking-tight md:text-5xl">
            Native power.
            <br />A familiar canvas.
          </h1>
          <p className="mt-5 max-w-lg text-base leading-7 text-muted-foreground">
            Build with Go. Create with React. A lightweight desktop foundation, ready for your
            ideas.
          </p>
        </div>
        <div className="grid items-start gap-6 md:grid-cols-[1.4fr_1fr]">
          <Card>
            <CardHeader>
              <CardTitle>Make the connection</CardTitle>
              <CardDescription>
                Call a Go method through your typed TypeScript client.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={greet} className="space-y-4">
                <label htmlFor="name" className="text-sm font-medium">
                  Your name
                </label>
                <Input
                  id="name"
                  autoComplete="given-name"
                  placeholder="Ada Lovelace"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  maxLength={200}
                />
                <Button type="submit" disabled={pending || !info} className="w-full">
                  {pending ? <Loader2 className="animate-spin" /> : <ArrowUpRight />}{" "}
                  {pending ? "Calling Go…" : "Say hello from Go"}
                </Button>
              </form>
              <div aria-live="polite" className="mt-5 min-h-12 text-sm">
                {reply && (
                  <p className="flex items-start gap-2 rounded-md bg-muted p-3">
                    <Check className="mt-0.5 size-4 shrink-0" />
                    {reply}
                  </p>
                )}
                {error && (
                  <p role="alert" className="text-destructive">
                    {error}
                  </p>
                )}
                {!reply && !error && (
                  <p className="text-muted-foreground">Your response will appear here.</p>
                )}
              </div>
            </CardContent>
          </Card>
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Cpu size={18} /> Runtime
                </CardTitle>
                <CardDescription>
                  {info ? "Connected to your desktop backend" : "Waiting for desktop connection"}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <dl className="space-y-3 text-sm">
                  {[
                    ["Platform", info ? `${info.os} / ${info.arch}` : "—"],
                    ["Go version", info?.goVersion ?? "—"],
                    ["Frontend", "TanStack Start + shadcn"],
                  ].map(([key, value]) => (
                    <div key={key} className="flex justify-between gap-4">
                      <dt className="text-muted-foreground">{key}</dt>
                      <dd className="font-mono text-xs">{value}</dd>
                    </div>
                  ))}
                </dl>
              </CardContent>
            </Card>
            <div className="flex gap-3 px-1 text-sm text-muted-foreground">
              <Code2 className="mt-0.5 size-5 shrink-0" />
              <p className="leading-6">
                Start in <code className="text-foreground">src/components/app.tsx</code>. Your
                components live in your project. Make them yours.
              </p>
            </div>
          </div>
        </div>
      </main>
      <footer className="flex flex-wrap justify-between gap-3 border-t py-6 text-xs text-muted-foreground">
        <span>Golet · Go beyond the browser.</span>
        <span>System WebView · Typed IPC · No bundled browser</span>
      </footer>
    </div>
  );
}
