// Author: HyperparameterSI
#pragma once
#include <string>
#include "include/cef_browser.h"

// Platform shells and CEF callbacks run on the same UI thread.
void CreateBrowserShell();
void ConfigureShellChild(CefWindowInfo& info);
void ShellAttachBrowser(int id, CefWindowHandle view);
void ShellRemoveBrowser(int id);
void ShellDetachBrowser(int id);
void ShellSelectBrowser(int id);
void ShellUpdateBrowser(int id, const std::string& title, const std::string& url,
                        bool loading, bool back, bool forward);
void ShellDestroy();
// Author: HyperparameterSI. Applies a validated theme on the host UI thread.
bool ShellSetTheme(const std::string& json);
std::string ShellHomeURL();
void BrowserUIAction(const std::string& method, int id = 0, const std::string& value = "");

// Author: HyperparameterSI. UI-thread-only native consent and window identity.
int ShellNativeWindowID(CefWindowHandle view);
bool ShellComputerUseEnabled();
