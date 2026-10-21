// Author: HyperparameterSI

package golet

import "testing"

func BenchmarkListenerSnapshot(b *testing.B) {
	var l listeners[func()]
	for range 16 {
		l.add(func() {}, false)
	}
	b.ReportAllocs()
	for b.Loop() {
		for _, fn := range l.snapshot() {
			fn()
		}
	}
}
