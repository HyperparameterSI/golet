// Author: HyperparameterSI

// Greeter and events are generated from the Go code by `golet generate`.
import { Greeter, events } from "./golet";

const $ = <T extends HTMLElement>(selector: string) => document.querySelector<T>(selector)!;

$("#greet").addEventListener("submit", async (e) => {
  e.preventDefault();
  $("#greeting").textContent = await Greeter.greet($<HTMLInputElement>("#name").value);
});

const info = await Greeter.info();
$("#info").textContent = `${info.os}/${info.arch} · ${info.goVersion}`;

events.tick.on((time) => {
  $("#clock").textContent = new Date(time).toLocaleTimeString();
});
