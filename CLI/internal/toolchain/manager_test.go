package toolchain

import (
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/catalog"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/platform"
)

type fakeRunner struct {
	installed      bool
	version        string
	interactiveErr error
	commands       [][]string
}

func (runner *fakeRunner) Capture(_ context.Context, _ string, _ ...string) ([]byte, error) {
	if !runner.installed {
		return nil, errors.New("runtime missing")
	}
	return []byte(runner.version), nil
}

func (runner *fakeRunner) Interactive(_ context.Context, program string, args ...string) error {
	command := append([]string{program}, args...)
	runner.commands = append(runner.commands, command)
	if runner.interactiveErr != nil {
		return runner.interactiveErr
	}
	if len(args) > 0 && args[0] == "install" {
		runner.installed = true
	}
	if len(args) > 0 && args[0] == "uninstall" {
		runner.installed = false
	}
	return nil
}

func TestDetectReportsVersion(t *testing.T) {
	runner := &fakeRunner{installed: true, version: "Python 3.13.4\n"}
	manager := testManager(t, runner, true)
	implementation, _ := catalog.Find("python")
	toolchain, _ := catalog.FindToolchain("python")

	status := manager.Detect(context.Background(), toolchain, implementation)
	if !status.Installed || !status.Compatible || status.Version != "3.13.4" || status.Program != "python.exe" {
		t.Fatalf("Detect() = %#v", status)
	}
}

func TestDetectRejectsIncompatiblePythonVersion(t *testing.T) {
	runner := &fakeRunner{installed: true, version: "Python 2.7.18\n"}
	manager := testManager(t, runner, true)
	implementation, _ := catalog.Find("python")
	toolchain, _ := catalog.FindToolchain("python")

	status := manager.Detect(context.Background(), toolchain, implementation)
	if !status.Installed || status.Compatible || status.Version != "2.7.18" {
		t.Fatalf("Detect() = %#v, want installed but incompatible Python", status)
	}
}

func TestJavaMinimumVersionMatchesClassFileRequirement(t *testing.T) {
	toolchain, ok := catalog.FindToolchain("java")
	if !ok || toolchain.MinimumVersion != "21" {
		t.Fatalf("Java toolchain minimum = %#v", toolchain)
	}
	if !versionAtLeast("openjdk version \"21.0.2\"", toolchain.MinimumVersion) || versionAtLeast("openjdk version \"17.0.9\"", toolchain.MinimumVersion) {
		t.Fatal("Java version compatibility did not enforce the class-file minimum")
	}
}

func TestInstallRecordsOwnershipAndUninstallsOnlyManagedRuntime(t *testing.T) {
	runner := &fakeRunner{version: "Python 3.13.4\n"}
	manager := testManager(t, runner, false)
	toolchain, _ := catalog.FindToolchain("python")

	status, err := manager.Install(context.Background(), toolchain)
	if err != nil {
		t.Fatal(err)
	}
	if !status.Installed || !status.Managed || status.Version != "3.13.4" {
		t.Fatalf("Install() status = %#v", status)
	}
	wantInstall := []string{"winget", "install", "--id", "Python.Python.3.13", "--exact", "--accept-package-agreements", "--accept-source-agreements"}
	if !reflect.DeepEqual(runner.commands[0], wantInstall) {
		t.Fatalf("install command = %v, want %v", runner.commands[0], wantInstall)
	}

	if err := manager.Uninstall(context.Background(), toolchain); err != nil {
		t.Fatal(err)
	}
	wantUninstall := []string{"winget", "uninstall", "--id", "Python.Python.3.13", "--exact", "--accept-source-agreements"}
	if !reflect.DeepEqual(runner.commands[1], wantUninstall) {
		t.Fatalf("uninstall command = %v, want %v", runner.commands[1], wantUninstall)
	}
	owned, err := manager.readOwnership()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := owned[toolchain.ID]; ok {
		t.Fatalf("ownership record remains after uninstall: %#v", owned)
	}
}

func TestInstallDoesNotReinstallAvailableRuntime(t *testing.T) {
	runner := &fakeRunner{installed: true, version: "Python 3.13.4\n"}
	manager := testManager(t, runner, true)
	toolchain, _ := catalog.FindToolchain("python")

	status, err := manager.Install(context.Background(), toolchain)
	if err != nil {
		t.Fatal(err)
	}
	if !status.Installed || len(runner.commands) != 0 {
		t.Fatalf("Install() status = %#v, commands = %v", status, runner.commands)
	}
}

func TestInstallFailureDoesNotRecordOwnership(t *testing.T) {
	runner := &fakeRunner{interactiveErr: errors.New("package installation failed")}
	manager := testManager(t, runner, false)
	toolchain, _ := catalog.FindToolchain("python")

	if _, err := manager.Install(context.Background(), toolchain); err == nil || !strings.Contains(err.Error(), "install Python 3 with winget failed") {
		t.Fatalf("Install() error = %v, want installation failure", err)
	}
	if manager.isManaged(toolchain.ID) {
		t.Fatal("failed installation was recorded as CLI-managed")
	}
}

func TestUninstallRefusesUnmanagedRuntime(t *testing.T) {
	manager := testManager(t, &fakeRunner{}, true)
	toolchain, _ := catalog.FindToolchain("java")
	if err := manager.Uninstall(context.Background(), toolchain); err == nil || !strings.Contains(err.Error(), "not recorded as installed by this CLI") {
		t.Fatalf("Uninstall() error = %v, want ownership refusal", err)
	}
}

func TestDetectPreservesManagedStateWhenRuntimeIsMissing(t *testing.T) {
	manager := testManager(t, &fakeRunner{}, false)
	toolchain, _ := catalog.FindToolchain("python")
	if err := manager.updateOwnership(toolchain.ID, &ownership{PackageManager: "winget", PackageID: "Python.Python.3.13"}); err != nil {
		t.Fatal(err)
	}
	implementation, _ := catalog.Find("python")
	status := manager.Detect(context.Background(), toolchain, implementation)
	if status.Installed || !status.Managed {
		t.Fatalf("Detect() = %#v, want missing but CLI-managed runtime", status)
	}
}

func TestInstallReportsMissingPackageManager(t *testing.T) {
	runner := &fakeRunner{}
	manager := testManager(t, runner, false)
	manager.lookPath = func(program string) (string, error) {
		if program == "winget" {
			return "", errors.New("not found")
		}
		return program + ".exe", nil
	}
	toolchain, _ := catalog.FindToolchain("julia")
	if _, err := manager.Install(context.Background(), toolchain); err == nil || !strings.Contains(err.Error(), "winget is required") {
		t.Fatalf("Install() error = %v, want missing package manager", err)
	}
	if len(runner.commands) != 0 {
		t.Fatalf("package manager was invoked unexpectedly: %v", runner.commands)
	}
}

func TestInstallRejectsUnsupportedLinuxDistribution(t *testing.T) {
	manager := testManager(t, &fakeRunner{}, false)
	manager.target = platform.Target{OS: "linux", Arch: "amd64"}
	manager.distro = "arch"
	toolchain, _ := catalog.FindToolchain("java")
	if _, err := manager.Install(context.Background(), toolchain); err == nil || !strings.Contains(err.Error(), "Debian and Ubuntu Linux only") {
		t.Fatalf("Install() error = %v, want unsupported distribution", err)
	}
}

func TestPrivilegedCommandUsesAptGetDirectlyWhenRoot(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("root-specific package command")
	}
	manager := testManager(t, &fakeRunner{}, false)
	manager.target.OS = "linux"
	program, prefix, err := manager.privilegedCommand("apt-get")
	if err != nil {
		t.Fatal(err)
	}
	if program != "apt-get" || len(prefix) != 0 {
		t.Fatalf("privilegedCommand() = %q %v, want apt-get with no prefix", program, prefix)
	}
}

func testManager(t *testing.T, runner *fakeRunner, installed bool) *Manager {
	t.Helper()
	runner.installed = installed
	return &Manager{
		target: platform.Target{OS: "windows", Arch: "amd64"},
		lookPath: func(program string) (string, error) {
			if program == "winget" || runner.installed && (program == "python" || program == "py") {
				return program + ".exe", nil
			}
			return "", errors.New("not found")
		},
		runner: runner,
		config: t.TempDir() + "\\toolchains.json",
	}
}
