// Author: HyperparameterSI

package main

import (
	"os"
	"runtime"
	"slices"
	"sync"
	"time"
)

// ProcessMonitor collects process metrics off the UI thread, with bounded history.
type ProcessMonitor struct {
	mu                   sync.RWMutex
	history              sampleHistory
	stopOnce             sync.Once
	gcPauseP95MS         float64
	startedAt            time.Time
	stop                 chan struct{}
	done                 chan struct{}
	previousSampleAt     time.Time
	previousCPU          float64
	previousCPUAvailable bool
	previousAlloc        uint64
	previousGC           uint32
}

func newProcessMonitor() *ProcessMonitor {
	m := &ProcessMonitor{startedAt: time.Now(), stop: make(chan struct{}), done: make(chan struct{})}
	m.collectSample()
	go func() {
		defer close(m.done)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.collectSample()
			case <-m.stop:
				return
			}
		}
	}()
	return m
}

// close is idempotent and waits for the sampling goroutine to release its timer.
func (m *ProcessMonitor) close() {
	m.stopOnce.Do(func() { close(m.stop) })
	<-m.done
}

// Snapshot returns real metrics for the Go process, excluding WebView children.
func (m *ProcessMonitor) Snapshot() ProcessSnapshot {
	return m.Updates(0)
}

// Updates returns observations newer than sequence. Clients keep a bounded
// history locally, avoiding repeated encoding and transmission of old samples.
func (m *ProcessMonitor) Updates(sequence uint64) ProcessSnapshot {
	m.mu.RLock()
	samples := m.history.since(sequence)
	m.mu.RUnlock()
	return ProcessSnapshot{
		Samples:    samples,
		PID:        os.Getpid(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		GoVersion:  runtime.Version(),
		CPUs:       runtime.NumCPU(),
		GOMAXPROCS: runtime.GOMAXPROCS(0),
		Uptime:     time.Since(m.startedAt).Seconds(),
		StartedAt:  m.startedAt.UnixMilli(),
		IntervalMS: 1000,
	}
}

func calculateCPUPercent(previous, current, elapsed float64) float64 {
	if elapsed <= 0 || current < previous {
		return 0
	}
	return (current - previous) / elapsed * 100
}

func (m *ProcessMonitor) collectSample() {
	start := time.Now()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	cpu, err := processCPUSeconds()
	threads, _ := runtime.ThreadCreateProfile(nil)
	s := ProcessSample{
		Timestamp:      start.UnixMilli(),
		Heap:           mem.HeapAlloc,
		HeapInuse:      mem.HeapInuse,
		HeapIdle:       mem.HeapIdle,
		HeapReleased:   mem.HeapReleased,
		RuntimeMemory:  mem.Sys - mem.HeapReleased,
		Stack:          mem.StackInuse,
		Objects:        mem.HeapObjects,
		TotalAlloc:     mem.TotalAlloc,
		Goroutines:     runtime.NumGoroutine(),
		ThreadRecords:  threads,
		GCCycles:       mem.NumGC,
		GCTotalPauseMS: float64(mem.PauseTotalNs) / 1e6,
		NextGC:         mem.NextGC,
		GCCPUFraction:  mem.GCCPUFraction,
	}
	if !m.previousSampleAt.IsZero() {
		elapsed := start.Sub(m.previousSampleAt).Seconds()
		s.CPUAvailable = err == nil && m.previousCPUAvailable
		if s.CPUAvailable {
			s.CPU = calculateCPUPercent(m.previousCPU, cpu, elapsed)
		}
		s.AllocRate = float64(mem.TotalAlloc-m.previousAlloc) / elapsed
		s.GCDelta = mem.NumGC - m.previousGC
	}
	if mem.NumGC > 0 {
		s.GCPauseMS = float64(mem.PauseNs[(mem.NumGC-1)%256]) / 1e6
		// GC statistics only change after a collection. Sorting a fixed stack
		// buffer avoids per-sample allocations and repeated percentile work.
		if mem.NumGC != m.previousGC {
			var pauses [256]uint64
			count := min(int(mem.NumGC), len(pauses))
			copy(pauses[:count], mem.PauseNs[:count])
			slices.Sort(pauses[:count])
			m.gcPauseP95MS = float64(pauses[(count*95+99)/100-1]) / 1e6
		}
		s.GCPauseP95MS = m.gcPauseP95MS
	}
	s.SampleDurationMS = float64(time.Since(start).Nanoseconds()) / 1e6
	m.previousSampleAt = start
	m.previousCPU = cpu
	m.previousCPUAvailable = err == nil
	m.previousAlloc = mem.TotalAlloc
	m.previousGC = mem.NumGC
	m.mu.Lock()
	m.history.append(s)
	m.mu.Unlock()
}
