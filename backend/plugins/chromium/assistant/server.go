// Author: HyperparameterSI
package assistant

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/HyperparameterSI/golet/backend/plugins/chromium"
	"io/fs"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Server serves the local assistant with a per-launch bearer token, no CORS,
// and no credentials or browser bridge exposed to ordinary web pages.
type Server struct {
	mu            sync.Mutex
	browser       *chromium.Browser
	provider      Provider
	extensions    []Extension
	disabled      map[int]bool
	plans         map[string]*plan
	token, origin string
	http          *http.Server
}
type plan struct {
	Tab     int
	Page    Page
	Actions []Action
	Expires time.Time
}
type request struct {
	Question string `json:"question"`
	Tab      int    `json:"tab"`
	Share    bool   `json:"share"`
	Plan     string `json:"plan"`
	Index    int    `json:"index"`
	Enabled  bool   `json:"enabled"`
}

func nonce() string {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
func Start(assets fs.FS, p Provider, extensions []Extension) (*Server, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	s := &Server{provider: p, extensions: extensions, disabled: map[int]bool{}, plans: map[string]*plan{}, token: nonce(), origin: "http://" + ln.Addr().String()}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", s.api)
	files := http.FileServer(http.FS(assets))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		files.ServeHTTP(w, r)
	})
	s.http = &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != ln.Addr().String() {
			http.Error(w, "invalid host", 403)
			return
		}
		mux.ServeHTTP(w, r)
	}), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 75 * time.Second, IdleTimeout: 30 * time.Second}
	go s.http.Serve(ln)
	return s, nil
}
func (s *Server) URL() string              { return s.origin + "/#" + s.token }
func (s *Server) Bind(b *chromium.Browser) { s.mu.Lock(); s.browser = b; s.mu.Unlock() }
func (s *Server) Close() error             { return s.http.Close() }
func (s *Server) api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if subtle.ConstantTimeCompare([]byte(token), []byte(s.token)) != 1 || (r.Header.Get("Origin") != "" && r.Header.Get("Origin") != s.origin) {
		http.Error(w, "unauthorized", 403)
		return
	}
	s.mu.Lock()
	b := s.browser
	s.mu.Unlock()
	if b == nil {
		http.Error(w, "browser starting", 503)
		return
	}
	if r.URL.Path == "/api/state" && r.Method == "GET" {
		s.mu.Lock()
		disabled := map[int]bool{}
		for k, v := range s.disabled {
			disabled[k] = v
		}
		s.mu.Unlock()
		json.NewEncoder(w).Encode(map[string]any{"state": b.State(), "configured": s.provider.Endpoint != "" && s.provider.Model != "", "model": s.provider.Model, "extensions": s.extensions, "disabled": disabled})
		return
	}
	if r.Method != "POST" {
		http.Error(w, "POST required", 405)
		return
	}
	var req request
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 16000)).Decode(&req) != nil {
		http.Error(w, "invalid request", 400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 65*time.Second)
	defer cancel()
	var result any = map[string]bool{"ok": true}
	var err error
	switch r.URL.Path {
	case "/api/ask":
		if len(req.Question) == 0 || len(req.Question) > 8000 {
			err = errors.New("question must be 1–8000 bytes")
			break
		}
		var page *Page
		if req.Share {
			var snapshot Page
			snapshot, err = Snapshot(ctx, b, req.Tab)
			if err != nil {
				break
			}
			if strings.HasPrefix(snapshot.URL, s.origin) {
				err = errors.New("select a web page, not the assistant")
				break
			}
			page = &snapshot
		}
		var answer Answer
		answer, err = s.provider.Ask(ctx, req.Question, page)
		if err != nil {
			break
		}
		id := ""
		if page != nil && len(answer.Actions) > 0 {
			id = nonce()
			s.mu.Lock()
			for key, p := range s.plans {
				if time.Now().After(p.Expires) {
					delete(s.plans, key)
				}
			}
			if len(s.plans) >= 32 {
				err = errors.New("too many pending plans")
			} else {
				s.plans[id] = &plan{req.Tab, *page, answer.Actions, time.Now().Add(5 * time.Minute)}
			}
			s.mu.Unlock()
		}
		result = map[string]any{"answer": answer, "plan": id, "page": page}
	case "/api/execute":
		s.mu.Lock()
		p := s.plans[req.Plan]
		delete(s.plans, req.Plan)
		s.mu.Unlock()
		if p == nil || time.Now().After(p.Expires) || req.Index < 0 || req.Index >= len(p.Actions) {
			err = errors.New("plan expired; ask again")
			break
		}
		err = Apply(ctx, b, p.Tab, p.Page, p.Actions[req.Index])
	case "/api/extension/toggle", "/api/extension/run":
		if req.Index < 0 || req.Index >= len(s.extensions) {
			err = errors.New("extension not found")
			break
		}
		s.mu.Lock()
		if r.URL.Path == "/api/extension/toggle" {
			s.disabled[req.Index] = !req.Enabled
		}
		disabled := s.disabled[req.Index]
		s.mu.Unlock()
		if r.URL.Path == "/api/extension/run" {
			if disabled {
				err = errors.New("extension disabled")
			} else {
				err = s.extensions[req.Index].Run(ctx, b, req.Tab)
			}
		}
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(result)
}
