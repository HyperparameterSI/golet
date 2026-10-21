// Author: HyperparameterSI

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func copyFile(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	info, err := input.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("not a regular file: %s", source)
	}
	if err = os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return err
	}
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func copyTree(ctx context.Context, source, destination string) error {
	return filepath.WalkDir(source, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err = ctx.Err(); err != nil {
			return err
		}
		rel, err := filepath.Rel(source, file)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		if entry.Type()&os.ModeSymlink != 0 {
			link, err := os.Readlink(file)
			if err != nil {
				return err
			}
			resolved, err := filepath.Rel(source, filepath.Join(filepath.Dir(file), link))
			if err != nil {
				return err
			}
			if filepath.IsAbs(link) || !filepath.IsLocal(resolved) {
				return fmt.Errorf("framework symlink escapes its root: %s", file)
			}
			return os.Symlink(link, target)
		}
		return copyFile(file, target)
	})
}

func xmlText(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

func bundle(directory, name, binary, helperSuffix string, helper bool, opts options) error {
	contents := filepath.Join(directory, "Contents")
	if err := copyFile(binary, filepath.Join(contents, "MacOS", name)); err != nil {
		return err
	}
	id := opts.identifier
	if id == "" {
		id = "dev.golet.chromium"
	}
	display := opts.appName
	if display == "" {
		display = "Golet Browser"
	}
	extra := "<key>CFBundleIconFile</key><string>AppIcon.icns</string><key>NSPrincipalClass</key><string>GoletChromiumApplication</string>"
	if helper {
		id += ".helper" + helperSuffix
		display += " Helper" + helperSuffix
		extra = "<key>LSUIElement</key><true/>"
	}
	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>CFBundleName</key><string>%s</string>
<key>CFBundleDisplayName</key><string>%s</string>
<key>CFBundleExecutable</key><string>%s</string>
<key>CFBundleIdentifier</key><string>%s</string>
<key>CFBundlePackageType</key><string>APPL</string>
<key>CFBundleVersion</key><string>0.1.0</string>
<key>CFBundleShortVersionString</key><string>0.1.0</string>
<key>NSHighResolutionCapable</key><true/>
<key>LSMinimumSystemVersion</key><string>12.0</string>
%s
</dict></plist>
`, xmlText(display), xmlText(display), xmlText(name), xmlText(id), extra)
	return os.WriteFile(filepath.Join(contents, "Info.plist"), []byte(plist), 0644)
}

func packageApp(ctx context.Context, root, sdk, work, destination, temp string, opts options) error {
	app := filepath.Join(destination, "Golet Browser.app")
	if err := bundle(app, "Golet Browser", filepath.Join(work, "golet-chromium"), "", false, opts); err != nil {
		return err
	}
	frameworks := filepath.Join(app, "Contents/Frameworks")
	if err := copyTree(ctx, filepath.Join(sdk, "Release/Chromium Embedded Framework.framework"), filepath.Join(frameworks, "Chromium Embedded Framework.framework")); err != nil {
		return err
	}
	for _, role := range []string{"", "Alerts", "GPU", "Plugin", "Renderer"} {
		name := "Golet Browser Helper"
		suffix := ""
		if role != "" {
			name += " (" + role + ")"
			suffix = "." + strings.ToLower(role)
		}
		if err := bundle(filepath.Join(frameworks, name+".app"), name, filepath.Join(work, "golet-chromium-helper"), suffix, true, opts); err != nil {
			return err
		}
	}
	resources := filepath.Join(app, "Contents/Resources")
	if err := os.MkdirAll(resources, 0755); err != nil {
		return err
	}
	newTab := opts.newTabDir
	if newTab == "" {
		newTab = filepath.Join(root, "frontend/examples/browser/newtab")
	}
	if _, err := os.Stat(filepath.Join(newTab, "index.html")); err != nil {
		return err
	}
	if err := copyTree(ctx, newTab, filepath.Join(resources, "newtab")); err != nil {
		return err
	}
	icon := opts.icon
	if icon == "" {
		icon = filepath.Join(root, "backend/cmd/golet/assets/icon.png")
	}
	if opts.newTabDir == "" {
		logo := ""
		if opts.icon != "" {
			logo = "product-icon.png"
			if err := copyFile(icon, filepath.Join(resources, "newtab", logo)); err != nil {
				return err
			}
		}
		branding, err := json.Marshal(map[string]string{"name": opts.appName, "icon": logo})
		if err != nil {
			return err
		}
		// Replace the copied (possibly read-only module-cache) default atomically.
		file := filepath.Join(resources, "newtab/branding.js")
		if err := os.Remove(file); err != nil {
			return err
		}
		if err := os.WriteFile(file, append(append([]byte("window.browserBranding = "), branding...), ';'), 0644); err != nil {
			return err
		}
	}
	iconset := filepath.Join(temp, "AppIcon.iconset")
	if err := os.Mkdir(iconset, 0755); err != nil {
		return err
	}
	for _, size := range []int{16, 32, 128, 256, 512} {
		for _, scale := range []int{1, 2} {
			name := fmt.Sprintf("icon_%dx%d", size, size)
			if scale == 2 {
				name += "@2x"
			}
			name += ".png"
			pixels := strconv.Itoa(size * scale)
			if err := run(ctx, "sips", "-z", pixels, pixels, icon, "--out", filepath.Join(iconset, name)); err != nil {
				return err
			}
		}
	}
	if err := run(ctx, "iconutil", "-c", "icns", iconset, "-o", filepath.Join(resources, "AppIcon.icns")); err != nil {
		return err
	}
	for _, name := range []string{"LICENSE.txt", "CREDITS.html"} {
		if err := copyFile(filepath.Join(sdk, name), filepath.Join(resources, name)); err != nil {
			return err
		}
	}
	return run(ctx, "codesign", "--force", "--deep", "--sign", "-", app)
}
