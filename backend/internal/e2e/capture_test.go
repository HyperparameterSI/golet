// Author: HyperparameterSI

//go:build !js

package e2e

import (
	"os"
	"runtime"
	"slices"
	"testing"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/internal/wait"
)

// capturePage captures with getUserMedia() at app://localhost/capture, a
// secure context the app trusts, which gets the devices it asks for.
// capture answers "ok" or the error, and track tells the state of the
// stream's track of a kind.
const capturePage = `<!doctype html><p id=state>idle</p><script>
window.capture = async (constraints) => {
	try {
		window.stream = await navigator.mediaDevices.getUserMedia(constraints);
		return "ok";
	} catch (e) {
		return "error: " + e.name + ": " + e.message;
	}
};
window.track = (kind) => {
	const t = window.stream && window.stream.getTracks().find((t) => t.kind === kind);
	return t ? t.readyState + (t.muted ? " muted" : "") : "none";
};
</script>`

// A page that captures nothing captures no kind, and setting a kind's
// state changes nothing.
func TestCaptureStateWithoutCapture(t *testing.T) {
	w := newWindow(t, golet.WindowOptions{Title: "No capture", Width: 300, Height: 200})
	p := w.Page()
	changes := make(chan golet.CaptureKind, 16) // listeners must not block
	defer p.OnCaptureStateChanged(func(kind golet.CaptureKind, _ golet.CaptureState) { changes <- kind })()
	p.LoadHTML(`<p>nothing</p>`, "")
	waitFor(t, w, `document.body.textContent === "nothing" && document.readyState === "complete"`)
	kinds := []golet.CaptureKind{golet.CaptureCamera, golet.CaptureMicrophone, golet.CaptureDisplay}
	for _, kind := range kinds {
		p.SetCaptureState(kind, golet.CaptureMuted)
	}
	mustEval(t, w, `1`) // a round trip, for the changes to land
	for _, kind := range kinds {
		if s := p.CaptureState(kind); s != golet.CaptureNone {
			t.Errorf("CaptureState(%d) = %d for a page capturing nothing", kind, s)
		}
	}
	select {
	case kind := <-changes:
		t.Errorf("a page capturing nothing changed its capture of %d", kind)
	default:
	}
}

// The page captures the camera and the microphone, which the app mutes,
// unmutes and stops, and the page stops the microphone itself. It runs
// with WebKit's mock devices: on Linux, and on macOS with
// GOLET_E2E_CAPTURE=1, where WebKit may yet ask the system; WebView2 tells
// no capture.
func TestCaptureState(t *testing.T) {
	if runtime.GOOS == "darwin" && os.Getenv("GOLET_E2E_CAPTURE") != "1" {
		t.Skip("captures WebKit's mock devices: set GOLET_E2E_CAPTURE=1")
	}
	w := newWindow(t, golet.WindowOptions{Title: "Capture", Width: 300, Height: 200})
	if !mockCapture(w) {
		t.Skip("no mock capture devices")
	}
	p := w.Page()
	type change struct {
		kind  golet.CaptureKind
		state golet.CaptureState
	}
	changes := make(chan change, 64) // listeners must not block
	defer p.OnCaptureStateChanged(func(kind golet.CaptureKind, s golet.CaptureState) { changes <- change{kind, s} })()
	// told waits until the last change listeners heard of kind is to s,
	// and CaptureState says so too: a change on the way, as SetAudioMuted
	// on macOS mutes the camera with the microphone and unmutes it again,
	// may be heard after the one before it is done. The changes of other
	// kinds, which may come first, are kept for later.
	var heard []change
	told := func(what string, kind golet.CaptureKind, s golet.CaptureState) {
		t.Helper()
		for {
			last := -1
			for i, c := range heard {
				if c.kind == kind {
					last = i
				}
			}
			if last >= 0 && heard[last].state == s && p.CaptureState(kind) == s {
				heard = slices.DeleteFunc(heard, func(c change) bool { return c.kind == kind })
				break
			}
			heard = append(heard, wait.For(t, what, changes))
		}
	}
	if err := p.LoadURL("app://localhost/capture"); err != nil {
		t.Fatal(err)
	}
	waitFor(t, w, `location.pathname === "/capture" && window.capture`)
	ctx := wait.Context(t)
	if got, err := p.EvalContext(ctx, `capture({audio: true, video: true})`); err != nil || got != "ok" {
		info := mustEval(t, w, `[isSecureContext, typeof navigator.mediaDevices].join(" ")`)
		t.Skipf("the page can't capture: %v %v (secure context, navigator.mediaDevices: %v)", got, err, info)
	}
	told("the camera capturing", golet.CaptureCamera, golet.CaptureActive)
	told("the microphone capturing", golet.CaptureMicrophone, golet.CaptureActive)
	if s := p.CaptureState(golet.CaptureDisplay); s != golet.CaptureNone {
		t.Errorf("CaptureState(CaptureDisplay) = %d without getDisplayMedia()", s)
	}

	p.SetCaptureState(golet.CaptureMicrophone, golet.CaptureMuted)
	told("the microphone muted", golet.CaptureMicrophone, golet.CaptureMuted)
	waitFor(t, w, `track("audio") === "live muted"`)
	if s := p.CaptureState(golet.CaptureCamera); s != golet.CaptureActive {
		t.Errorf("muting the microphone: the camera %d", s)
	}

	// Muting the page's sound keeps the microphone muted, and the camera
	// live. Muting the camera next tells it once the page took both.
	p.SetAudioMuted(true)
	p.SetCaptureState(golet.CaptureCamera, golet.CaptureMuted)
	told("the camera muted", golet.CaptureCamera, golet.CaptureMuted)
	if s := p.CaptureState(golet.CaptureMicrophone); s != golet.CaptureMuted {
		t.Errorf("SetAudioMuted(true) unmuted the microphone: %d", s)
	}
	p.SetCaptureState(golet.CaptureCamera, golet.CaptureActive)
	told("the camera unmuted", golet.CaptureCamera, golet.CaptureActive)
	p.SetAudioMuted(false)
	p.SetCaptureState(golet.CaptureCamera, golet.CaptureMuted)
	told("the camera muted again", golet.CaptureCamera, golet.CaptureMuted)
	if s := p.CaptureState(golet.CaptureMicrophone); s != golet.CaptureMuted {
		t.Errorf("SetAudioMuted(false) unmuted the microphone: %d", s)
	}
	p.SetCaptureState(golet.CaptureCamera, golet.CaptureActive)
	told("the camera unmuted again", golet.CaptureCamera, golet.CaptureActive)
	// So does a burst of them, whose changes WebKit may tell on the way.
	for range 5 {
		p.SetAudioMuted(true)
		p.SetAudioMuted(false)
	}
	p.SetCaptureState(golet.CaptureCamera, golet.CaptureMuted)
	told("the camera muted after a burst", golet.CaptureCamera, golet.CaptureMuted)
	if s := p.CaptureState(golet.CaptureMicrophone); s != golet.CaptureMuted {
		t.Errorf("a burst of SetAudioMuted unmuted the microphone: %d", s)
	}
	p.SetCaptureState(golet.CaptureCamera, golet.CaptureActive)
	told("the camera unmuted after a burst", golet.CaptureCamera, golet.CaptureActive)

	p.SetCaptureState(golet.CaptureMicrophone, golet.CaptureActive)
	told("the microphone unmuted", golet.CaptureMicrophone, golet.CaptureActive)
	waitFor(t, w, `track("audio") === "live"`)

	p.SetCaptureState(golet.CaptureCamera, golet.CaptureNone)
	told("the camera stopped", golet.CaptureCamera, golet.CaptureNone)
	waitFor(t, w, `track("video") === "ended"`)
	if s := p.CaptureState(golet.CaptureMicrophone); s != golet.CaptureActive {
		t.Errorf("stopping the camera: the microphone %d", s)
	}
	p.SetCaptureState(golet.CaptureCamera, golet.CaptureMuted) // the page captures none
	mustEval(t, w, `stream.getAudioTracks()[0].stop()`)
	told("the page stopping the microphone", golet.CaptureMicrophone, golet.CaptureNone)

	// The page captures the camera again, unmuted.
	if got, err := p.EvalContext(ctx, `capture({video: true})`); err != nil || got != "ok" {
		t.Fatalf("capturing again: %v %v", got, err)
	}
	told("the camera capturing again", golet.CaptureCamera, golet.CaptureActive)
}
