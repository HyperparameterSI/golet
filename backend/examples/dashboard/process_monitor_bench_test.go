// Author: HyperparameterSI

package main

import (
	"encoding/json"
	"runtime"
	"testing"
	"time"
)

func benchmarkMonitor() *ProcessMonitor {
	m := &ProcessMonitor{startedAt: time.Now()}
	for i := range historyLimit {
		m.history.append(ProcessSample{Timestamp: int64(i + 1), Heap: 1 << 20, Goroutines: 4})
	}
	return m
}

func BenchmarkSnapshotHistory(b *testing.B) {
	m := benchmarkMonitor()
	b.ReportAllocs()
	for b.Loop() {
		_ = m.Snapshot()
	}
}

func BenchmarkSnapshotJSON(b *testing.B) {
	m := benchmarkMonitor()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := json.Marshal(m.Snapshot()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCollectFullHistory(b *testing.B) {
	m := benchmarkMonitor()
	runtime.GC()
	b.ReportAllocs()
	for b.Loop() {
		m.collectSample()
	}
}

// BenchmarkUpdateJSON measures the steady-state UI polling path: one new sample.
func BenchmarkUpdateJSON(b *testing.B) {
	m := benchmarkMonitor()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := json.Marshal(m.Updates(historyLimit - 1)); err != nil {
			b.Fatal(err)
		}
	}
}
