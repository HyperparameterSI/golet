// Author: HyperparameterSI

package main

import (
	"github.com/HyperparameterSI/golet/backend"
	"log"
)

func main() {
	monitor := newProcessMonitor()
	defer monitor.close()
	golet.Bind(monitor)
	golet.App.WhenReady(func() {
		golet.NewWindow(golet.WindowOptions{Title: "Golet Observatory", Width: 1440, Height: 960, MinWidth: 760, MinHeight: 600, BackgroundColor: "#f6f7f9", StateKey: "dashboard", URL: "/"})
	})
	if err := golet.App.Run(); err != nil {
		log.Fatal(err)
	}
}
