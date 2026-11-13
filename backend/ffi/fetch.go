// Author: HyperparameterSI

//go:build !golet_nodownload

package ffi

import (
	"context"
	"time"

	"github.com/HyperparameterSI/golet/backend/internal/download"
	"github.com/HyperparameterSI/golet/backend/internal/nativelib"
)

// fetch downloads a published library into the user's cache, for programs
// that are not packaged apps. Production builds of the CLI leave it out
// (golet_nodownload), and HTTP with it, as packaged apps never download.
var fetch = func(f nativelib.File, path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	return download.Cached(ctx, download.File{URL: f.URL, SHA256: f.SHA256}, path)
}
