// Author: HyperparameterSI

package e2e

import (
	"testing"
	"time"

	"github.com/HyperparameterSI/golet/backend"
)

// TestNetwork reads the network status through the real backend. The
// computers that run the tests are online (they downloaded the module), so
// the first Status, which starts the monitor on macOS and waits for its
// first path, says so. Changes are covered through internal/fake.
func TestNetwork(t *testing.T) {
	start := time.Now()
	s := golet.Network.Status()
	t.Logf("%+v in %v", s, time.Since(start))
	if !s.Online {
		t.Errorf("Status = %+v, want online", s)
	}
	golet.Network.OnChanged(func(golet.NetworkStatus) {})()
	if again := golet.Network.Status(); again != s {
		t.Logf("the status changed to %+v", again)
	}
}
