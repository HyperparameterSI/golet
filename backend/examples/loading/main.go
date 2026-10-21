// Author: HyperparameterSI
package main

import (
	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/examples/internal/loadingdemo"
	"log"
)

func main() {
	golet.App.WhenReady(func() { loadingdemo.Open() })
	if err := golet.App.Run(); err != nil {
		log.Fatal(err)
	}
}
