// Author: HyperparameterSI

// Todo is a small but complete Golet app: typed services and events shared
// with a TypeScript frontend, persistence, dialogs, menus and multiple
// windows kept in sync. In examples/todo:
//
//	bun run dev    # live reload with the Vite dev server
//	bun run build  # the app and a disk image in build/
package main

import (
	"log"

	"github.com/HyperparameterSI/golet/backend"
)

func main() {
	app := golet.App
	// Set the name first: it determines the user data directory.
	app.SetName("Todo")

	store, err := OpenStore()
	if err != nil {
		log.Fatal(err)
	}
	golet.Bind(store)

	app.WhenReady(func() {
		app.SetMenu(golet.NewMenu([]*golet.MenuItem{
			{Role: golet.RoleAppMenu},
			{Label: "File", Submenu: []*golet.MenuItem{
				{Label: "New Window", Accelerator: "CmdOrCtrl+N", Click: func(*golet.MenuItem, *golet.Window) { openWindow() }},
				{Label: "Export…", Accelerator: "CmdOrCtrl+E", Click: func(_ *golet.MenuItem, w *golet.Window) {
					go store.export(w)
				}},
				golet.Separator(),
				{Role: golet.RoleClose},
			}},
			{Role: golet.RoleEditMenu},
			{Role: golet.RoleViewMenu},
			{Role: golet.RoleWindowMenu},
		}))
		openWindow()
	})
	// Keep the app alive in the Dock and reopen a window when clicked.
	app.OnWindowAllClosed(func() {})
	app.OnActivate(func(hasVisibleWindows bool) {
		if !hasVisibleWindows {
			openWindow()
		}
	})
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func openWindow() {
	golet.NewWindow(golet.WindowOptions{
		Title:           "Todo",
		Width:           460,
		Height:          640,
		MinWidth:        360,
		MinHeight:       420,
		TitleBarStyle:   golet.TitleBarHidden,
		TitleBarHeight:  52,                             // the header's
		BackgroundColor: "light-dark(#f5f5f7, #1c1c1e)", // the page's --bg
		URL:             "/",                            // the frontend: devUrl in development, frontendDist once built
	})
}
