// Author: HyperparameterSI

package golet

import (
	"slices"
	"sync"
	"testing"
)

func TestListenerSnapshotMutation(t *testing.T) {
	var l listeners[func()]
	var calls []int
	remove := l.add(func() { calls = append(calls, 1) }, false)
	l.add(func() { calls = append(calls, 2) }, true)
	l.add(func() { calls = append(calls, 3) }, false)
	first := l.snapshot()
	remove()
	remove() // Unsubscribe is idempotent.
	l.add(func() { calls = append(calls, 4) }, false)
	for _, fn := range first {
		fn()
	}
	fire(&l)
	if !slices.Equal(calls, []int{1, 2, 3, 3, 4}) {
		t.Fatal(calls)
	}
}

func TestListenerOnceReentrant(t *testing.T) {
	var l listeners[func()]
	once, regular := 0, 0
	l.add(func() { once++; fire(&l) }, true)
	l.add(func() { regular++ }, false)
	fire(&l)
	if once != 1 || regular != 2 {
		t.Fatalf("once=%d regular=%d", once, regular)
	}
}

func TestListenerConcurrentMutation(t *testing.T) {
	var l listeners[func()]
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 100 {
				remove := l.add(func() {}, true)
				fire(&l)
				remove()
			}
		})
	}
	wg.Wait()
	if l.len() != 0 {
		t.Fatal("listeners leaked")
	}
}

// Cached snapshots must remain valid during reentrant subscription changes.
func TestListenerCachedSnapshotIsolation(t *testing.T) {
	var listeners listeners[func()]
	calls := 0
	remove := listeners.add(func() { calls++ }, false)
	before := listeners.snapshot()
	again := listeners.snapshot()
	if &before[0] != &again[0] {
		t.Fatal("stable subscription did not reuse snapshot")
	}
	remove()
	if listeners.cached != nil {
		t.Fatal("removed listener retained by cache")
	}
	listeners.add(func() { calls += 10 }, false)
	for _, fn := range before {
		fn()
	}
	fire(&listeners)
	if calls != 11 {
		t.Fatalf("snapshot changed during mutation: %d", calls)
	}
	listeners.clear()
	if listeners.cached != nil || len(listeners.snapshot()) != 0 {
		t.Fatal("clear retained callbacks")
	}
}
