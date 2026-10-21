// Author: HyperparameterSI
package main

import (
	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/examples/internal/editordemo"
	"log"
)

func main() {
	golet.App.WhenReady(func() {
		golet.App.SetMenu(golet.NewMenu([]*golet.MenuItem{{Role: golet.RoleAppMenu}, {Role: golet.RoleEditMenu}, {Role: golet.RoleWindowMenu}}))
		editordemo.Open()
	})
	if err := golet.App.Run(); err != nil {
		log.Fatal(err)
	}
}
