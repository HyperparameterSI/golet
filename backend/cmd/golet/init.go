// Author: HyperparameterSI

package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"
)

//go:embed all:template
var templateFS embed.FS

type templateData struct {
	Name       string
	Slug       string
	Module     string
	Identifier string
	// Runtime is the version of the golet-runtime npm package.
	Runtime string
	// CLI is the version of the golet-cli npm package, which runs golet in
	// the scripts of package.json, and Golet the command they run.
	CLI, Golet string
	// ConfigImport is the module golet.config.ts imports defineConfig from.
	ConfigImport string
}

// The templates of golet init: a TypeScript frontend, and native UI.
const (
	webTemplate    = "web"
	shadcnTemplate = "shadcn"
	nativeTemplate = "native"
)

// cliTool is the package of the CLI, which projects of native UI depend on
// as a tool of their module.
const cliTool = "github.com/HyperparameterSI/golet/backend/cmd/golet"

func runInit(args []string) error {
	flags := newFlags("init", "[flags] <dir>", "Creates a new Golet project. The default shadcn template uses React, Tailwind and\neditable UI components. The web template has a TypeScript frontend built\nwith Vite: Bun installs its dependencies, the golet-cli package among them,\nand runs its scripts, bun run dev and bun run build. The native template\nshows a user interface Golet draws itself, written in Go with package ui:\nthe module has the CLI as a tool, for go tool golet dev and go tool golet\nbuild.")
	name := flags.String("name", "", "application name (default: directory name)")
	module := flags.String("module", "", "Go module path (default: directory name)")
	local := flags.String("golet", "", "path to a local checkout of Golet to use via a replace directive")
	tmpl := flags.String("template", shadcnTemplate, "the project: shadcn (React + Tailwind, default), web (vanilla TypeScript), or native (Go UI)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		flags.Usage()
		return flag2Err("missing project directory")
	}
	if *tmpl != webTemplate && *tmpl != nativeTemplate && *tmpl != shadcnTemplate {
		return flag2Err(fmt.Sprintf("unknown template %q: shadcn, web or native", *tmpl))
	}
	dir, data, err := newProject(flags.Arg(0), *name, *module)
	if err != nil {
		return err
	}
	if *tmpl == nativeTemplate {
		return initNative(dir, data, *local)
	}
	data.Runtime, data.CLI, data.Golet, data.ConfigImport = "^"+version, "^"+version, "golet", "golet-cli"
	if *local != "" {
		// The packages of the local checkout, like the Go module, and its
		// CLI, which go run builds from the replaced module.
		checkout, err := filepath.Abs(*local)
		if err != nil {
			return err
		}
		runtime := filepath.Join(checkout, "frontend", "packages", "runtime")
		if _, err := os.Stat(filepath.Join(runtime, "dist", "index.js")); err != nil {
			return fmt.Errorf("the packages of %s are not built: run pnpm install && pnpm --dir frontend/packages/runtime build there first", checkout)
		}
		data.Runtime = "file:" + runtime
		data.CLI, data.Golet = "", "go run github.com/HyperparameterSI/golet/backend/cmd/golet"
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		data.ConfigImport = moduleSpecifier(dir, filepath.Join(checkout, "frontend", "packages", "cli", "index.js"))
	}

	logf("Creating %s", dir)
	if err := writeProject(dir, *tmpl, data, *local, false); err != nil {
		return err
	}

	if _, err := exec.LookPath("bun"); err != nil {
		warnf("Dependencies not installed automatically: run pnpm install or npm install (the legacy web template also needs Bun)")
	} else {
		t := con.start("Installing the dependencies")
		if err := t.end(run(t, "bun", dir, "bun install"), "Installed the dependencies"); err != nil {
			return err
		}
	}
	if c, err := loadConfig(dir); err != nil {
		warnf("%v; run golet generate once it can be read", err)
	} else {
		bin := tempBinary(c.executableName())
		defer os.Remove(bin)
		if err := buildBinary(c, bin, nil); err == nil {
			_, _ = generateBindings(c, bin)
		}
	}

	fmt.Printf("\nCreated %s. Next steps:\n\n  cd %s\n  npm run dev      # develop with live reload\n  npm run build    # package the app\n\n", data.Name, relPath(dir))
	return nil
}

// newProject returns the absolute directory of a new project, which must
// be empty or not exist, and what its template is filled with.
func newProject(path, name, module string) (string, templateData, error) {
	dir, err := filepath.Abs(path)
	if err != nil {
		return "", templateData{}, err
	}
	if entries, err := os.ReadDir(dir); err == nil && len(entries) > 0 {
		return "", templateData{}, fmt.Errorf("%s is not empty", dir)
	}
	data := templateData{Name: name, Module: module}
	if data.Name == "" {
		data.Name = filepath.Base(dir)
	}
	data.Slug = slugify(data.Name)
	if data.Module == "" {
		data.Module = data.Slug
	}
	data.Identifier = "com.example." + strings.ReplaceAll(data.Slug, "-", "")
	return dir, data, nil
}

// initNative creates a project of native UI in dir: a Go module, without
// a frontend.
func initNative(dir string, data templateData, local string) error {
	logf("Creating %s", dir)
	if err := writeProject(dir, nativeTemplate, data, local, true); err != nil {
		return err
	}
	fmt.Printf("\nCreated %s. Next steps:\n\n  cd %s\n  go tool golet dev      # develop with live reload\n  go test               # test the view without a window\n  go tool golet build    # package the app\n\n", data.Name, relPath(dir))
	return nil
}

// writeProject writes a template into dir with the app's icon and its Go
// module, which depends on Golet, from local when it is a checkout, and,
// with tool, has the CLI as a tool.
func writeProject(dir, tmpl string, data templateData, local string, tool bool) error {
	if err := writeTemplate(dir, tmpl, data); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(dir, resourcesDir), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, resourcesDir, "icon.png"), defaultIcon(), 0o644); err != nil {
		return err
	}
	gomod := fmt.Sprintf("module %s\n\ngo %s\n", data.Module, goVersion())
	if local != "" {
		abs, err := filepath.Abs(local)
		if err != nil {
			return err
		}
		if tool {
			gomod += "\ntool " + cliTool + "\n"
		}
		gomod += fmt.Sprintf("\nrequire github.com/HyperparameterSI/golet v0.0.0\n\nreplace github.com/HyperparameterSI/golet => %s\n", abs)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0o644); err != nil {
		return err
	}
	if local == "" {
		get := []string{"get", "github.com/HyperparameterSI/golet@latest"}
		if tool {
			get = []string{"get", "-tool", cliTool + "@latest"}
		}
		if err := goCommand(dir, nil, get...).Run(); err != nil {
			warnf("Could not fetch github.com/HyperparameterSI/golet: %v", err)
		}
	}
	if err := goCommand(dir, nil, "mod", "tidy").Run(); err != nil {
		warnf("go mod tidy failed: %v", err)
	}
	return nil
}

// relPath returns dir relative to the working directory when it is below
// it, for showing.
func relPath(dir string) string {
	if wd, err := os.Getwd(); err == nil {
		if r, err := filepath.Rel(wd, dir); err == nil && !strings.HasPrefix(r, "..") {
			return r
		}
	}
	return dir
}

// moduleSpecifier returns how a module in dir imports the file target: a
// relative path between their real locations, which runtimes resolve
// imports from, or a file URL where there is none, as between the drives
// of Windows.
func moduleSpecifier(dir, target string) string {
	if real, err := filepath.EvalSymlinks(dir); err == nil {
		dir = real
	}
	if real, err := filepath.EvalSymlinks(target); err == nil {
		target = real
	}
	rel, err := filepath.Rel(dir, target)
	if err != nil {
		path := filepath.ToSlash(target)
		if !strings.HasPrefix(path, "/") {
			path = "/" + path // C:/…
		}
		return (&url.URL{Scheme: "file", Path: path}).String()
	}
	rel = filepath.ToSlash(rel)
	if !strings.HasPrefix(rel, "../") {
		rel = "./" + rel
	}
	return rel
}

// templateFuncs are the functions of the templates: json renders a value
// as JSON, which is also a JavaScript literal.
var templateFuncs = template.FuncMap{
	"json": func(v any) (string, error) {
		b, err := json.Marshal(v)
		return string(b), err
	},
}

// writeTemplate installs the bundled skills and template/<tmpl> into dir.
func writeTemplate(dir, tmpl string, data templateData) error {
	if err := installSkills(dir); err != nil {
		return err
	}
	return writeTemplateDir(dir, tmpl, data)
}

func writeTemplateDir(dir, tmpl string, data templateData) error {
	root := "template/" + tmpl
	return fs.WalkDir(templateFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		target := filepath.Join(dir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		content, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.HasSuffix(target, ".tmpl") {
			target = strings.TrimSuffix(target, ".tmpl")
			t, err := template.New(rel).Funcs(templateFuncs).Parse(string(content))
			if err != nil {
				return err
			}
			var b bytes.Buffer
			if err := t.Execute(&b, data); err != nil {
				return err
			}
			content = b.Bytes()
		}
		if filepath.Base(target) == "gitignore" {
			target = filepath.Join(filepath.Dir(target), ".gitignore")
		}
		return os.WriteFile(target, content, 0o644)
	})
}

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.Trim(slugRe.ReplaceAllString(strings.ToLower(s), "-"), "-")
	if s == "" {
		return "app"
	}
	return s
}

func goVersion() string {
	v := strings.TrimPrefix(runtime.Version(), "go")
	if i := strings.IndexAny(v, " -"); i >= 0 {
		v = v[:i]
	}
	return v
}

type flagErr string

func (e flagErr) Error() string { return string(e) }

func flag2Err(s string) error { return flagErr(s) }

// defaultIcon returns the shared Golet app icon used by newly scaffolded projects.
// The PNG matches the dashboard brand mark. Regenerate it from the repository root:
// rsvg-convert backend/cmd/golet/assets/icon.svg -o backend/cmd/golet/assets/icon.png
func defaultIcon() []byte {
	return defaultIconPNG
}

//go:embed assets/icon.png
var defaultIconPNG []byte
