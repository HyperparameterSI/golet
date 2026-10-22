// Author: HyperparameterSI
#include "shell.h"
#include <windows.h>
#include <commctrl.h>
#include <shlwapi.h>
#include <algorithm>
#include <map>
#include <vector>
#include "include/cef_command_line.h"
#include "include/cef_app.h"
#include "include/cef_parser.h"
#include "include/cef_version.h"

namespace {
HWND window = nullptr, address = nullptr, tabs = nullptr;
HWND back_button = nullptr, forward_button = nullptr;
std::map<int, HWND> pages;
std::vector<int> tab_ids;
int selected = 0;
bool app_mode = false;
constexpr int kBack = 101, kForward = 102, kReload = 103, kHome = 104,
              kNewTab = 105, kCloseTab = 106, kGo = 107, kAssistant = 108;
int ContentTop() { return app_mode ? 0 : 82; }
std::wstring Wide(const std::string& text) { return CefString(text).ToWString(); }
std::wstring ProductName() {
  auto name = CefCommandLine::GetGlobalCommandLine()->GetSwitchValue("golet-app-name");
  return name.empty() ? L"Golet Browser" : name.ToWString();
}
void Resize() {
  if (!window) return;
  RECT r; GetClientRect(window, &r);
  if (!app_mode) {
    MoveWindow(tabs, 8, 4, std::max(100L, r.right - 16), 30, TRUE);
    MoveWindow(address, 172, 42, std::max(80L, r.right - 414), 28, TRUE);
    const int ids[] = {kGo, kNewTab, kCloseTab, kAssistant};
    const int offsets[] = {234, 190, 152, 72};
    const int widths[] = {40, 34, 76, 64};
    for (int i = 0; i < 4; ++i)
      MoveWindow(GetDlgItem(window, ids[i]), r.right - offsets[i], 42, widths[i], 28, TRUE);
  }
  for (const auto& [id, view] : pages)
    if (view) MoveWindow(view, 0, ContentTop(), r.right,
                         std::max(0L, r.bottom - ContentTop()), TRUE);
}
void Navigate() {
  int size = GetWindowTextLengthW(address);
  std::wstring text(size + 1, L'\0');
  GetWindowTextW(address, text.data(), size + 1);
  text.resize(size);
  if (text.empty()) return;
  auto value = CefString(text).ToString();
  if (value.find(":") == std::string::npos) value = "https://" + value;
  BrowserUIAction("navigate", 0, value);
}
LRESULT CALLBACK AddressProc(HWND hwnd, UINT msg, WPARAM wp, LPARAM lp,
                              UINT_PTR, DWORD_PTR) {
  if (msg == WM_KEYDOWN && wp == VK_RETURN) { Navigate(); return 0; }
  return DefSubclassProc(hwnd, msg, wp, lp);
}
LRESULT CALLBACK WindowProc(HWND hwnd, UINT msg, WPARAM wp, LPARAM lp) {
  switch (msg) {
    case WM_SIZE: Resize(); return 0;
    case WM_CLOSE: BrowserUIAction("close"); return 0;
    case WM_SETFOCUS:
      if (pages.count(selected)) SetFocus(pages[selected]);
      return 0;
    case WM_GETMINMAXINFO: {
      auto info = reinterpret_cast<MINMAXINFO*>(lp);
      info->ptMinTrackSize = {640, 400};
      return 0;
    }
    case WM_NOTIFY:
      if (reinterpret_cast<NMHDR*>(lp)->hwndFrom == tabs &&
          reinterpret_cast<NMHDR*>(lp)->code == TCN_SELCHANGE) {
        int at = TabCtrl_GetCurSel(tabs);
        if (at >= 0 && at < static_cast<int>(tab_ids.size()))
          BrowserUIAction("select-tab", tab_ids[at]);
      }
      break;
    case WM_COMMAND:
      switch (LOWORD(wp)) {
        case kBack: BrowserUIAction("back"); break;
        case kForward: BrowserUIAction("forward"); break;
        case kReload: BrowserUIAction("reload"); break;
        case kHome: BrowserUIAction("home"); break;
        case kNewTab: BrowserUIAction("newtab"); break;
        case kCloseTab: BrowserUIAction("close-tab", selected); break;
        case kGo: Navigate(); break;
        case kAssistant: BrowserUIAction("assistant"); break;
      }
      return 0;
  }
  return DefWindowProcW(hwnd, msg, wp, lp);
}
HWND Control(const wchar_t* type, const wchar_t* title, DWORD style,
             int x, int y, int width, int id) {
  auto hwnd = CreateWindowExW(0, type, title, WS_CHILD | WS_VISIBLE | style,
      x, y, width, 28, window, reinterpret_cast<HMENU>(static_cast<INT_PTR>(id)),
      GetModuleHandleW(nullptr), nullptr);
  SendMessageW(hwnd, WM_SETFONT, reinterpret_cast<WPARAM>(GetStockObject(DEFAULT_GUI_FONT)), TRUE);
  return hwnd;
}
}

void CreateBrowserShell() {
  app_mode = CefCommandLine::GetGlobalCommandLine()->HasSwitch("golet-app-mode");
  INITCOMMONCONTROLSEX controls{sizeof(controls), ICC_TAB_CLASSES};
  InitCommonControlsEx(&controls);
  WNDCLASSW cls{};
  cls.lpfnWndProc = WindowProc;
  cls.hInstance = GetModuleHandleW(nullptr);
  cls.hCursor = LoadCursor(nullptr, IDC_ARROW);
  cls.hbrBackground = reinterpret_cast<HBRUSH>(COLOR_WINDOW + 1);
  cls.lpszClassName = L"GoletBrowserShell";
  RegisterClassW(&cls);
  window = CreateWindowExW(0, cls.lpszClassName, ProductName().c_str(),
      WS_OVERLAPPEDWINDOW | WS_CLIPCHILDREN, CW_USEDEFAULT, CW_USEDEFAULT,
      1200, 820, nullptr, nullptr, cls.hInstance, nullptr);
  if (!window) { CefQuitMessageLoop(); return; }
  if (!app_mode) {
    tabs = Control(WC_TABCONTROLW, L"", WS_TABSTOP, 8, 4, 1000, 0);
    back_button = Control(L"BUTTON", L"<", WS_TABSTOP, 8, 42, 34, kBack);
    forward_button = Control(L"BUTTON", L">", WS_TABSTOP, 48, 42, 34, kForward);
    Control(L"BUTTON", L"Reload", WS_TABSTOP, 88, 42, 76, kReload);
    address = Control(L"EDIT", L"", WS_BORDER | ES_AUTOHSCROLL | WS_TABSTOP, 172, 42, 700, 0);
    SetWindowSubclass(address, AddressProc, 1, 0);
    Control(L"BUTTON", L"Go", WS_TABSTOP, 900, 42, 40, kGo);
    Control(L"BUTTON", L"+", WS_TABSTOP, 944, 42, 34, kNewTab);
    Control(L"BUTTON", L"Close tab", WS_TABSTOP, 982, 42, 76, kCloseTab);
    auto ai = Control(L"BUTTON", L"AI", WS_TABSTOP, 1062, 42, 40, kAssistant);
    EnableWindow(ai, CefCommandLine::GetGlobalCommandLine()->HasSwitch("assistant-url"));
  }
  Resize();
  ShowWindow(window, SW_SHOW);
}
void ConfigureShellChild(CefWindowInfo& info) {
  RECT r{};
  if (window) GetClientRect(window, &r);
  info.SetAsChild(window, CefRect(0, ContentTop(), r.right, std::max(0L, r.bottom - ContentTop())));
  info.runtime_style = CEF_RUNTIME_STYLE_ALLOY;
}
void ShellAttachBrowser(int id, CefWindowHandle view) {
  if (!window) return;
  pages[id] = view;
  tab_ids.push_back(id);
  if (tabs) {
    TCITEMW item{}; item.mask = TCIF_TEXT; item.pszText = const_cast<wchar_t*>(L"New tab");
    TabCtrl_InsertItem(tabs, static_cast<int>(tab_ids.size()) - 1, &item);
  }
  Resize();
}
void ShellDetachBrowser(int id) {
  auto it = pages.find(id);
  if (it != pages.end() && it->second) { ShowWindow(it->second, SW_HIDE); DestroyWindow(it->second); it->second = nullptr; }
}
void ShellRemoveBrowser(int id) {
  pages.erase(id);
  auto it = std::find(tab_ids.begin(), tab_ids.end(), id);
  if (it != tab_ids.end()) {
    if (tabs) TabCtrl_DeleteItem(tabs, static_cast<int>(it - tab_ids.begin()));
    tab_ids.erase(it);
  }
}
void ShellSelectBrowser(int id) {
  selected = id;
  for (const auto& [key, view] : pages) if (view) ShowWindow(view, key == id ? SW_SHOW : SW_HIDE);
  auto it = std::find(tab_ids.begin(), tab_ids.end(), id);
  if (tabs && it != tab_ids.end()) TabCtrl_SetCurSel(tabs, static_cast<int>(it - tab_ids.begin()));
}
void ShellUpdateBrowser(int id, const std::string& title, const std::string& url,
                        bool, bool back, bool forward) {
  auto it = std::find(tab_ids.begin(), tab_ids.end(), id);
  auto text = Wide(title.empty() ? "New tab" : title);
  if (tabs && it != tab_ids.end()) {
    TCITEMW item{}; item.mask = TCIF_TEXT; item.pszText = text.data();
    TabCtrl_SetItem(tabs, static_cast<int>(it - tab_ids.begin()), &item);
  }
  if (id != selected || !window) return;
  SetWindowTextW(window, (text + L" - " + ProductName()).c_str());
  if (address && GetFocus() != address) SetWindowTextW(address, Wide(url).c_str());
  if (back_button) EnableWindow(back_button, back);
  if (forward_button) EnableWindow(forward_button, forward);
}
void ShellDestroy() {
  if (window) DestroyWindow(window);
  window = address = tabs = back_button = forward_button = nullptr;
  pages.clear(); tab_ids.clear(); selected = 0;
}
std::string ShellHomeURL() {
  wchar_t module[32768];
  DWORD length = GetModuleFileNameW(nullptr, module, 32768);
  if (!length || length >= 32768) return "about:blank";
  std::wstring path(module, length);
  path.resize(path.find_last_of(L"\\/") + 1);
  path += L"newtab\\index.html";
  wchar_t url[32768]; DWORD size = 32768;
  if (FAILED(UrlCreateFromPathW(path.c_str(), url, &size, 0))) return "about:blank";
  auto result = CefString(url).ToString() + "?version=" +
      std::to_string(CHROME_VERSION_MAJOR) + "." + std::to_string(CHROME_VERSION_MINOR) + "." +
      std::to_string(CHROME_VERSION_BUILD) + "." + std::to_string(CHROME_VERSION_PATCH);
  auto name = CefCommandLine::GetGlobalCommandLine()->GetSwitchValue("golet-app-name");
  if (!name.empty()) result += "&appName=" + CefURIEncode(name, false).ToString();
  return result;
}

// Author: HyperparameterSI. The current native computer-use driver is macOS-only.
int ShellNativeWindowID(CefWindowHandle) { return 0; }
bool ShellComputerUseEnabled() { return false; }

// Author: HyperparameterSI. Windows computer-use consent is not implemented;
// fail closed rather than exposing an HWND through the macOS window-ID API.
int ShellNativeWindowID(CefWindowHandle) { return 0; }
bool ShellComputerUseEnabled() { return false; }
