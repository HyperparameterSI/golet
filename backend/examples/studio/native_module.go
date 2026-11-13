// Author: HyperparameterSI
package main

import (
	"fmt"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/examples/native-module/modules/places"
	"github.com/HyperparameterSI/golet/backend/ui"
)

// openNativeModule hosts a MapKit view alongside Go-rendered controls.
// The optional Swift module must be built by the Golet CLI on macOS.
func openNativeModule() (*golet.Window, error) {
	m, err := places.NewMap()
	if err != nil {
		return nil, fmt.Errorf("原生地图需要 macOS 和已构建的 Swift 模块；请用 golet dev 启动 Studio: %w", err)
	}
	m.Show(places.Place{Latitude: 48.8566, Longitude: 2.3522}, 0.2)
	win := golet.NewWindow(golet.WindowOptions{
		Title: "Golet · Native Module", Width: 960, Height: 640,
		TitleBarStyle:        golet.TitleBarHiddenInset,
		TrafficLightPosition: &golet.Point{X: 12, Y: 11},
		Content: ui.View(func(c *ui.Context) {
			ui.Column(c).Fill().Padding(24).Gap(16).Children(func() {
				ui.Text(c, "NativeView · MapKit").FontSize(22)
				ui.Row(c).Gap(12).Children(func() {
					if ui.Button(c, "巴黎").Clicked() {
						m.Show(places.Place{Latitude: 48.8566, Longitude: 2.3522}, 0.2)
					}
					if ui.Button(c, "东京").Clicked() {
						m.Show(places.Place{Latitude: 35.6762, Longitude: 139.6503}, 0.3)
					}
				})
				ui.NativeView(c, m).Grow(1).Radius(12)
			})
		}),
	})
	return win, nil
}
