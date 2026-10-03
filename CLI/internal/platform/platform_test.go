package platform

import (
	"runtime"
	"testing"
)

func TestDetectReturnsGoRuntimeTarget(t *testing.T) {
	if got := Detect(); got.OS != runtime.GOOS || got.Arch != runtime.GOARCH {
		t.Fatalf("Detect() = %#v, want %s/%s", got, runtime.GOOS, runtime.GOARCH)
	}
}
