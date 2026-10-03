package toolchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/catalog"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/platform"
)

var versionPattern = regexp.MustCompile(`[0-9]+(?:\.[0-9]+){0,3}(?:[-+][A-Za-z0-9.-]+)?`)

type Status struct {
	Installed  bool
	Compatible bool
	Managed    bool
	Program    string
	Version    string
}

type processRunner interface {
	Capture(context.Context, string, ...string) ([]byte, error)
	Interactive(context.Context, string, ...string) error
}

type systemRunner struct{}

func (systemRunner) Capture(ctx context.Context, program string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, program, args...).CombinedOutput()
}

func (systemRunner) Interactive(ctx context.Context, program string, args ...string) error {
	command := exec.CommandContext(ctx, program, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

type ownership struct {
	PackageManager string `json:"packageManager"`
	PackageID      string `json:"packageId"`
}

type Manager struct {
	target   platform.Target
	distro   string
	lookPath func(string) (string, error)
	runner   processRunner
	config   string
}

func NewManager(target platform.Target) (*Manager, error) {
	configRoot, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("locate user configuration directory: %w", err)
	}
	manager := &Manager{
		target:   target,
		lookPath: exec.LookPath,
		runner:   systemRunner{},
		config:   filepath.Join(configRoot, "CalculatePiScripts", "toolchains.json"),
	}
	if target.OS == "linux" {
		manager.distro, err = linuxDistribution()
		if err != nil {
			return nil, err
		}
	}
	return manager, nil
}

func (manager *Manager) Detect(ctx context.Context, spec catalog.Toolchain, implementation catalog.Implementation) Status {
	status := Status{Managed: manager.isManaged(spec.ID)}
	var incompatible Status
	for _, candidate := range implementation.RuntimeCommandsFor(manager.target) {
		program, err := manager.lookPath(candidate.Program)
		if err != nil {
			continue
		}
		args := append([]string(nil), candidate.Args...)
		args = append(args, spec.VersionArgs...)
		output, runErr := manager.runner.Capture(ctx, program, args...)
		if runErr != nil && len(bytesTrimSpace(output)) == 0 {
			continue
		}
		version := versionPattern.FindString(string(output))
		if version == "" {
			version = strings.TrimSpace(string(output))
		}
		if version == "" {
			version = "version unavailable"
		}
		candidateStatus := Status{
			Installed:  true,
			Compatible: versionAtLeast(version, spec.MinimumVersion),
			Managed:    status.Managed,
			Program:    program,
			Version:    version,
		}
		if candidateStatus.Compatible {
			return candidateStatus
		}
		if !incompatible.Installed {
			incompatible = candidateStatus
		}
	}
	if incompatible.Installed {
		return incompatible
	}
	return status
}

func (manager *Manager) Install(ctx context.Context, spec catalog.Toolchain) (Status, error) {
	if implementation := implementationForToolchain(spec.ID); implementation.ID != "" {
		if status := manager.Detect(ctx, spec, implementation); status.Compatible {
			return status, nil
		}
	}
	if len(spec.PackageSource) == 0 {
		return Status{}, fmt.Errorf("no supported installation method is defined for %s", spec.Name)
	}
	for _, source := range spec.PackageSource {
		if source.OS != manager.target.OS || source.Distribution != "" && source.Distribution != manager.distro {
			continue
		}
		return manager.installWith(ctx, spec, source)
	}
	if manager.target.OS == "linux" {
		return Status{}, fmt.Errorf("automatic installation of %s is supported on Debian and Ubuntu Linux only; detected %q", spec.Name, manager.distro)
	}
	return Status{}, fmt.Errorf("no supported installation method for %s on %s", spec.Name, manager.target.OS)
}

func (manager *Manager) installWith(ctx context.Context, spec catalog.Toolchain, source catalog.PackageSource) (Status, error) {
	if _, err := manager.lookPath(source.Manager); err != nil {
		return Status{}, fmt.Errorf("%s is required to install %s but was not found on PATH", source.Manager, spec.Name)
	}
	program, prefix, err := manager.privilegedCommand(source.Manager)
	if err != nil {
		return Status{}, err
	}
	if source.Manager == "apt-get" {
		if err := manager.runner.Interactive(ctx, program, append(append([]string(nil), prefix...), "update")...); err != nil {
			return Status{}, fmt.Errorf("apt package index update failed: %w", err)
		}
	}
	args := installArgs(source)
	args = append(prefix, args...)
	if err := manager.runner.Interactive(ctx, program, args...); err != nil {
		return Status{}, fmt.Errorf("install %s with %s failed: %w", spec.Name, source.Manager, err)
	}
	if err := manager.updateOwnership(spec.ID, &ownership{PackageManager: source.Manager, PackageID: source.PackageID}); err != nil {
		return Status{}, fmt.Errorf("%s was installed, but its ownership could not be recorded: %w", spec.Name, err)
	}

	status := manager.Detect(ctx, spec, implementationForToolchain(spec.ID))
	status.Managed = true
	if !status.Compatible {
		return status, fmt.Errorf("%s installation completed, but its runtime is not available in this process; restart the CLI and check PATH", spec.Name)
	}
	return status, nil
}

func versionAtLeast(version, minimum string) bool {
	versionParts := versionPattern.FindAllString(version, 1)
	if len(versionParts) == 0 {
		return false
	}
	return compareVersions(versionParts[0], minimum) >= 0
}

func compareVersions(left, right string) int {
	leftParts := strings.Split(left, ".")
	rightParts := strings.Split(right, ".")
	partCount := len(leftParts)
	if len(rightParts) > partCount {
		partCount = len(rightParts)
	}
	for index := 0; index < partCount; index++ {
		leftValue, rightValue := 0, 0
		if index < len(leftParts) {
			leftValue, _ = strconv.Atoi(leftParts[index])
		}
		if index < len(rightParts) {
			rightValue, _ = strconv.Atoi(rightParts[index])
		}
		if leftValue < rightValue {
			return -1
		}
		if leftValue > rightValue {
			return 1
		}
	}
	return 0
}

func (manager *Manager) Uninstall(ctx context.Context, spec catalog.Toolchain) error {
	installedByCLI, err := manager.readOwnership()
	if err != nil {
		return err
	}
	owned, ok := installedByCLI[spec.ID]
	if !ok {
		return fmt.Errorf("refusing to uninstall %s: it was not recorded as installed by this CLI", spec.Name)
	}
	source, ok := manager.packageSource(spec, owned)
	if !ok {
		return fmt.Errorf("refusing to uninstall %s: its recorded package source is not supported for this platform", spec.Name)
	}
	if _, err := manager.lookPath(source.Manager); err != nil {
		return fmt.Errorf("%s is required to uninstall %s but was not found on PATH", source.Manager, spec.Name)
	}
	program, prefix, err := manager.privilegedCommand(source.Manager)
	if err != nil {
		return err
	}
	args := uninstallArgs(source)
	args = append(prefix, args...)
	if err := manager.runner.Interactive(ctx, program, args...); err != nil {
		return fmt.Errorf("uninstall %s with %s failed: %w", spec.Name, source.Manager, err)
	}
	delete(installedByCLI, spec.ID)
	return manager.writeOwnership(installedByCLI)
}

func (manager *Manager) packageSource(spec catalog.Toolchain, owned ownership) (catalog.PackageSource, bool) {
	for _, source := range spec.PackageSource {
		if source.OS == manager.target.OS && source.Distribution == manager.distro && source.Manager == owned.PackageManager && source.PackageID == owned.PackageID {
			return source, true
		}
	}
	return catalog.PackageSource{}, false
}

func (manager *Manager) privilegedCommand(packageManager string) (string, []string, error) {
	if manager.target.OS == "windows" {
		return packageManager, nil, nil
	}
	if os.Geteuid() == 0 {
		return packageManager, nil, nil
	}
	if _, err := manager.lookPath("sudo"); err != nil {
		return "", nil, errors.New("package installation requires root privileges or sudo, which was not found on PATH")
	}
	return "sudo", []string{packageManager + "-get"}, nil
}

func installArgs(source catalog.PackageSource) []string {
	switch source.Manager {
	case "winget":
		return []string{"install", "--id", source.PackageID, "--exact", "--accept-package-agreements", "--accept-source-agreements"}
	case "apt-get":
		return []string{"install", "--yes", source.PackageID}
	default:
		return nil
	}
}

func uninstallArgs(source catalog.PackageSource) []string {
	switch source.Manager {
	case "winget":
		return []string{"uninstall", "--id", source.PackageID, "--exact", "--accept-source-agreements"}
	case "apt-get":
		return []string{"remove", source.PackageID}
	default:
		return nil
	}
}

func (manager *Manager) isManaged(id string) bool {
	owned, err := manager.readOwnership()
	if err != nil {
		return false
	}
	_, ok := owned[id]
	return ok
}

func (manager *Manager) readOwnership() (map[string]ownership, error) {
	data, err := os.ReadFile(manager.config)
	if errors.Is(err, os.ErrNotExist) {
		return make(map[string]ownership), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read toolchain ownership: %w", err)
	}
	owned := make(map[string]ownership)
	if err := json.Unmarshal(data, &owned); err != nil {
		return nil, fmt.Errorf("parse toolchain ownership: %w", err)
	}
	return owned, nil
}

func (manager *Manager) updateOwnership(id string, item *ownership) error {
	owned, err := manager.readOwnership()
	if err != nil {
		return err
	}
	if item == nil {
		delete(owned, id)
	} else {
		owned[id] = *item
	}
	return manager.writeOwnership(owned)
}

func (manager *Manager) writeOwnership(owned map[string]ownership) error {
	if err := os.MkdirAll(filepath.Dir(manager.config), 0o700); err != nil {
		return fmt.Errorf("create toolchain state directory: %w", err)
	}
	data, err := json.MarshalIndent(owned, "", "  ")
	if err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(manager.config), ".toolchains-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := replaceStateFile(temporaryPath, manager.config); err != nil {
		return err
	}
	return nil
}

func replaceStateFile(source, destination string) error {
	if _, err := os.Lstat(destination); errors.Is(err, os.ErrNotExist) {
		return os.Rename(source, destination)
	} else if err != nil {
		return err
	}
	backup, err := os.CreateTemp(filepath.Dir(destination), ".toolchains-previous-*")
	if err != nil {
		return err
	}
	backupPath := backup.Name()
	if err := backup.Close(); err != nil {
		os.Remove(backupPath)
		return err
	}
	if err := os.Remove(backupPath); err != nil {
		return err
	}
	if err := os.Rename(destination, backupPath); err != nil {
		return err
	}
	if err := os.Rename(source, destination); err != nil {
		return errors.Join(err, os.Rename(backupPath, destination))
	}
	return os.Remove(backupPath)
}

func linuxDistribution() (string, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "", fmt.Errorf("read /etc/os-release to select a supported Linux package source: %w", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("read Linux distribution information: %w", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "ID=") {
			return strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "ID=")), `"'`), nil
		}
	}
	return "", errors.New("Linux distribution ID is missing from /etc/os-release")
}

func implementationForToolchain(id string) catalog.Implementation {
	for _, implementation := range catalog.All() {
		if implementation.ToolchainID == id {
			return implementation
		}
	}
	return catalog.Implementation{}
}

func bytesTrimSpace(value []byte) []byte {
	return []byte(strings.TrimSpace(string(value)))
}
