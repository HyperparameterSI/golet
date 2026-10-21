// Author: HyperparameterSI

package golet

import (
	"io/fs"
	"net/url"
	"os"
	"strings"
	"sync"
)

// The app's frontend is its web UI. Windows load it with URLs without a
// scheme, such as "/" or "/settings" (see WindowOptions.URL), which resolve
// against:
//
//   - the dev server at devUrl (golet.config.ts) while `golet dev` runs the app;
//   - otherwise golet://localhost/, which serves the frontendDist files that
//     `golet build` embeds into the app (see SetFrontend).
//
// The golet scheme serves the frontend unless the app handles it itself with
// Protocol.Handle.
const frontendScheme = "golet"

var frontend struct {
	sync.Mutex
	fsys fs.FS
}

// SetFrontend serves fsys as the app's frontend at golet://localhost/. Paths
// without a file extension that match no file serve index.html, for
// client-side routing.
//
// `golet build` calls it with the frontendDist files it embeds into the app,
// so apps rarely need to. Call it to ship a frontend without the golet CLI,
// e.g. from an embed.FS.
func SetFrontend(fsys fs.FS) {
	frontend.Lock()
	frontend.fsys = fsys
	frontend.Unlock()
}

// devURL returns the dev server `golet dev` points the app at, or "".
func devURL() string {
	if !IsDev() {
		return ""
	}
	return os.Getenv("GOLET_DEV_URL")
}

// resolveURL resolves a URL without a scheme against the frontend.
func resolveURL(raw string) (string, error) {
	ref, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if ref.Scheme != "" {
		return raw, nil
	}
	base := frontendScheme + "://localhost/"
	if d := devURL(); d != "" {
		base = d
	}
	b, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	return b.ResolveReference(ref).String(), nil
}

// isDevOrigin reports whether origin ("scheme://host") is the dev server's.
func isDevOrigin(origin string) bool {
	d := devURL()
	if d == "" {
		return false
	}
	u, err := url.Parse(d)
	return err == nil && strings.EqualFold(u.Scheme+"://"+u.Host, origin)
}
