// Author: HyperparameterSI

package golet

import (
	"log"
	"sync/atomic"
	"time"
)

// LogLevel is the severity of a framework diagnostic.
type LogLevel uint8

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
)

func (l LogLevel) String() string {
	switch l {
	case LogDebug:
		return "DEBUG"
	case LogInfo:
		return "INFO"
	case LogWarn:
		return "WARN"
	case LogError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogRecord describes an application, window, IPC or protocol diagnostic.
// Event is a stable dotted name. WindowID and Method are zero when inapplicable.
// Arguments, results, bridge tokens and request URLs are not included. Errors
// and panic stacks may contain application data; filter them before exporting.
type LogRecord struct {
	Time     time.Time
	Level    LogLevel
	Event    string
	WindowID int
	Method   string
	Err      error
	Stack    string
}

type logConfig struct {
	minimum LogLevel
	handler func(LogRecord)
}

var frameworkLog atomic.Pointer[logConfig]

// SetLogHandler replaces the handler for core framework diagnostics. A nil
// handler restores the default: warnings and errors through the standard log
// package. Detailed lifecycle events are opt-in; no background logger is started.
//
// The handler runs synchronously on the reporting goroutine, possibly the UI
// thread. It must be concurrency-safe, return promptly and not call UI methods
// or panic. Replacing the handler is safe while the application is running;
// already-started calls may still finish on the previous handler.
func SetLogHandler(minimum LogLevel, handler func(LogRecord)) {
	if handler == nil {
		frameworkLog.Store(nil)
		return
	}
	frameworkLog.Store(&logConfig{minimum: minimum, handler: handler})
}

func reportLog(record LogRecord) {
	config := frameworkLog.Load()
	if config != nil {
		if record.Level >= config.minimum {
			record.Time = time.Now()
			config.handler(record)
		}
		return
	}
	if record.Level < LogWarn {
		return
	}
	errorText := ""
	if record.Err != nil {
		errorText = record.Err.Error()
	}
	log.Printf("golet: level=%s event=%s window_id=%d method=%q error=%q stack=%q",
		record.Level, record.Event, record.WindowID, record.Method, errorText, record.Stack)
}
