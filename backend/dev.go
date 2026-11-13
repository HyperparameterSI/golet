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

// devRelaunchFile is the file in which `golet dev` looks for the arguments
// of the build that App.Relaunch has it start again, GOLET_DEV_RELAUNCH,
// which the core takes out of its environment as GOLET_DEV.
var devRelaunchFile = func() string {
	f := os.Getenv("GOLET_DEV_RELAUNCH")
	os.Unsetenv("GOLET_DEV_RELAUNCH")
	return f
}()

// devBackground tells that `golet dev -background` launched the app, with
// GOLET_DEV_BACKGROUND=1: it starts without coming to the front, and its
// windows show without the focus until it takes it (backgroundShow), so
// that the rebuild of a save leaves the focus in the editor. Production
// builds ignore it.
var devBackground = func() bool {
	set := os.Getenv("GOLET_DEV_BACKGROUND") == "1"
	os.Unsetenv("GOLET_DEV_BACKGROUND")
	return set && production != "1"
}()

// launchedByDev reports whether `golet dev` launched the app.
func launchedByDev() bool { return devLaunched }
