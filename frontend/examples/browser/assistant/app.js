// Author: HyperparameterSI
const token = location.hash.slice(1) || sessionStorage.getItem("golet-assistant-token") || "";
if (token) sessionStorage.setItem("golet-assistant-token", token);
history.replaceState(null, "", location.pathname);
const $ = (id) => document.getElementById(id);
async function api(path, data) {
  const response = await fetch(`/api/${path}`, {
    method: data ? "POST" : "GET",
    headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
    body: data ? JSON.stringify(data) : undefined,
  });
  const result = await response.json();
  if (!response.ok) throw Error(result.error || "连接失败，请重新打开助手");
  return result;
}
function message(text, kind = "") {
  const node = document.createElement("div");
  node.className = `message ${kind}`;
  node.textContent = text;
  $("conversation").append(node);
  node.scrollIntoView({ block: "nearest" });
  return node;
}
function status(error) {
  $("status").textContent = error?.message || error || "";
}
async function refresh() {
  try {
    const data = await api("state");
    $("model").textContent = data.configured ? data.model : "模型尚未配置";
    const selected = $("tabs").value;
    $("tabs").replaceChildren();
    for (const tab of data.state.tabs || []) {
      if (!/^https?:/.test(tab.url) || tab.url.startsWith(location.origin)) continue;
      const option = document.createElement("option");
      option.value = tab.browserId;
      option.textContent = `${tab.title || tab.url} — ${tab.url}`;
      $("tabs").append(option);
    }
    if ([...$("tabs").options].some((o) => o.value === selected)) $("tabs").value = selected;
    $("extension-list").replaceChildren();
    for (const [index, extension] of (data.extensions || []).entries()) {
      const row = document.createElement("div");
      row.className = "extension";
      const title = document.createElement("strong");
      title.textContent = extension.name;
      const desc = document.createElement("p");
      desc.textContent = `${extension.description} · ${extension.origins.join(", ")}`;
      const run = document.createElement("button");
      run.textContent = "在所选页面运行";
      run.disabled = !!data.disabled[index];
      run.onclick = async () => {
        if (!confirm(`运行可信扩展「${extension.name}」？它可以读取并修改所选页面。`)) return;
        try {
          await api("extension/run", { index, tab: Number($("tabs").value) });
          status("扩展已运行");
        } catch (e) {
          status(e);
        }
      };
      const toggle = document.createElement("button");
      toggle.textContent = data.disabled[index] ? "启用" : "停用";
      toggle.onclick = async () => {
        try {
          await api("extension/toggle", { index, enabled: !!data.disabled[index] });
          await refresh();
        } catch (e) {
          status(e);
        }
      };
      row.append(title, desc, run, toggle);
      $("extension-list").append(row);
    }
    if (!data.extensions?.length)
      $("extension-list").textContent = "尚未加载扩展。启动时添加 -user-extension <目录>。";
  } catch (e) {
    status(e);
  }
}
$("refresh").onclick = refresh;
$("chat").onsubmit = async (event) => {
  event.preventDefault();
  const question = $("question").value.trim();
  if (!question) return;
  $("send").disabled = true;
  status("正在思考…");
  document.querySelector(".welcome")?.remove();
  message(question, "user");
  try {
    const data = await api("ask", {
      question,
      share: $("share").checked,
      tab: Number($("tabs").value),
    });
    message(data.answer.reply);
    for (const [index, action] of (data.answer.actions || []).entries()) {
      const row = document.createElement("div");
      row.className = "action";
      const text = document.createElement("code");
      const label = data.page?.elements?.[action.element]?.label || "页面";
      text.textContent = `${data.page?.title || ""} · ${data.page?.url || ""}\n${label}\n${JSON.stringify(action)}`;
      const button = document.createElement("button");
      button.textContent = "确认执行这一步";
      button.onclick = async () => {
        if (
          !confirm(
            `确认在原目标页面执行：\n${JSON.stringify(action)}\n执行一步后请重新提问，以获取最新页面状态。`,
          )
        )
          return;
        button.disabled = true;
        try {
          await api("execute", { plan: data.plan, index });
          status("已执行。请重新提问获取最新页面状态。");
        } catch (e) {
          status(e);
        }
      };
      row.append(text, button);
      $("conversation").append(row);
    }
    status("");
    $("question").value = "";
  } catch (e) {
    status(e);
  } finally {
    $("send").disabled = false;
  }
};
refresh();
