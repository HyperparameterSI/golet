// Author: HyperparameterSI

package main

import (
	"sync"
	"testing"
)

func TestCPUPercent(t *testing.T) {
	for _, tt := range []struct{ previous, current, elapsed, want float64 }{{1, 1.5, 1, 50}, {1, 3, 1, 200}, {1, 1, 1, 0}, {2, 1, 1, 0}, {1, 2, 0, 0}} {
		if got := calculateCPUPercent(tt.previous, tt.current, tt.elapsed); got != tt.want {
			t.Fatalf("%+v: got %v", tt, got)
		}
	}
}
func TestHistoryBoundAndSnapshotIsolation(t *testing.T) {
	m := benchmarkMonitor()
	m.collectSample()
	snapshot := m.Snapshot()
	if len(snapshot.Samples) != historyLimit || snapshot.Samples[0].Sequence != 2 {
		t.Fatal("history not bounded")
	}
	snapshot.Samples[0].Heap = 123
	if m.Snapshot().Samples[0].Heap == 123 {
		t.Fatal("snapshot aliases internal history")
	}
	last := snapshot.Samples[len(snapshot.Samples)-1]
	if last.Heap == 0 || last.Goroutines == 0 || last.CPUAvailable {
		t.Fatalf("invalid first sample: %+v", last)
	}
	updates := m.Updates(historyLimit)
	if len(updates.Samples) != 1 || updates.Samples[0].Sequence != historyLimit+1 {
		t.Fatal("incremental cursor failed")
	}
	if len(m.Updates(historyLimit+1).Samples) != 0 {
		t.Fatal("duplicate updates")
	}
	if len(m.Updates(historyLimit+100).Samples) != historyLimit {
		t.Fatal("stale process cursor did not reset")
	}
}

func TestHistoryWrapAndClockRollback(t *testing.T) {
	var history sampleHistory
	for i := range historyLimit*3 + 7 {
		history.append(ProcessSample{Timestamp: int64(-i)})
	}
	samples := history.since(0)
	for i, sample := range samples {
		want := uint64(historyLimit*2 + 8 + i)
		if sample.Sequence != want {
			t.Fatalf("sequence %d, want %d", sample.Sequence, want)
		}
	}
	if len(history.since(history.sequence)) != 0 {
		t.Fatal("current cursor must be empty")
	}
}

func TestMonitorConcurrentCloseAndRead(t *testing.T) {
	monitor := newProcessMonitor()
	var group sync.WaitGroup
	for range 8 {
		group.Go(func() { _ = monitor.Updates(0); monitor.close() })
	}
	group.Wait()
}
