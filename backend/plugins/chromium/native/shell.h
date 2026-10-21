// Author: HyperparameterSI
#pragma once
#include <string>
#include "include/cef_browser.h"

// The Cocoa shell and CEF callbacks run on the same UI thread.
void CreateBrowserShell();
void ConfigureShellChild(CefWindowInfo& info);
void ShellAttachBrowser(int id, CefWindowHandle view);
void ShellRemoveBrowser(int id);
void ShellDetachBrowser(int id);
void ShellSelectBrowser(int id);
void ShellUpdateBrowser(int id, const std::string& title, const std::string& url,
                        bool loading, bool back, bool forward);
void ShellDestroy();
std::string ShellHomeURL();
void BrowserUIAction(const std::string& method, int id = 0, const std::string& value = "");
