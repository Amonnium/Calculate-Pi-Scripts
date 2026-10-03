package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/artifact"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/catalog"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/execution"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/platform"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/toolchain"
)

type Service struct {
	target   platform.Target
	cacheDir string
	args     []string
}

type RuntimeRequirement struct {
	Toolchain catalog.Toolchain
	Status    toolchain.Status
}

type ToolchainStatus struct {
	Toolchain catalog.Toolchain
	Status    toolchain.Status
}

type Result struct {
	Message             string
	Stdout              string
	Stderr              string
	ExitCode            int
	NeedsRuntimeInstall bool
	RuntimePrompt       string
	RuntimeID           string
	Err                 error
}

func New(args []string) (*Service, error) {
	target := platform.Detect()
	cacheDir, err := artifact.CacheDirectory(target, catalog.RepositoryRevision)
	if err != nil {
		return nil, err
	}
	return &Service{target: target, cacheDir: cacheDir, args: append([]string(nil), args...)}, nil
}

func (service *Service) Implementations() []catalog.Implementation {
	return catalog.All()
}

func (service *Service) Toolchains(ctx context.Context) ([]ToolchainStatus, error) {
	manager, err := toolchain.NewManager(service.target)
	if err != nil {
		return nil, err
	}
	statuses := make([]ToolchainStatus, 0)
	for _, spec := range catalog.Toolchains() {
		implementation, ok := catalog.Find(spec.ID)
		if !ok {
			continue
		}
		statuses = append(statuses, ToolchainStatus{
			Toolchain: spec,
			Status:    manager.Detect(ctx, spec, implementation),
		})
	}
	return statuses, nil
}

func (service *Service) CheckRuntime(ctx context.Context, implementationID string) (RuntimeRequirement, error) {
	implementation, ok := catalog.Find(implementationID)
	if !ok {
		return RuntimeRequirement{}, fmt.Errorf("unknown implementation %q", implementationID)
	}
	spec, ok := implementation.ManagedToolchain()
	if !ok {
		return RuntimeRequirement{}, nil
	}
	manager, err := toolchain.NewManager(service.target)
	if err != nil {
		return RuntimeRequirement{}, err
	}
	return RuntimeRequirement{Toolchain: spec, Status: manager.Detect(ctx, spec, implementation)}, nil
}

func (service *Service) Prepare(ctx context.Context, implementationID string) (RuntimeRequirement, error) {
	implementation, ok := catalog.Find(implementationID)
	if !ok {
		return RuntimeRequirement{}, fmt.Errorf("unknown implementation %q", implementationID)
	}
	requirement, err := service.CheckRuntime(ctx, implementationID)
	if err != nil {
		return RuntimeRequirement{}, err
	}
	downloader := artifact.NewDownloader(service.cacheDir)
	for _, file := range implementation.DownloadsFor(service.target) {
		if _, err := downloader.Ensure(ctx, file); err != nil {
			return RuntimeRequirement{}, fmt.Errorf("obtain %s artifact: %w", implementation.Name, err)
		}
	}
	return requirement, nil
}

func (service *Service) Run(ctx context.Context, implementationID string, installMissingRuntime bool, terminalStdout, terminalStderr io.Writer) Result {
	implementation, ok := catalog.Find(implementationID)
	if !ok {
		return Result{Message: "Unable to run implementation", Err: fmt.Errorf("unknown implementation %q", implementationID)}
	}

	runtimeRequirement, err := service.CheckRuntime(ctx, implementationID)
	if err != nil {
		return Result{Message: "Unable to check runtime requirements", Err: err}
	}

	downloader := artifact.NewDownloader(service.cacheDir)
	for _, file := range implementation.DownloadsFor(service.target) {
		if _, err := downloader.Ensure(ctx, file); err != nil {
			return Result{Message: "Unable to obtain implementation artifact", Err: err}
		}
	}

	if runtimeRequirement.Toolchain.ID != "" && !runtimeRequirement.Status.Compatible {
		if !installMissingRuntime {
			return Result{
				Message:             "A runtime is required",
				NeedsRuntimeInstall: true,
				RuntimePrompt:       fmt.Sprintf("Install %s %s using the system package manager?", runtimeRequirement.Toolchain.Name, runtimeRequirement.Toolchain.MinimumVersion),
				RuntimeID:           runtimeRequirement.Toolchain.ID,
			}
		}
		installResult := service.InstallToolchain(ctx, runtimeRequirement.Toolchain.ID)
		if installResult.Err != nil {
			return installResult
		}
	}

	command, err := execution.Resolve(implementation, service.target, service.cacheDir, service.args...)
	if err != nil {
		return Result{Message: "Unable to prepare implementation", Err: err}
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	command.Stdout = outputWriter(stdout, terminalStdout)
	command.Stderr = outputWriter(stderr, terminalStderr)
	err = execution.Run(ctx, command)
	result := Result{Stdout: stdout.String(), Stderr: stderr.String(), ExitCode: 0}
	if err != nil {
		result.Err = err
		result.Message = "Implementation failed"
		var processError *execution.ProcessError
		if errors.As(err, &processError) {
			result.ExitCode = processError.ExitCode()
		}
		return result
	}
	result.Message = "Implementation completed successfully"
	return result
}

func outputWriter(capture *bytes.Buffer, terminal io.Writer) io.Writer {
	if terminal == nil {
		return capture
	}
	return io.MultiWriter(capture, terminal)
}

func (service *Service) InstallToolchain(ctx context.Context, id string) Result {
	spec, ok := catalog.FindToolchain(id)
	if !ok {
		return Result{Message: "Unable to install toolchain", Err: fmt.Errorf("toolchain %q is not managed by this CLI", id)}
	}
	manager, err := toolchain.NewManager(service.target)
	if err != nil {
		return Result{Message: "Unable to prepare toolchain installation", Err: err}
	}
	status, err := manager.Install(ctx, spec)
	if err != nil {
		return Result{Message: "Toolchain installation failed", Err: err}
	}
	return Result{Message: fmt.Sprintf("Installed %s %s at %s", spec.Name, status.Version, status.Program)}
}

func (service *Service) UninstallToolchain(ctx context.Context, id string) Result {
	spec, ok := catalog.FindToolchain(id)
	if !ok {
		return Result{Message: "Unable to uninstall toolchain", Err: fmt.Errorf("toolchain %q is not managed by this CLI", id)}
	}
	manager, err := toolchain.NewManager(service.target)
	if err != nil {
		return Result{Message: "Unable to prepare toolchain removal", Err: err}
	}
	if err := manager.Uninstall(ctx, spec); err != nil {
		return Result{Message: "Toolchain removal failed", Err: err}
	}
	return Result{Message: fmt.Sprintf("Uninstalled %s", spec.Name)}
}
