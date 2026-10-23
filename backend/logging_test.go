// Author: HyperparameterSI

//go:build !js

package golet

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestLogFiltering(t *testing.T) {
	t.Cleanup(func() { SetLogHandler(LogWarn, nil) })
	var records []LogRecord
	SetLogHandler(LogWarn, func(r LogRecord) { records = append(records, r) })
	reportLog(LogRecord{Level: LogDebug, Event: "test.debug"})
	failure := errors.New("failure")
	reportLog(LogRecord{Level: LogError, Event: "test.error", WindowID: 42, Method: "Service.Method", Err: failure})
	if len(records) != 1 {
		t.Fatalf("got %d records", len(records))
	}
	r := records[0]
	if r.Time.IsZero() || r.Event != "test.error" || r.WindowID != 42 || r.Method != "Service.Method" || r.Err != failure {
		t.Fatalf("lost diagnostic context: %+v", r)
	}
}

func TestLogDefault(t *testing.T) {
	var output bytes.Buffer
	old := log.Writer()
	log.SetOutput(&output)
	t.Cleanup(func() { log.SetOutput(old); SetLogHandler(LogWarn, nil) })
	SetLogHandler(LogDebug, func(LogRecord) { t.Fatal("replaced handler invoked") })
	SetLogHandler(LogDebug, nil)
	reportLog(LogRecord{Level: LogInfo, Event: "test.quiet"})
	if output.Len() != 0 {
		t.Fatal("default logger emitted info")
	}
	reportLog(LogRecord{Level: LogWarn, Event: "test.warning", WindowID: 17, Err: errors.New("first\nsecond")})
	got := output.String()
	if !strings.Contains(got, "level=WARN event=test.warning window_id=17") || !strings.Contains(got, `error="first\nsecond"`) {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestLogConcurrentReplacement(t *testing.T) {
	t.Cleanup(func() { SetLogHandler(LogWarn, nil) })
	var count atomic.Int64
	handler := func(LogRecord) { count.Add(1) }
	SetLogHandler(LogDebug, handler)
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 100 {
				SetLogHandler(LogDebug, handler)
				reportLog(LogRecord{Level: LogDebug, Event: "test.concurrent"})
			}
		})
	}
	wg.Wait()
	if count.Load() != 800 {
		t.Fatalf("lost records: %d", count.Load())
	}
}

func TestLogDisabledAllocations(t *testing.T) {
	SetLogHandler(LogWarn, nil)
	if n := testing.AllocsPerRun(100, func() {
		reportLog(LogRecord{Level: LogDebug, Event: "ipc.started", WindowID: 1, Method: "Service.Method"})
	}); n != 0 {
		t.Fatalf("disabled logging allocates: %v", n)
	}
}

func TestLogIPCPanic(t *testing.T) {
	bindForTest(t, "LogSvc", &testService{})
	w, fw := testWindow(t, WindowOptions{})
	var mu sync.Mutex
	var records []LogRecord
	SetLogHandler(LogDebug, func(r LogRecord) {
		mu.Lock()
		records = append(records, r)
		mu.Unlock()
	})
	t.Cleanup(func() { SetLogHandler(LogWarn, nil) })
	reply := call(t, fw, 1, "LogSvc.Panic")
	if reply["ok"] != false || reply["e"] != "panic: boom" {
		t.Fatalf("unexpected reply: %v", reply)
	}
	mu.Lock()
	defer mu.Unlock()
	for _, r := range records {
		if r.Event == "ipc.panic" {
			if r.WindowID != w.ID() || r.Method != "LogSvc.Panic" || r.Err == nil || !strings.Contains(r.Stack, "(*testService).Panic") {
				t.Fatalf("incomplete panic context: %+v", r)
			}
			return
		}
	}
	t.Fatal("missing panic diagnostic")
}

func BenchmarkLogDisabled(b *testing.B) {
	SetLogHandler(LogWarn, nil)
	b.ReportAllocs()
	for b.Loop() {
		reportLog(LogRecord{Level: LogDebug, Event: "ipc.started", WindowID: 1, Method: "Service.Method"})
	}
}

type canceledLogService struct{}

func (*canceledLogService) Cancel() error { return context.Canceled }

func TestLogIPCCancellation(t *testing.T) {
	bindForTest(t, "LogCancel", &canceledLogService{})
	_, fw := testWindow(t, WindowOptions{})
	var mu sync.Mutex
	var records []LogRecord
	SetLogHandler(LogDebug, func(r LogRecord) {
		mu.Lock()
		records = append(records, r)
		mu.Unlock()
	})
	t.Cleanup(func() { SetLogHandler(LogWarn, nil) })
	reply := call(t, fw, 1, "LogCancel.Cancel")
	if reply["ok"] != false {
		t.Fatalf("unexpected reply: %v", reply)
	}
	mu.Lock()
	defer mu.Unlock()
	found := false
	for _, r := range records {
		if r.Method != "LogCancel.Cancel" {
			continue
		}
		if r.Level >= LogWarn {
			t.Fatalf("cancellation emitted warning: %+v", r)
		}
		if r.Event == "ipc.canceled" {
			found = true
		}
	}
	if !found {
		t.Fatal("missing cancellation diagnostic")
	}
}
