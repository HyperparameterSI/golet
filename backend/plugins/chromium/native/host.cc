// Author: HyperparameterSI
#include "host.h"
#include "shell.h"
#include <iostream>
#include <map>
#include <thread>
#if defined(OS_WIN)
#include <process.h>
#else
#include <unistd.h>
#endif
#include "include/base/cef_callback.h"
#include "include/cef_browser.h"
#include "include/cef_devtools_message_observer.h"
#include "include/cef_command_line.h"
#include "include/cef_parser.h"
#include "include/wrapper/cef_closure_task.h"
#include "include/wrapper/cef_helpers.h"

namespace {
// Only the CEF UI thread touches this object or writes protocol responses.
class Host final : public CefApp, public CefBrowserProcessHandler,
                   public CefClient, public CefLifeSpanHandler,
                   public CefDisplayHandler, public CefLoadHandler,
                   public CefFocusHandler, public CefDevToolsMessageObserver {
 public:
  CefRefPtr<CefBrowserProcessHandler> GetBrowserProcessHandler() override { return this; }
  CefRefPtr<CefClient> GetDefaultClient() override { return this; }
  CefRefPtr<CefLifeSpanHandler> GetLifeSpanHandler() override { return this; }
  CefRefPtr<CefDisplayHandler> GetDisplayHandler() override { return this; }
  CefRefPtr<CefLoadHandler> GetLoadHandler() override { return this; }
  CefRefPtr<CefFocusHandler> GetFocusHandler() override { return this; }

  void OnContextInitialized() override {
    CEF_REQUIRE_UI_THREAD();
    chrome_style_ = CefCommandLine::GetGlobalCommandLine()->HasSwitch("golet-chrome-style");
    app_mode_ = CefCommandLine::GetGlobalCommandLine()->HasSwitch("golet-app-mode");
    if (!chrome_style_) CreateBrowserShell();
    auto url = CefCommandLine::GetGlobalCommandLine()->GetSwitchValue("url");
    NewTab(url.ToString());
    // Blocking reads use a dedicated thread; there is no polling timer. Parsing
    // and browser operations are posted to CEF's UI thread.
    CefRefPtr<Host> self = this;
    std::thread([self]() {
      std::string line;
      while (std::getline(std::cin, line)) {
        if (line.size() > 65536) break;
        CefPostTask(TID_UI, base::BindOnce(&Host::Command, self, line));
      }
      CefPostTask(TID_UI, base::BindOnce(&Host::CloseAll, self));
    }).detach();
  }

  void OnAfterCreated(CefRefPtr<CefBrowser> browser) override {
    CEF_REQUIRE_UI_THREAD();
    browsers_[browser->GetIdentifier()] = browser;
    // Author: HyperparameterSI. Make rendered page controls available to AX
    // when the embedding app opts into native computer-use support.
    if (CefCommandLine::GetGlobalCommandLine()->HasSwitch("golet-computer-use"))
      browser->GetHost()->SetAccessibilityState(STATE_ENABLED);
    observers_[browser->GetIdentifier()] = browser->GetHost()->AddDevToolsMessageObserver(this);
    if (!chrome_style_) ShellAttachBrowser(browser->GetIdentifier(), browser->GetHost()->GetWindowHandle());
    active_ = browser;
    if (!chrome_style_) ShellSelectBrowser(browser->GetIdentifier());
    Publish(ready_ ? "state" : "ready");
    ready_ = true;
  }
  bool DoClose(CefRefPtr<CefBrowser> browser) override {
    if (chrome_style_) return false;
    // A tab closes its child view, not the shared top-level browser window.
    ShellDetachBrowser(browser->GetIdentifier());
    return true;
  }
  bool OnBeforePopup(CefRefPtr<CefBrowser> browser, CefRefPtr<CefFrame> frame,
                     int popup_id, const CefString& target_url,
                     const CefString& target_frame_name, WindowOpenDisposition disposition,
                     bool user_gesture, const CefPopupFeatures& features, CefWindowInfo& info,
                     CefRefPtr<CefClient>& client, CefBrowserSettings& settings,
                     CefRefPtr<CefDictionaryValue>& extra_info, bool* no_javascript_access) override {
    if (app_mode_) return true;
    if (chrome_style_) info.runtime_style = CEF_RUNTIME_STYLE_CHROME;
    else ConfigureShellChild(info);
    return false;
  }
  void OnBeforeClose(CefRefPtr<CefBrowser> browser) override {
    observers_.erase(browser->GetIdentifier());
    for (auto it = pending_.begin(); it != pending_.end();) {
      if (it->second.first == browser->GetIdentifier()) {
        auto reply = it->second.second; reply->SetString("error", "tab closed"); Write(reply); it = pending_.erase(it);
      } else ++it;
    }
    if (!chrome_style_) ShellRemoveBrowser(browser->GetIdentifier());
    browsers_.erase(browser->GetIdentifier());
    titles_.erase(browser->GetIdentifier());
    if (active_ && active_->IsSame(browser)) {
      active_ = browsers_.empty() ? nullptr : browsers_.begin()->second;
    }
    if (browsers_.empty()) { ShellDestroy(); CefQuitMessageLoop(); }
    else { ShellSelectBrowser(active_->GetIdentifier()); Publish("state"); }
  }
  void OnGotFocus(CefRefPtr<CefBrowser> browser) override {
    if (active_ && active_->IsSame(browser)) return;
    active_ = browser; ShellSelectBrowser(browser->GetIdentifier()); Publish("state");
  }
  void OnTitleChange(CefRefPtr<CefBrowser> browser, const CefString& title) override {
    titles_[browser->GetIdentifier()] = title;
    UpdateShell(browser);
    if (active_ && active_->IsSame(browser)) Publish("state");
  }
  void OnAddressChange(CefRefPtr<CefBrowser> browser, CefRefPtr<CefFrame> frame,
                       const CefString& url) override {
    if (frame->IsMain()) UpdateShell(browser);
    if (frame->IsMain() && active_ && active_->IsSame(browser)) Publish("state");
  }
  void OnLoadingStateChange(CefRefPtr<CefBrowser> browser, bool loading,
                            bool back, bool forward) override {
    if (active_ && active_->IsSame(browser)) Publish("state");
  }

  void NewTab(const std::string& target = "") {
    if (app_mode_ && tab_requested_) return;
    tab_requested_ = true;
    CefWindowInfo window;
#if defined(OS_WIN)
    if (chrome_style_) window.SetAsPopup(nullptr, "Golet Browser");
#endif
    if (chrome_style_) window.runtime_style = CEF_RUNTIME_STYLE_CHROME;
    else ConfigureShellChild(window);
    CefBrowserSettings settings;
    const auto url = target.empty() || target == "about:blank" ? (app_mode_ ? "about:blank" : ShellHomeURL()) : target;
    CefBrowserHost::CreateBrowser(window, this, url, settings, nullptr, nullptr);
  }
  void Action(const std::string& method, int id, const std::string& value) {
    CEF_REQUIRE_UI_THREAD();
    if (method == "assistant") {
      auto url = CefCommandLine::GetGlobalCommandLine()->GetSwitchValue("assistant-url");
      if (!url.empty()) NewTab(url.ToString());
      return;
    }
    if (method == "close") { CloseAll(); return; }
    if (method == "newtab") { NewTab(value); return; }
    if (method == "select-tab") {
      auto found = browsers_.find(id);
      if (found != browsers_.end()) {
        active_ = found->second; ShellSelectBrowser(id);
        active_->GetHost()->SetFocus(true); Publish("state");
      }
      return;
    }
    if (method == "close-tab") {
      auto found = browsers_.find(id);
      if (found != browsers_.end()) found->second->GetHost()->CloseBrowser(false);
      return;
    }
    if (!active_) return;
    if (method == "home") active_->GetMainFrame()->LoadURL(ShellHomeURL());
    else if (method == "navigate") active_->GetMainFrame()->LoadURL(value);
    else if (method == "back") active_->GoBack();
    else if (method == "forward") active_->GoForward();
    else if (method == "reload") active_->Reload();
    else if (method == "stop") active_->StopLoad();
  }
  void CloseAll() {
    CEF_REQUIRE_UI_THREAD();
    // CloseBrowser may eventually remove entries; hold a stable copy.
    auto browsers = browsers_;
    for (auto& entry : browsers) entry.second->GetHost()->CloseBrowser(true);
  }

 private:
  void UpdateShell(CefRefPtr<CefBrowser> browser) {
    if (chrome_style_) return;
    ShellUpdateBrowser(browser->GetIdentifier(), titles_[browser->GetIdentifier()].ToString(),
                       browser->GetMainFrame()->GetURL().ToString(), browser->IsLoading(),
                       browser->CanGoBack(), browser->CanGoForward());
  }
  void Write(CefRefPtr<CefDictionaryValue> dict) {
    auto value = CefValue::Create(); value->SetDictionary(dict);
    std::cout << CefWriteJSON(value, JSON_WRITER_DEFAULT).ToString() << std::endl;
  }
  void Publish(const std::string& event) {
    if (!active_) return;
    UpdateShell(active_);
    auto state = CefDictionaryValue::Create();
    state->SetInt("browserId", active_->GetIdentifier());
    state->SetString("url", active_->GetMainFrame()->GetURL());
    state->SetString("title", titles_[active_->GetIdentifier()]);
    state->SetBool("loading", active_->IsLoading());
    state->SetBool("canGoBack", active_->CanGoBack());
    state->SetBool("canGoForward", active_->CanGoForward());
    auto tabs = CefListValue::Create();
    size_t index = 0;
    for (const auto& [id, browser] : browsers_) {
      auto tab = CefDictionaryValue::Create();
      tab->SetInt("browserId", id);
      tab->SetString("url", browser->GetMainFrame()->GetURL());
      tab->SetString("title", titles_[id]);
      tabs->SetDictionary(index++, tab);
    }
    state->SetList("tabs", tabs);
    auto msg = CefDictionaryValue::Create();
    msg->SetString("event", event); msg->SetDictionary("state", state); Write(msg);
  }
  void Command(std::string line) {
    CEF_REQUIRE_UI_THREAD();
    auto value = CefParseJSON(line, JSON_PARSER_RFC);
    if (!value || value->GetType() != VTYPE_DICTIONARY) return;
    auto request = value->GetDictionary();
    auto reply = CefDictionaryValue::Create();
    reply->SetValue("id", request->GetValue("id"));
    const auto method = request->GetString("method").ToString();
    if (method == "close") { Write(reply); CloseAll(); return; }
    if (method == "set-theme") {
      if (chrome_style_ || !ShellSetTheme(request->GetString("url").ToString()))
        reply->SetString("error", "custom theme unsupported or invalid; requires macOS Alloy");
      Write(reply); return;
    }
    if (method == "newtab") {
      if (app_mode_) reply->SetString("error", "new tabs are unavailable in AppMode");
      else NewTab(request->GetString("url").ToString());
      Write(reply); return;
    }
    if (method == "select-tab" || method == "close-tab") {
      const int id = request->GetInt("browserId");
      if (!browsers_.count(id)) reply->SetString("error", "tab not found");
      else Action(method, id, "");
      Write(reply); return;
    }
    // Author: HyperparameterSI. Read consent and identity atomically on the UI
    // thread. No protocol command can enable the native consent control.
    if (method == "computer-use-target") {
      if (!active_) { reply->SetString("error", "no active browser"); Write(reply); return; }
      auto result = CefDictionaryValue::Create();
#if defined(OS_WIN)
      result->SetInt("pid", _getpid());
#else
      result->SetInt("pid", getpid());
#endif
      const int window = ShellNativeWindowID(active_->GetHost()->GetWindowHandle());
      result->SetInt("windowId", window);
      result->SetInt("browserId", active_->GetIdentifier());
      result->SetBool("computerUse", window > 0 && !chrome_style_ && !app_mode_ && ShellComputerUseEnabled());
      reply->SetDictionary("result", result); Write(reply); return;
    }
    if (method == "devtools") {
      auto it = browsers_.find(request->GetInt("browserId"));
      if (it == browsers_.end()) { reply->SetString("error", "tab not found"); Write(reply); return; }
      if (pending_.size() >= 128) { reply->SetString("error", "too many pending DevTools commands"); Write(reply); return; }
      int id = ++devtools_id_;
      pending_[id] = {it->first, reply};
      if (!it->second->GetHost()->ExecuteDevToolsMethod(id, request->GetString("devToolsMethod"), request->GetDictionary("params"))) {
        pending_.erase(id); reply->SetString("error", "DevTools rejected command"); Write(reply);
      }
      return;
    }
    if (!active_) { reply->SetString("error", "no active browser"); Write(reply); return; }
    if (method == "navigate") {
      const auto url = request->GetString("url").ToString();
      CefURLParts parts;
      if (url != "about:blank" && (!CefParseURL(url, parts) ||
          (CefString(&parts.scheme) != "https" && CefString(&parts.scheme) != "http"))) {
        reply->SetString("error", "unsupported URL");
      } else { active_->GetMainFrame()->LoadURL(url); }
    } else if (method == "back") active_->GoBack();
    else if (method == "forward") active_->GoForward();
    else if (method == "reload") active_->Reload();
    else if (method == "stop") active_->StopLoad();
    else reply->SetString("error", "unknown command");
    Write(reply);
  }
  void OnDevToolsMethodResult(CefRefPtr<CefBrowser> browser, int id, bool success,
      const void* data, size_t size) override {
    auto it = pending_.find(id);
    if (it == pending_.end()) return;
    auto reply = it->second.second; pending_.erase(it);
    if (size > 900 * 1024) reply->SetString("error", "DevTools result exceeds 900 KiB");
    else if (!success) reply->SetString("error", std::string(static_cast<const char*>(data), size));
    else {
      auto result = CefParseJSON(std::string(static_cast<const char*>(data), size), JSON_PARSER_RFC);
      if (result) reply->SetValue("result", result);
    }
    Write(reply);
  }
  bool chrome_style_ = false;
  bool app_mode_ = false;
  bool tab_requested_ = false;
  int devtools_id_ = 0;
  std::map<int, CefRefPtr<CefRegistration>> observers_;
  std::map<int, std::pair<int, CefRefPtr<CefDictionaryValue>>> pending_;
  bool ready_ = false;
  CefRefPtr<CefBrowser> active_;
  std::map<int, CefRefPtr<CefBrowser>> browsers_;
  std::map<int, CefString> titles_;
  IMPLEMENT_REFCOUNTING(Host);
};
CefRefPtr<Host> host;
}

CefRefPtr<CefApp> CreateGoletChromiumApp() { host = new Host; return host; }
void CloseGoletChromium() { if (host) host->CloseAll(); }
void ConfigureGoletChromium(CefSettings& settings, CefRefPtr<CefCommandLine> args) {
  CefString(&settings.cache_path) = args->GetSwitchValue("cache-path");
  CefString(&settings.root_cache_path) = args->GetSwitchValue("cache-path");
  settings.log_severity = LOGSEVERITY_WARNING;
}

void BrowserUIAction(const std::string& method, int id, const std::string& value) {
  if (host) host->Action(method, id, value);
}
