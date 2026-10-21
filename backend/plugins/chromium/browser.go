// Author: HyperparameterSI

// Package chromium runs the optional CEF browser host in its own process.
// It provides a complete Chromium window, not a system-WebView fallback.
// The native host owns CEF's UI thread; Go controls it over private stdio pipes.
package chromium

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
)

var ErrClosed = errors.New("chromium: browser closed")

// Options selects a separately installed CEF host and a persistent profile.
// A profile directory must not be shared by simultaneously running hosts.
type Options struct {
	// EphemeralDemo creates a disposable development profile. On macOS it uses
	// a mock keychain to suppress authorization prompts. Data is not protected
	// by Keychain; do not use for real credentials. Normal exit removes the profile.
	// It cannot be combined with ProfileDir. The default remains secure storage.
	EphemeralDemo bool
	// Style selects the window implementation. The zero value uses StyleAlloy.
	Style BrowserStyle
	// ChromeStyle enables the standard Chrome window.
	// Deprecated: use Style: StyleChrome instead.
	ChromeStyle bool
	// Extensions contains explicitly trusted unpacked extension directories.
	Extensions []string
	// AssistantURL is a local assistant page exposed by the embedding application.
	AssistantURL string
	HostPath     string
	URL          string
	ProfileDir   string
	// Stderr receives CEF diagnostics; nil uses os.Stderr.
	Stderr io.Writer
}

// Tab identifies one open page within the browser window.
type Tab struct {
	BrowserID int    `json:"browserId"`
	URL       string `json:"url"`
	Title     string `json:"title"`
}

// State describes the focused tab and open tabs. Events are snapshots, not a durable log.
type State struct {
	Tabs         []Tab  `json:"tabs"`
	BrowserID    int    `json:"browserId"`
	URL          string `json:"url"`
	Title        string `json:"title"`
	Loading      bool   `json:"loading"`
	CanGoBack    bool   `json:"canGoBack"`
	CanGoForward bool   `json:"canGoForward"`
}

type message struct {
	DevToolsMethod string          `json:"devToolsMethod,omitempty"`
	Params         json.RawMessage `json:"params,omitempty"`
	Result         json.RawMessage `json:"result,omitempty"`
	BrowserID      int             `json:"browserId,omitempty"`
	ID             uint64          `json:"id,omitempty"`
	Method         string          `json:"method,omitempty"`
	URL            string          `json:"url,omitempty"`
	Event          string          `json:"event,omitempty"`
	Error          string          `json:"error,omitempty"`
	State          State           `json:"state,omitempty"`
}

// Browser owns a host process. Methods are safe to call concurrently. Call Close
// or cancel the context passed to Launch to release the host and its children.
type Browser struct {
	cmd       *exec.Cmd
	input     io.WriteCloser
	writes    chan []byte
	mu        sync.Mutex
	nextID    uint64
	pending   map[uint64]chan message
	state     State
	exitErr   error
	readDone  chan struct{}
	readErr   error
	ready     chan struct{}
	readyOnce sync.Once
	done      chan struct{}
	events    chan State
}

// DefaultHostPath returns the optional host's installation location.
// GOLET_CHROMIUM_HOST can point to a host shipped alongside an application.
func DefaultHostPath() (string, error) {
	if p := os.Getenv("GOLET_CHROMIUM_HOST"); p != "" {
		return filepath.Abs(p)
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	root := filepath.Join(cache, "golet", "chromium", "host-v1", runtime.GOOS+"-"+runtime.GOARCH)
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(root, "Golet Browser.app", "Contents", "MacOS", "Golet Browser"), nil
	case "windows":
		return filepath.Join(root, "golet-chromium.exe"), nil
	case "linux":
		return filepath.Join(root, "golet-chromium"), nil
	default:
		return "", fmt.Errorf("chromium: unsupported platform %s", runtime.GOOS)
	}
}

// Launch starts the CEF host and waits for its first browser. ctx controls the
// entire session. It never downloads binaries implicitly.
func Launch(ctx context.Context, opts Options) (*Browser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	extra, err := extensionArguments(opts)
	if err != nil {
		return nil, err
	}
	target, err := navigationURL(opts.URL)
	if err != nil {
		return nil, err
	}
	host := opts.HostPath
	if host == "" {
		host, err = DefaultHostPath()
		if err != nil {
			return nil, err
		}
	}
	host, err = filepath.Abs(host)
	if err != nil {
		return nil, err
	}
	if info, e := os.Stat(host); e != nil || info.IsDir() {
		return nil, fmt.Errorf("chromium: host not installed at %s; run the Chromium setup command", host)
	}
	profile, profileArgs, cleanupProfile, err := prepareProfile(opts)
	if err != nil {
		return nil, err
	}
	started := false
	defer func() {
		if !started {
			cleanupProfile()
		}
	}()
	cmd := exec.Command(host, "--url="+target, "--cache-path="+profile)
	cmd.Args = append(cmd.Args, extra...)
	cmd.Args = append(cmd.Args, profileArgs...)
	if opts.AssistantURL != "" {
		cmd.Args = append(cmd.Args, "--assistant-url="+opts.AssistantURL)
	}
	cmd.Dir = filepath.Dir(host)
	cmd.WaitDelay = 5 * time.Second
	prepareProcess(cmd)
	cmd.Stderr = opts.Stderr
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}
	input, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	output, outputWriter, err := os.Pipe()
	if err != nil {
		input.Close()
		return nil, err
	}
	cmd.Stdout = outputWriter
	b := &Browser{
		cmd:      cmd,
		input:    input,
		pending:  make(map[uint64]chan message),
		readDone: make(chan struct{}),
		ready:    make(chan struct{}),
		done:     make(chan struct{}),
		events:   make(chan State, 1),
		writes:   make(chan []byte, 16),
	}
	if err = cmd.Start(); err != nil {
		input.Close()
		output.Close()
		outputWriter.Close()
		return nil, err
	}
	started = true
	_ = outputWriter.Close()
	go b.read(output)
	go func() {
		err := cmd.Wait()
		_ = input.Close()
		_ = output.Close()
		<-b.readDone
		b.mu.Lock()
		b.exitErr = errors.Join(err, b.readErr)
		b.mu.Unlock()
		cleanupProfile()
		close(b.events)
		close(b.done)
	}()
	go b.write()
	go func() {
		select {
		case <-ctx.Done():
			closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = b.Close(closeCtx)
		case <-b.done:
		}
	}()
	startup, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	select {
	case <-b.ready:
		return b, nil
	case <-b.done:
		return nil, errors.Join(errors.New("chromium: host exited before ready"), b.Wait())
	case <-startup.Done():
		_ = killProcess(cmd)
		<-b.done
		return nil, fmt.Errorf("chromium: startup: %w", startup.Err())
	}
}

func navigationURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "about:blank", nil
	}
	if raw == "about:blank" {
		return raw, nil
	}
	u, err := url.Parse(raw)
	if err != nil || u == nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil {
		return "", fmt.Errorf("chromium: expected an http(s) URL or about:blank")
	}
	if len(raw) > 16384 {
		return "", errors.New("chromium: URL exceeds 16 KiB")
	}
	return u.String(), nil
}

func (b *Browser) read(r io.Reader) {
	defer close(b.readDone)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 1<<20)
	for scanner.Scan() {
		var m message
		if json.Unmarshal(scanner.Bytes(), &m) != nil {
			continue
		}
		b.mu.Lock()
		if m.Event == "state" || m.Event == "ready" {
			b.state = m.State
			select {
			case b.events <- cloneState(m.State):
			default:
				select {
				case <-b.events:
				default:
				}
				select {
				case b.events <- cloneState(m.State):
				default:
				}
			}
		}
		reply := b.pending[m.ID]
		b.mu.Unlock()
		if reply != nil {
			select {
			case reply <- m:
			default:
			}
		}
		if m.Event == "ready" {
			b.readyOnce.Do(func() { close(b.ready) })
		}
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, os.ErrClosed) {
		b.mu.Lock()
		b.readErr = err
		b.mu.Unlock()
		_ = killProcess(b.cmd)
	}

}

// Events returns a bounded stream of the newest focused-tab state. Slow readers
// skip intermediate states. The channel closes after the host exits.
func (b *Browser) Events() <-chan State { return b.events }

// State returns the most recently reported focused-tab state.
func (b *Browser) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return cloneState(b.state)
}

// Done closes once the host has exited and its pipes have been released.
func (b *Browser) Done() <-chan struct{} { return b.done }

// Wait waits for process exit and returns its error, if any.
func (b *Browser) Wait() error {
	<-b.done
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.exitErr
}

func (b *Browser) call(ctx context.Context, method, target string, tab ...int) error {
	request := message{Method: method, URL: target}
	if len(tab) > 0 {
		request.BrowserID = tab[0]
	}
	_, err := b.request(ctx, request)
	return err
}

func (b *Browser) request(ctx context.Context, request message) (json.RawMessage, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	b.mu.Lock()
	b.nextID++
	id := b.nextID
	reply := make(chan message, 1)
	b.pending[id] = reply
	b.mu.Unlock()
	defer func() { b.mu.Lock(); delete(b.pending, id); b.mu.Unlock() }()
	request.ID = id
	data, _ := json.Marshal(request)
	select {
	case b.writes <- append(data, '\n'):
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-b.done:
		return nil, ErrClosed
	}
	select {
	case m := <-reply:
		if m.Error != "" {
			return nil, errors.New("chromium: " + m.Error)
		}
		return m.Result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-b.done:
		return nil, ErrClosed
	}
}

// Navigate requests navigation; observe Events for loading completion.
func (b *Browser) Navigate(ctx context.Context, target string) error {
	u, err := navigationURL(target)
	if err != nil {
		return err
	}
	return b.call(ctx, "navigate", u)
}

// Back navigates to the previous history entry.
func (b *Browser) Back(ctx context.Context) error { return b.call(ctx, "back", "") }

// Forward navigates to the next history entry.
func (b *Browser) Forward(ctx context.Context) error { return b.call(ctx, "forward", "") }

// Reload refreshes the focused tab.
func (b *Browser) Reload(ctx context.Context) error { return b.call(ctx, "reload", "") }

// Stop cancels loading in the focused tab.
func (b *Browser) Stop(ctx context.Context) error { return b.call(ctx, "stop", "") }

// Close asks CEF to close all tabs and flush its profile. If ctx expires the host
// is killed. Close may be called more than once.
func (b *Browser) Close(ctx context.Context) error {
	select {
	case <-b.done:
		return b.Wait()
	default:
	}
	err := b.call(ctx, "close", "")
	if err != nil && !errors.Is(err, ErrClosed) {
		_ = killProcess(b.cmd)
	}
	select {
	case <-b.done:
		return b.Wait()
	case <-ctx.Done():
		_ = killProcess(b.cmd)
		return ctx.Err()
	}
}

// write is the single pipe writer. Bounded queuing keeps a stalled native host
// from blocking callers beyond their contexts or creating unbounded goroutines.
func (b *Browser) write() {
	for {
		select {
		case data := <-b.writes:
			if _, err := b.input.Write(data); err != nil {
				_ = killProcess(b.cmd)
				return
			}
		case <-b.done:
			return
		}
	}
}

// NewTab creates and selects a tab. An empty URL opens the bundled new-tab page.
// Observe Events for its BrowserID and loading completion.
func (b *Browser) NewTab(ctx context.Context, target string) error {
	url, err := navigationURL(target)
	if err != nil {
		return err
	}
	return b.call(ctx, "newtab", url)
}

// SelectTab focuses an existing tab using the BrowserID from State or Events.
func (b *Browser) SelectTab(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("chromium: invalid tab ID")
	}
	return b.call(ctx, "select-tab", "", id)
}

// CloseTab closes one tab; closing the last tab exits the host.
func (b *Browser) CloseTab(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("chromium: invalid tab ID")
	}
	return b.call(ctx, "close-tab", "", id)
}

func cloneState(state State) State { state.Tabs = slices.Clone(state.Tabs); return state }
