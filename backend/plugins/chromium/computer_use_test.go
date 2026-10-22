// Author: HyperparameterSI
package chromium

import (
	"context"
	"os"
	"slices"
	"testing"
	"time"
)

func TestComputerUseRequiresExplicitShellSupport(t *testing.T) {
	args, err := extensionArguments(Options{})
	if err != nil || slices.Contains(args, "--golet-computer-use") {
		t.Fatalf("default enabled: %v %v", args, err)
	}
	args, err = extensionArguments(Options{EnableComputerUse: true})
	if err != nil || !slices.Contains(args, "--golet-computer-use") {
		t.Fatalf("missing opt-in: %v %v", args, err)
	}
	for _, opts := range []Options{{EnableComputerUse: true, AppMode: true}, {EnableComputerUse: true, Style: StyleChrome}} {
		if _, err := extensionArguments(opts); err == nil {
			t.Fatal("unsupported shell accepted")
		}
	}
}
func TestNativeTargetIdentity(t *testing.T) {
	for _, invalid := range []string{"", "pid", "window"} {
		t.Run(invalid, func(t *testing.T) {
			t.Setenv("GOLET_CHROMIUM_TEST_HOST", "1")
			t.Setenv("GOLET_TARGET_INVALID", invalid)
			exe, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			b, err := Launch(ctx, Options{HostPath: exe, ProfileDir: t.TempDir()})
			if err != nil {
				t.Fatal(err)
			}
			defer b.Close(ctx)
			target, err := b.ComputerUseTarget(ctx)
			if invalid != "" {
				if err == nil {
					t.Fatal("invalid native identity accepted")
				}
				return
			}
			if err != nil || target.PID != b.cmd.Process.Pid || target.WindowID != 123 || target.Enabled {
				t.Fatalf("target=%+v error=%v", target, err)
			}
		})
	}
}
