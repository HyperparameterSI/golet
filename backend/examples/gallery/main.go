// Author: HyperparameterSI
package main

import (
	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/examples/internal/gallerydemo"
	"log"
)

func main() {
	gallerydemo.ConfigureLanguages()
	golet.App.WhenReady(func() { gallerydemo.Open("Overview") })
	if err := golet.App.Run(); err != nil {
		log.Fatal(err)
	}
}
