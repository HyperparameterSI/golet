// Author: HyperparameterSI

//go:build js

package golet

import "github.com/HyperparameterSI/golet/backend/internal/platform"

// Browsers show no web pages, for which `golet generate` writes the client:
// builds for them leave out the TypeScript generator, with go/parser.

// GenerateTypeScript renders the typed TypeScript client, which browsers
// do not generate: it returns platform.ErrUnsupported.
func GenerateTypeScript() ([]byte, error) { return nil, platform.ErrUnsupported }

// WriteTypeScript writes the TypeScript client to path. Bind binds nothing
// in browsers, so there is no client and path is not written.
func WriteTypeScript(path string) error { return nil }

// writeClient writes nothing: there is no client in browsers.
func writeClient(path string) error { return nil }
