// Author: HyperparameterSI
package main

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type PerformanceSample struct {
	At           int64    `json:"at"`
	ProcessCPU   *float64 `json:"processCPU"`
	SystemCPU    *float64 `json:"systemCPU"`
	RSS          *uint64  `json:"rss"`
	SystemUsed   *uint64  `json:"systemUsed"`
	SystemTotal  *uint64  `json:"systemTotal"`
	Heap         uint64   `json:"heap"`
	HeapReserved uint64   `json:"heapReserved"`
	GoReserved   uint64   `json:"goReserved"`
	Goroutines   int      `json:"goroutines"`
	GC           uint32   `json:"gc"`
	GCPauseMS    float64  `json:"gcPauseMS"`
	Errors       []string `json:"errors"`
}
type PerformanceSnapshot struct {
	Samples    []PerformanceSample `json:"samples"`
	Cores      int                 `json:"cores"`
	Platform   string              `json:"platform"`
	IntervalMS int                 `json:"intervalMS"`
}
type raw struct {
	cpuSeconds          float64
	ticks               [4]uint32
	processOK, systemOK bool
	rss, total, used    *uint64
	errors              []string
}
type source interface {
	Read() raw
	Close()
}
type Performance struct {
	started  atomic.Bool
	mu       sync.RWMutex
	samples  []PerformanceSample
	previous raw
	at       time.Time
	source   source
}

// NewPerformance allocates no native resources. Run owns and releases its sampler.
func NewPerformance() *Performance { return &Performance{} }
func (m *Performance) run(ctx context.Context) error {
	if ctx == nil {
		return errors.New("metrics: nil context")
	}
	if !m.started.CompareAndSwap(false, true) {
		return errors.New("metrics: monitor is single-use")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if m.source == nil {
		m.source = newSource()
	}
	defer m.source.Close()
	m.take(time.Now())
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-t.C:
			m.take(now)
		}
	}
}
func (m *Performance) take(now time.Time) {
	r := m.source.Read()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	s := PerformanceSample{At: now.UnixMilli(), RSS: r.rss, SystemUsed: r.used, SystemTotal: r.total, Heap: mem.HeapAlloc, HeapReserved: mem.HeapSys, GoReserved: mem.Sys, Goroutines: runtime.NumGoroutine(), GC: mem.NumGC, Errors: r.errors}
	if mem.NumGC > 0 {
		s.GCPauseMS = float64(mem.PauseNs[(mem.NumGC-1)%256]) / 1e6
	}
	s.ProcessCPU, s.SystemCPU = cpuDelta(m.previous, r, now.Sub(m.at).Seconds(), !m.at.IsZero())
	m.previous = r
	m.at = now
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.samples) == 120 {
		copy(m.samples, m.samples[1:])
		m.samples = m.samples[:119]
	}
	m.samples = append(m.samples, s)
}
func cpuDelta(prev, next raw, seconds float64, ready bool) (*float64, *float64) {
	var process, system *float64
	if ready && seconds > 0 && prev.processOK && next.processOK && next.cpuSeconds >= prev.cpuSeconds {
		v := (next.cpuSeconds - prev.cpuSeconds) / seconds * 100
		process = &v
	}
	if ready && prev.systemOK && next.systemOK {
		var total, idle uint64
		for i := range next.ticks {
			d := uint64(next.ticks[i] - prev.ticks[i])
			total += d
			if i == 2 {
				idle = d
			}
		}
		if total > 0 {
			v := 100 * float64(total-idle) / float64(total)
			system = &v
		}
	}
	return process, system
}
func (m *Performance) Snapshot() PerformanceSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	samples := append([]PerformanceSample{}, m.samples...)
	for i := range samples {
		s := &samples[i]
		s.ProcessCPU, s.SystemCPU = copyValue(s.ProcessCPU), copyValue(s.SystemCPU)
		s.RSS, s.SystemUsed, s.SystemTotal = copyValue(s.RSS), copyValue(s.SystemUsed), copyValue(s.SystemTotal)
		s.Errors = append([]string(nil), s.Errors...)
	}
	return PerformanceSnapshot{samples, runtime.NumCPU(), runtime.GOOS + "/" + runtime.GOARCH, 2000}
}
func copyValue[T any](p *T) *T {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}
