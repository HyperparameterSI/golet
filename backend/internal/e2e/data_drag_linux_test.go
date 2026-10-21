// Author: HyperparameterSI

//go:build linux && (amd64 || arm64)

package e2e

import (
	"bytes"
	"path/filepath"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/internal/linux"
	"github.com/HyperparameterSI/golet/backend/internal/wait"
	"github.com/HyperparameterSI/golet/backend/transfer"
	"github.com/HyperparameterSI/golet/backend/ui"
)

func gtkPressSource(t *testing.T, w *golet.Window) {
	t.Helper()
	var supported bool
	golet.RunOnMain(func() {
		supported = linux.TestMoveSurfacePointer(w.NativeHandle(), 30, 30) && linux.TestPressButton(true)
	})
	if !supported {
		t.Skip("GTK drag tracking needs X11 and XTEST")
	}
	eventually(t, "the source pointer press", func() bool {
		ready := false
		golet.RunOnMain(func() { ready = linux.TestDataDragReady(w.NativeHandle()) })
		return ready
	})
}
func gtkDragResult(t *testing.T, result <-chan transfer.Result) transfer.Result {
	t.Helper()
	select {
	case r := <-result:
		return r
	case <-time.After(wait.Guard):
		t.Fatal("GTK drag did not complete")
	}
	return transfer.Result{}
}

// TestNativeGTKDataDrag uses GTK's tracker and XTEST input, including the
// asynchronous selection path used for sources outside the Golet registry.
func TestNativeGTKDataDrag(t *testing.T) {
	custom := transfer.Format("application/vnd.golet.e2e+json")
	var frames, providers atomic.Int32
	var drops []transfer.Drop
	var over atomic.Bool
	targetOperations := transfer.Copy | transfer.Move
	source := newWindow(t, golet.WindowOptions{Title: "GTK source", X: 30, Y: 40, Width: 240, Height: 200, Content: ui.View(func(c *ui.Context) { ui.Box(c).Fill() })})
	target := newWindow(t, golet.WindowOptions{Title: "GTK destination", X: 340, Y: 40, Width: 240, Height: 200, Content: ui.View(func(c *ui.Context) {
		frames.Add(1)
		e := ui.Box(c).Fill()
		opts := transfer.DropOptions{Formats: []transfer.Format{custom, transfer.Text, transfer.FileList, transfer.URIList}, Operations: targetOperations}
		_, hover := ui.DataDragOver(e, opts)
		over.Store(hover)
		if d, ok := ui.DropData(e, opts); ok {
			drops = append(drops, d)
		}
	})})
	t.Cleanup(func() {
		golet.RunOnMain(func() { linux.TestPressButton(false); linux.TestCancelDataDrag(source.NativeHandle()) })
	})
	eventually(t, "a frame", func() bool { return frames.Load() > 0 })
	url, _ := transfer.URLData("https://example.com/a#b")
	path := filepath.Join(t.TempDir(), "file #%.txt")
	files, _ := transfer.FileData(path)
	customData := transfer.New(transfer.NewItem(transfer.Lazy(custom, func() ([]byte, error) { providers.Add(1); return []byte("{\"id\":7}"), nil })))
	cases := []struct {
		name string
		data transfer.Data
		want transfer.Operation
	}{
		{"text", transfer.TextData("hello 日本語"), transfer.Copy},
		{"url", url, transfer.Copy},
		{"files", files, transfer.Copy},
		{"custom move", customData, transfer.Move},
	}
	for n, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.want == transfer.Move {
				before := frames.Load()
				target.Update(func() { targetOperations = transfer.Move })
				eventually(t, "move-only destination frame", func() bool { return frames.Load() > before })
			}
			gtkPressSource(t, source)
			result := make(chan transfer.Result, 2)
			var started bool
			golet.RunOnMain(func() {
				started = linux.TestStartSerializedDrag(source.NativeHandle(), tc.data, transfer.Copy|transfer.Move, func(r transfer.Result) { result <- r })
			})
			if !started {
				t.Fatal("native source did not start")
			}
			golet.RunOnMain(func() { linux.TestMoveSurfacePointer(target.NativeHandle(), 60, 60) })
			eventually(t, "the native destination hover", over.Load)
			if providers.Load() != 0 {
				t.Fatal("GTK fetched custom data during hover")
			}
			golet.RunOnMain(func() { linux.TestPressButton(false) })
			r := gtkDragResult(t, result)
			if r.Operation != tc.want || r.Err != nil {
				t.Fatalf("native result %+v", r)
			}
			eventually(t, "the serialized drop", func() bool { count := 0; golet.RunOnMain(func() { count = len(drops) }); return count == n+1 })
			var got transfer.Data
			golet.RunOnMain(func() { got = drops[n].Data })
			switch tc.name {
			case "files":
				paths, err := got.Files()
				if err != nil || !slices.Equal(paths, []string{path}) {
					t.Fatalf("files %v: %v", paths, err)
				}
			case "custom move":
				b, err := got.Read(custom)
				if err != nil || !bytes.Equal(b, []byte("{\"id\":7}")) {
					t.Fatalf("custom %q: %v", b, err)
				}
			default:
				for _, f := range tc.data.Formats() {
					want, _ := tc.data.Read(f)
					b, err := got.Read(f)
					if err != nil || !bytes.Equal(b, want) {
						t.Errorf("%s %q: %v", f, b, err)
					}
				}
			}
			if tc.name == "custom move" && providers.Load() != 1 {
				t.Fatalf("provider called %d times", providers.Load())
			}
			select {
			case <-result:
				t.Fatal("source completed twice")
			default:
			}
		})
	}
}

func TestNativeGTKDragCancellationCleanup(t *testing.T) {
	for _, destroy := range []bool{false, true} {
		t.Run(map[bool]string{false: "cancel", true: "source destruction"}[destroy], func(t *testing.T) {
			source := newWindow(t, golet.WindowOptions{Title: "GTK cancel", X: 30, Y: 40, Width: 240, Height: 200, Content: ui.View(func(c *ui.Context) { ui.Box(c).Fill() })})
			gtkPressSource(t, source)
			result := make(chan transfer.Result, 2)
			golet.RunOnMain(func() {
				linux.TestStartSerializedDrag(source.NativeHandle(), transfer.TextData("cancel"), transfer.Copy, func(r transfer.Result) { result <- r })
			})
			if destroy {
				source.Destroy()
			} else {
				golet.RunOnMain(func() { linux.TestCancelDataDrag(source.NativeHandle()) })
			}
			golet.RunOnMain(func() { linux.TestPressButton(false) })
			r := gtkDragResult(t, result)
			if !r.Canceled || !r.Escaped || r.Operation != transfer.None {
				t.Fatalf("result %+v", r)
			}
			eventually(t, "released GTK drag resources", func() bool {
				count := 1
				golet.RunOnMain(func() { count = linux.TestDataDragResources() })
				return count == 0
			})
			select {
			case <-result:
				t.Fatal("source completed twice")
			default:
			}
		})
	}
}
