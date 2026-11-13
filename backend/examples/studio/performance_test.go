// Author: HyperparameterSI
package main

import (
	"math"
	"runtime"
	"testing"
	"time"
)

func TestCPUDelta(t *testing.T) {
	a := raw{cpuSeconds: 10, ticks: [4]uint32{10, 10, 60, 20}, processOK: true, systemOK: true}
	b := raw{cpuSeconds: 13, ticks: [4]uint32{20, 30, 110, 40}, processOK: true, systemOK: true}
	p, s := cpuDelta(a, b, 2, true)
	if p == nil || s == nil || *p != 150 || *s != 50 {
		t.Fatal(p, s)
	}
	p, s = cpuDelta(a, b, 2, false)
	if p != nil || s != nil {
		t.Fatal("first sample must be unavailable")
	}
	p, s = cpuDelta(a, raw{}, 2, true)
	if p != nil || s != nil {
		t.Fatal("failure must not become zero")
	}
	a.ticks = [4]uint32{math.MaxUint32 - 1, 0, 0, 0}
	b.ticks = [4]uint32{2, 0, 4, 0}
	_, s = cpuDelta(a, b, 2, true)
	if s == nil || *s != 50 {
		t.Fatal("tick rollover", s)
	}
}

type fakeSource struct{}

func (fakeSource) Read() raw { return raw{} }
func (fakeSource) Close()    {}
func TestBoundedHistory(t *testing.T) {
	m := &Performance{source: fakeSource{}}
	start := time.Unix(1000, 0)
	for i := 0; i < 130; i++ {
		m.take(start.Add(time.Duration(i) * 2 * time.Second))
	}
	snapshot := m.Snapshot()
	if len(snapshot.Samples) != 120 || snapshot.Samples[0].At != start.Add(20*time.Second).UnixMilli() {
		t.Fatal("history not bounded")
	}
	snapshot.Samples[0].Heap = 0
	if m.Snapshot().Samples[0].Heap == 0 {
		t.Fatal("snapshot aliases history")
	}
}

func TestSnapshotOwnsNullableValues(t *testing.T) {
	cpu := 12.5
	rss := uint64(100)
	m := &Performance{samples: []PerformanceSample{{ProcessCPU: &cpu, RSS: &rss, Errors: []string{"unavailable"}}}}
	s := m.Snapshot()
	*s.Samples[0].ProcessCPU = 0
	*s.Samples[0].RSS = 0
	s.Samples[0].Errors[0] = "changed"
	got := m.Snapshot()
	if *got.Samples[0].ProcessCPU != 12.5 || *got.Samples[0].RSS != 100 || got.Samples[0].Errors[0] != "unavailable" {
		t.Fatal("snapshot exposed sampler state")
	}
}

// TestNativeSample checks availability without requiring a time-based CPU interval.
func TestNativeSample(t *testing.T) {
	s := newSource()
	defer s.Close()
	r := s.Read()
	if runtime.GOOS == "darwin" {
		if r.rss == nil || r.total == nil || r.used == nil || !r.systemOK || !r.processOK {
			t.Fatalf("native counters unavailable: %v", r.errors)
		}
		if *r.used > *r.total || *r.rss == 0 {
			t.Fatal("invalid native memory counters")
		}
	}
}
