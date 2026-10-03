package app

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/execution"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/platform"
)

func TestRunStreamsAndCapturesCalculationOutput(t *testing.T) {
	target := platform.Detect()
	if target.OS != "windows" || target.Arch != "amd64" {
		t.Skip("checked-in Go artifact targets Windows amd64")
	}
	root, err := execution.FindProjectRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	service := &Service{target: target, cacheDir: root}
	var terminalOutput bytes.Buffer
	var terminalErrors bytes.Buffer
	result := service.Run(context.Background(), "go", false, &terminalOutput, &terminalErrors)
	if result.Err != nil {
		t.Fatal(result.Err)
	}
	want := "Leibniz Pi (10000000 iterations)"
	if !strings.Contains(terminalOutput.String(), want) {
		t.Fatalf("terminal output = %q, want %q", terminalOutput.String(), want)
	}
	if result.Stdout != terminalOutput.String() {
		t.Fatalf("captured output = %q, terminal output = %q", result.Stdout, terminalOutput.String())
	}
}
