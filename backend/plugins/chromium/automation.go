// Author: HyperparameterSI
package chromium

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DevTools executes a privileged CDP command on an explicit tab over the private
// host pipe. Embedders must authorize callers; never expose this to web pages.
func (b *Browser) DevTools(ctx context.Context, tab int, method string, params any) (json.RawMessage, error) {
	if tab <= 0 || method == "" {
		return nil, errors.New("chromium: tab and method required")
	}
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	if len(raw) > 60000 {
		return nil, errors.New("chromium: command exceeds 60 KB")
	}
	return b.request(ctx, message{Method: "devtools", BrowserID: tab, DevToolsMethod: method, Params: raw})
}

// Evaluate runs trusted JavaScript and decodes its by-value result. This is a
// privileged API, not an AI tool: assistant actions use fixed scripts instead.
func (b *Browser) Evaluate(ctx context.Context, tab int, expression string, out any) error {
	raw, err := b.DevTools(ctx, tab, "Runtime.evaluate", map[string]any{"expression": expression, "returnByValue": true, "awaitPromise": true, "timeout": 5000})
	if err != nil {
		return err
	}
	var result struct {
		Result struct {
			Value json.RawMessage `json:"value"`
		}
		Exception json.RawMessage `json:"exceptionDetails"`
	}
	if err = json.Unmarshal(raw, &result); err != nil {
		return err
	}
	if len(result.Exception) > 0 {
		return errors.New("chromium: page script failed")
	}
	if out == nil {
		return nil
	}
	if len(result.Result.Value) == 0 {
		return errors.New("chromium: script returned no serializable value")
	}
	return json.Unmarshal(result.Result.Value, out)
}

func extensionArguments(opts Options) ([]string, error) {
	style, err := opts.browserStyle()
	if err != nil {
		return nil, err
	}
	var args []string
	if opts.AppName != "" {
		args = append(args, "--golet-app-name="+opts.AppName)
	}
	if opts.AppMode {
		if style != StyleAlloy || len(opts.Extensions) != 0 || opts.AssistantURL != "" {
			return nil, errors.New("chromium: AppMode requires Alloy without extensions or AssistantURL")
		}
		args = append(args, "--golet-app-mode")
	}
	if style == StyleChrome {
		args = append(args, "--golet-chrome-style")
	}
	if len(opts.Extensions) == 0 {
		return args, nil
	}
	if style != StyleChrome {
		return nil, errors.New("chromium: Chrome extensions require StyleChrome")
	}
	var paths []string
	for _, dir := range opts.Extensions {
		path, err := filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
		if strings.Contains(path, ",") {
			return nil, errors.New("chromium: extension path cannot contain commas")
		}
		data, err := os.ReadFile(filepath.Join(path, "manifest.json"))
		if err != nil {
			return nil, err
		}
		var manifest struct {
			Version int    `json:"manifest_version"`
			Name    string `json:"name"`
		}
		if json.Unmarshal(data, &manifest) != nil || manifest.Version != 3 || manifest.Name == "" {
			return nil, fmt.Errorf("chromium: %s requires a Manifest V3 manifest", path)
		}
		paths = append(paths, path)
	}
	return append(args, "--load-extension="+strings.Join(paths, ",")), nil
}
