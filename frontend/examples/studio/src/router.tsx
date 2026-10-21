// Author: HyperparameterSI

import { createRouter } from "@tanstack/react-router";
import { LoadingScreen } from "./components/loading-screen";
import { routeTree } from "./routeTree.gen";
export function getRouter() {
  return createRouter({
    routeTree,
    defaultPreload: "intent",
    scrollRestoration: true,
    defaultPendingComponent: LoadingScreen,
  });
}
declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof getRouter>;
  }
}
