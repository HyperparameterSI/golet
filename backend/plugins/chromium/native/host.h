// Author: HyperparameterSI
#pragma once
#include "include/cef_app.h"
#include "include/cef_client.h"

CefRefPtr<CefApp> CreateGoletChromiumApp();
void CloseGoletChromium();
void ConfigureGoletChromium(CefSettings& settings, CefRefPtr<CefCommandLine> args);
