// Author: HyperparameterSI
package assistant

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/HyperparameterSI/golet/backend/plugins/chromium"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Extension is a manually invoked user script with explicit exact-origin grants.
// It is separate from Chrome Manifest V3 extensions and has no chrome.* APIs.
type Extension struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Origins     []string `json:"origins"`
	Script      string   `json:"-"`
}

func LoadExtension(dir string) (Extension, error) {
	var ext Extension
	root, err := os.OpenRoot(dir)
	if err != nil {
		return ext, err
	}
	defer root.Close()
	manifest, err := root.ReadFile("golet-extension.json")
	if err != nil {
		return ext, err
	}
	var data struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Origins     []string `json:"origins"`
		Script      string   `json:"script"`
	}
	if len(manifest) > 16384 || json.Unmarshal(manifest, &data) != nil || data.Name == "" || len(data.Origins) == 0 {
		return ext, errors.New("invalid Golet extension manifest")
	}
	if !filepath.IsLocal(data.Script) {
		return ext, errors.New("script must be within extension directory")
	}
	for _, origin := range data.Origins {
		u, e := url.Parse(origin)
		if e != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || origin != u.Scheme+"://"+u.Host {
			return ext, errors.New("extensions require exact HTTP(S) origins")
		}
	}
	script, err := root.ReadFile(data.Script)
	if err != nil {
		return ext, err
	}
	if len(script) > 40000 {
		return ext, errors.New("extension exceeds 40 KB")
	}
	return Extension{data.Name, data.Description, data.Origins, string(script)}, nil
}
func (e Extension) Run(ctx context.Context, b *chromium.Browser, tab int) error {
	var current string
	if err := b.Evaluate(ctx, tab, "location.origin", &current); err != nil {
		return err
	}
	allowed := false
	for _, origin := range e.Origins {
		if origin == current {
			allowed = true
		}
	}
	if !allowed {
		return errors.New("extension is not permitted on this origin")
	}
	origin, _ := json.Marshal(current)
	// Check the origin again in the same evaluation that invokes the script.
	script := "(() => {if(location.origin!==" + string(origin) + ")throw Error('Origin changed');\n" + strings.TrimSpace(e.Script) + "\n})()"
	return b.Evaluate(ctx, tab, script, nil)
}
