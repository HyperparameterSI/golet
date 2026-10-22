// Author: HyperparameterSI
#include <windows.h>
#include "host.h"
#include "include/cef_command_line.h"
#include "include/cef_sandbox_win.h"

// The official CEF bootstrap supplies sandbox state for every process. Keeping
// the console entry point preserves the private stdin/stdout JSON transport.
CEF_BOOTSTRAP_EXPORT int RunConsoleMain(int argc, char* argv[], void* sandbox_info,
                                        cef_version_info_t*) {
  CefMainArgs main_args(GetModuleHandleW(nullptr));
  const int child = CefExecuteProcess(main_args, nullptr, sandbox_info);
  if (child >= 0) return child;
  auto args = CefCommandLine::CreateCommandLine();
  args->InitFromString(GetCommandLineW());
  CefSettings settings;
  ConfigureGoletChromium(settings, args);
  if (!CefInitialize(main_args, settings, CreateGoletChromiumApp(), sandbox_info))
    return CefGetExitCode();
  CefRunMessageLoop();
  CefShutdown();
  return 0;
}
