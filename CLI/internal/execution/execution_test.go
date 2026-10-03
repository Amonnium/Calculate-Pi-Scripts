package execution

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/catalog"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/platform"
)

func TestResolveUsesMatchingArtifact(t *testing.T) {
	root := t.TempDir()
	artifactPath := filepath.Join(root, "program.exe")
	if err := os.WriteFile(artifactPath, []byte("artifact"), 0o600); err != nil {
		t.Fatal(err)
	}

	implementation := catalog.Implementation{
		Name:            "Test",
		SourcePath:      "missing-source.c",
		Artifacts:       []catalog.Artifact{{OS: "windows", Arch: "amd64", Path: "program.exe"}},
		RuntimeCommands: []catalog.RuntimeCommand{{Program: "missing-runtime"}},
		RunSource:       true,
	}
	command, err := Resolve(implementation, platform.Target{OS: "windows", Arch: "amd64"}, root)
	if err != nil {
		t.Fatal(err)
	}
	if command.Program != artifactPath || command.Dir != root || len(command.Args) != 0 {
		t.Fatalf("Resolve() = %#v", command)
	}
}

func TestResolveReportsMissingRuntime(t *testing.T) {
	implementation := catalog.Implementation{
		Name:            "Test",
		SourcePath:      "script.py",
		RuntimeCommands: []catalog.RuntimeCommand{{Program: "calculate-pi-runtime-that-does-not-exist"}},
		RunSource:       true,
	}

	_, err := Resolve(implementation, platform.Detect(), t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "required runtime for Test was not found on PATH") {
		t.Fatalf("Resolve() error = %v, want useful missing-runtime error", err)
	}
}

func TestResolveUsesConfiguredRuntime(t *testing.T) {
	root := t.TempDir()
	sourcePath := filepath.Join(root, "script.py")
	if err := os.WriteFile(sourcePath, []byte("print('ok')"), 0o600); err != nil {
		t.Fatal(err)
	}

	implementation := catalog.Implementation{
		Name:            "Test",
		SourcePath:      "script.py",
		RuntimeCommands: []catalog.RuntimeCommand{{Program: os.Args[0], Args: []string{"-runtime-option"}}},
		RunSource:       true,
	}
	command, err := Resolve(implementation, platform.Target{OS: "linux", Arch: "amd64"}, root, "argument with spaces", "second")
	if err != nil {
		t.Fatal(err)
	}
	wantArgs := []string{"-runtime-option", sourcePath, "argument with spaces", "second"}
	if command.Program != os.Args[0] || command.Dir != root || !reflect.DeepEqual(command.Args, wantArgs) {
		t.Fatalf("Resolve() = %#v", command)
	}
}

func TestResolveDoesNotUseArtifactForDifferentPlatform(t *testing.T) {
	implementation, ok := catalog.Find("c")
	if !ok {
		t.Fatal("C implementation not found")
	}

	_, err := Resolve(implementation, platform.Target{OS: "linux", Arch: "amd64"}, t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "linux/amd64") {
		t.Fatalf("Resolve() error = %v, want unsupported linux/amd64 artifact", err)
	}
}

func TestRunForwardsArgumentsStreamsAndExitCode(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	command := Command{
		Program: executable,
		Args:    []string{"-test.run=TestRunnerChildProcess", "--", "argument with spaces", "second"},
		Env: []string{
			"CALCULATE_PI_RUNNER_HELPER=1",
			"CALCULATE_PI_RUNNER_EXIT_CODE=7",
		},
		Stdout: stdout,
		Stderr: stderr,
	}

	err = Run(context.Background(), command)
	var processError *ProcessError
	if !errors.As(err, &processError) {
		t.Fatalf("Run() error = %v, want ProcessError", err)
	}
	if processError.ExitCode() != 7 {
		t.Fatalf("ExitCode() = %d, want 7", processError.ExitCode())
	}
	if stdout.String() != "argument with spaces|second" {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if stderr.String() != "child stderr" {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunnerChildProcess(t *testing.T) {
	if os.Getenv("CALCULATE_PI_RUNNER_HELPER") != "1" {
		return
	}

	var arguments []string
	for index, argument := range os.Args {
		if argument == "--" {
			arguments = os.Args[index+1:]
			break
		}
	}
	fmt.Fprint(os.Stdout, strings.Join(arguments, "|"))
	fmt.Fprint(os.Stderr, "child stderr")
	exitCode, _ := strconv.Atoi(os.Getenv("CALCULATE_PI_RUNNER_EXIT_CODE"))
	os.Exit(exitCode)
}

func TestFindProjectRootWalksUpToRepositoryMarker(t *testing.T) {
	root := t.TempDir()
	markerDirectory := filepath.Join(root, "Pi_in_Python")
	if err := os.Mkdir(markerDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(markerDirectory, "Chudnovsky.py"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	start := filepath.Join(root, "CLI", "internal")
	if err := os.MkdirAll(start, 0o700); err != nil {
		t.Fatal(err)
	}

	got, err := FindProjectRoot(start)
	if err != nil {
		t.Fatal(err)
	}
	if got != root {
		t.Fatalf("FindProjectRoot() = %q, want %q", got, root)
	}
}

func TestRunBundledJavaImplementation(t *testing.T) {
	if _, err := exec.LookPath("java"); err != nil {
		t.Skip("Java runtime is not available")
	}
	root, err := FindProjectRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	implementation, _ := catalog.Find("java")
	command, err := Resolve(implementation, platform.Detect(), root)
	if err != nil {
		t.Fatal(err)
	}
	stdout := &bytes.Buffer{}
	command.Stdout = stdout
	command.Stderr = &bytes.Buffer{}
	if err := Run(context.Background(), command); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Calculating Pi with 100 decimal places") || !strings.Contains(stdout.String(), "3.141592653589793") {
		t.Fatalf("unexpected Java output: %q", stdout.String())
	}
}

func TestRunBundledGoExecutable(t *testing.T) {
	target := platform.Detect()
	if target.OS != "windows" || target.Arch != "amd64" {
		t.Skip("checked-in Go executable targets Windows amd64")
	}
	root, err := FindProjectRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	implementation, _ := catalog.Find("go")
	command, err := Resolve(implementation, target, root)
	if err != nil {
		t.Fatal(err)
	}
	stdout := &bytes.Buffer{}
	command.Stdout = stdout
	command.Stderr = &bytes.Buffer{}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := Run(ctx, command); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Leibniz Pi (10000000 iterations)") {
		t.Fatalf("unexpected Go output: %q", stdout.String())
	}
}

func TestRunBundledNativeExecutables(t *testing.T) {
	target := platform.Detect()
	if target.OS != "windows" || target.Arch != "amd64" {
		t.Skip("checked-in C, C++, and Rust executables target Windows amd64")
	}
	root, err := FindProjectRoot(".")
	if err != nil {
		t.Fatal(err)
	}

	for _, id := range []string{"c", "cpp", "rust"} {
		t.Run(id, func(t *testing.T) {
			implementation, ok := catalog.Find(id)
			if !ok {
				t.Fatalf("implementation %q not found", id)
			}
			command, err := Resolve(implementation, target, root)
			if err != nil {
				t.Fatal(err)
			}
			stdout := &bytes.Buffer{}
			command.Stdout = stdout
			command.Stderr = &bytes.Buffer{}
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			if err := Run(ctx, command); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(stdout.String(), "Pi approximation:") {
				t.Fatalf("unexpected %s output: %q", implementation.Name, stdout.String())
			}
		})
	}
}

func TestRunPythonImplementationWithInput(t *testing.T) {
	target := platform.Detect()
	implementation, _ := catalog.Find("python")
	available := false
	for _, candidate := range implementation.RuntimeCommandsFor(target) {
		if _, err := exec.LookPath(candidate.Program); err == nil {
			available = true
			break
		}
	}
	if !available {
		t.Skip("Python runtime is not available")
	}
	root, err := FindProjectRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	command, err := Resolve(implementation, target, root)
	if err != nil {
		t.Fatal(err)
	}
	stdout := &bytes.Buffer{}
	command.Stdin = strings.NewReader("20\n")
	command.Stdout = stdout
	command.Stderr = &bytes.Buffer{}
	if err := Run(context.Background(), command); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "3.1415926535") {
		t.Fatalf("unexpected Python output: %q", stdout.String())
	}
}
