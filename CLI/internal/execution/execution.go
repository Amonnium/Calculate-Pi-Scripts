package execution

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/catalog"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/platform"
)

type Command struct {
	Program string
	Args    []string
	Dir     string
	Env     []string
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
}

type ProcessError struct {
	Program string
	Err     error
}

func (processError *ProcessError) Error() string {
	var exitError *exec.ExitError
	if errors.As(processError.Err, &exitError) {
		return fmt.Sprintf("%q exited with code %d", processError.Program, exitError.ExitCode())
	}
	return fmt.Sprintf("could not run %q: %v", processError.Program, processError.Err)
}

func (processError *ProcessError) Unwrap() error {
	return processError.Err
}

func (processError *ProcessError) ExitCode() int {
	var exitError *exec.ExitError
	if errors.As(processError.Err, &exitError) && exitError.ExitCode() >= 0 {
		return exitError.ExitCode()
	}
	return 1
}

func FindProjectRoot(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("resolve current directory: %w", err)
	}

	for {
		marker := filepath.Join(current, "Pi_in_Python", "Chudnovsky.py")
		if _, err := os.Stat(marker); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("could not find the Calculate Pi Scripts repository from %q", start)
		}
		current = parent
	}
}

func Resolve(implementation catalog.Implementation, target platform.Target, projectRoot string, implementationArgs ...string) (Command, error) {
	if artifact, ok := implementation.ArtifactFor(target); ok {
		program := filepath.Join(projectRoot, filepath.FromSlash(artifact.Path))
		if _, err := os.Stat(program); err != nil {
			return Command{}, fmt.Errorf("artifact for %s is unavailable at %q: %w", implementation.Name, program, err)
		}
		return Command{Program: program, Args: append([]string(nil), implementationArgs...), Dir: projectRoot}, nil
	}

	for _, runtimeCommand := range implementation.RuntimeCommandsFor(target) {
		program, err := exec.LookPath(runtimeCommand.Program)
		if err != nil {
			continue
		}
		for _, requiredFile := range implementation.RuntimeFiles {
			path := filepath.Join(projectRoot, filepath.FromSlash(requiredFile))
			if _, err := os.Stat(path); err != nil {
				return Command{}, fmt.Errorf("required artifact for %s is unavailable at %q: %w", implementation.Name, path, err)
			}
		}

		args := append([]string(nil), runtimeCommand.Args...)
		args = append(args, implementation.RuntimeArgs...)
		if implementation.RunSource {
			source := filepath.Join(projectRoot, filepath.FromSlash(implementation.SourcePath))
			if _, err := os.Stat(source); err != nil {
				return Command{}, fmt.Errorf("source for %s is unavailable at %q: %w", implementation.Name, source, err)
			}
			args = append(args, source)
		}
		args = append(args, implementationArgs...)
		return Command{
			Program: program,
			Args:    args,
			Dir:     projectRoot,
			Env:     append([]string(nil), implementation.RuntimeEnv...),
		}, nil
	}
	if len(implementation.RuntimeCommands) > 0 {
		candidates := implementation.RuntimeCommandsFor(target)
		programs := make([]string, 0, len(candidates))
		for _, candidate := range candidates {
			programs = append(programs, candidate.Program)
		}
		if len(programs) == 0 {
			return Command{}, fmt.Errorf("no runtime configured for %s on %s/%s", implementation.Name, target.OS, target.Arch)
		}
		return Command{}, fmt.Errorf("required runtime for %s was not found on PATH (checked: %s)", implementation.Name, strings.Join(programs, ", "))
	}

	return Command{}, fmt.Errorf(
		"no runnable artifact for %s on %s/%s; source is available at %s",
		implementation.Name,
		target.OS,
		target.Arch,
		implementation.SourcePath,
	)
}

func Run(ctx context.Context, command Command) error {
	process := exec.CommandContext(ctx, command.Program, command.Args...)
	process.Dir = command.Dir
	process.Env = append(os.Environ(), command.Env...)
	process.Stdin = command.Stdin
	if process.Stdin == nil {
		process.Stdin = os.Stdin
	}
	process.Stdout = command.Stdout
	if process.Stdout == nil {
		process.Stdout = os.Stdout
	}
	process.Stderr = command.Stderr
	if process.Stderr == nil {
		process.Stderr = os.Stderr
	}
	if err := process.Run(); err != nil {
		return &ProcessError{Program: command.Program, Err: err}
	}
	return nil
}
