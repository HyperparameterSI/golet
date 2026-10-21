// Author: HyperparameterSI
package main

import "testing"

func TestBrowserDemoRequiresLaunch(t *testing.T) {
	demo := &BrowserDemo{}
	if err := demo.Command("navigate", "https://go.dev"); err == nil {
		t.Fatal("command accepted without host")
	}
	if err := demo.Start("invalid-style", "https://go.dev"); err == nil {
		demo.shutdown()
		t.Fatal("invalid style launched")
	}
	if demo.browser != nil {
		t.Fatal("failed launch retained host")
	}
}
