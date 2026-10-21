// Author: HyperparameterSI

package golet

import "os"

// production is set to "1" by `golet build` through
// -ldflags "-X github.com/HyperparameterSI/golet/backend.production=1".
var production string

// IsDev reports whether the app runs in development: true unless it was
// built for production with `golet build`, or GOLET_ENV=production is set.
// Development builds enable the web inspector by default.
func IsDev() bool {
	if production == "1" {
		return false
	}
	return os.Getenv("GOLET_ENV") != "production"
}

// devLaunched tells that `golet dev` launched the app, which it does with
// GOLET_DEV=1: App.Relaunch then has golet dev start it again.
var devLaunched = devMarker()

func devMarker() bool {
	set := os.Getenv("GOLET_DEV") == "1"
	// Not meant for child processes.
	os.Unsetenv("GOLET_DEV")
	return set && production != "1"
}

// launchedByDev reports whether `golet dev` launched the app.
func launchedByDev() bool { return devLaunched }
