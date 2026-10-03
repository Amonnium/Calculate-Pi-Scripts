package catalog

import "github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/platform"

const RepositoryRevision = "5f6133455011b19077611c890566ff1e17815870"

const repositoryRawURL = "https://raw.githubusercontent.com/Amonnium/Calculate-Pi-Scripts/" + RepositoryRevision + "/"

type Download struct {
	Path        string
	URL         string
	GitBlobSHA1 string
	Executable  bool
}

type Artifact struct {
	OS          string
	Arch        string
	Path        string
	URL         string
	GitBlobSHA1 string
	Executable  bool
}

type RuntimeCommand struct {
	OS      string
	Program string
	Args    []string
}

type PackageSource struct {
	OS           string
	Distribution string
	Manager      string
	PackageID    string
}

type Toolchain struct {
	ID             string
	Name           string
	MinimumVersion string
	VersionArgs    []string
	PackageSource  []PackageSource
}

type Implementation struct {
	ID              string
	Name            string
	Description     string
	Algorithm       string
	SourcePath      string
	Downloads       []Download
	Artifacts       []Artifact
	RuntimeCommands []RuntimeCommand
	RuntimeArgs     []string
	RuntimeEnv      []string
	RuntimeFiles    []string
	RunSource       bool
	ToolchainID     string
}

var implementations = []Implementation{
	{
		ID:          "python",
		Name:        "Python",
		Description: "Calculate Pi to a requested precision.",
		Algorithm:   "Chudnovsky",
		SourcePath:  "Pi_in_Python/Chudnovsky.py",
		Downloads: []Download{
			repositoryDownload("Pi_in_Python/Chudnovsky.py", "5a39eea1dd73e3ec23cb72d93a3ee7d1e794d65e", false),
		},
		RuntimeCommands: []RuntimeCommand{
			{OS: "windows", Program: "python"},
			{OS: "windows", Program: "py", Args: []string{"-3"}},
			{OS: "linux", Program: "python3"},
			{OS: "linux", Program: "python"},
			{OS: "darwin", Program: "python3"},
			{OS: "darwin", Program: "python"},
			{Program: "python"},
		},
		RuntimeEnv:  []string{"PYTHONIOENCODING=utf-8"},
		RunSource:   true,
		ToolchainID: "python",
	},
	{
		ID:          "c",
		Name:        "C",
		Description: "Approximate Pi with a native executable when available.",
		Algorithm:   "Leibniz series",
		SourcePath:  "Pi_in_C_Folder/Pi.c",
		Artifacts: []Artifact{{
			OS:          "windows",
			Arch:        "amd64",
			Path:        "Pi_in_C_Folder/Pi_in_C.exe",
			URL:         repositoryRawURL + "Pi_in_C_Folder/Pi_in_C.exe",
			GitBlobSHA1: "67ba0e5e3c05fbcaf9d02e1146d40b00feb585a6",
			Executable:  true,
		}},
	},
	{
		ID:          "cpp",
		Name:        "C++",
		Description: "Approximate Pi with a native executable when available.",
		Algorithm:   "Leibniz series",
		SourcePath:  "Pi_in_C++_Folder/Pi.cpp",
		Artifacts: []Artifact{{
			OS:          "windows",
			Arch:        "amd64",
			Path:        "Pi_in_C++_Folder/Pi_in_C++.exe",
			URL:         repositoryRawURL + "Pi_in_C++_Folder/Pi_in_C++.exe",
			GitBlobSHA1: "bffe3674e6e669098e786e7728ac44307c83cb9b",
			Executable:  true,
		}},
	},
	{
		ID:          "rust",
		Name:        "Rust",
		Description: "Approximate Pi with a native executable when available.",
		Algorithm:   "Leibniz series",
		SourcePath:  "Pi_in_Rust_Folder/Pi.rs",
		Artifacts: []Artifact{{
			OS:          "windows",
			Arch:        "amd64",
			Path:        "Pi_in_Rust_Folder/Pi_in_Rust.exe",
			URL:         repositoryRawURL + "Pi_in_Rust_Folder/Pi_in_Rust.exe",
			GitBlobSHA1: "4b1208b5f9468f3e6880d95ab960bbb8f5b7a3da",
			Executable:  true,
		}},
	},
	{
		ID:          "go",
		Name:        "Go",
		Description: "Approximate Pi with a native executable when available.",
		Algorithm:   "Leibniz series",
		SourcePath:  "Pi_in_Go_Folder/Pi.go",
		Artifacts: []Artifact{{
			OS:          "windows",
			Arch:        "amd64",
			Path:        "Pi_in_Go_Folder/Pi_in_Go.exe",
			URL:         repositoryRawURL + "Pi_in_Go_Folder/Pi_in_Go.exe",
			GitBlobSHA1: "397a82f5c4486a59efbebd518afe5fdd8b4b89bb",
			Executable:  true,
		}},
	},
	{
		ID:          "java",
		Name:        "Java",
		Description: "Calculate Pi with arbitrary-precision decimal arithmetic.",
		Algorithm:   "Machin's formula",
		SourcePath:  "Pi_in_Java/Pi.java",
		Downloads: []Download{
			repositoryDownload("Pi_in_Java/Pi.class", "5095fb792732dfe00b87e3d13f07a148978773f8", false),
		},
		RuntimeCommands: []RuntimeCommand{{Program: "java"}},
		RuntimeArgs:     []string{"-cp", "Pi_in_Java", "Pi"},
		RuntimeFiles:    []string{"Pi_in_Java/Pi.class"},
		ToolchainID:     "java",
	},
	{
		ID:          "julia",
		Name:        "Julia",
		Description: "Calculate Pi to a configurable high precision.",
		Algorithm:   "Chudnovsky with binary splitting",
		SourcePath:  "Pi_in_Julia/Pi.jl",
		Downloads: []Download{
			repositoryDownload("Pi_in_Julia/Pi.jl", "36faf71a29005e0e58a08858033e4a1cbf4e61a0", false),
		},
		RuntimeCommands: []RuntimeCommand{{Program: "julia"}},
		RunSource:       true,
		ToolchainID:     "julia",
	},
}

var toolchains = []Toolchain{
	{
		ID:             "python",
		Name:           "Python 3",
		MinimumVersion: "3.8",
		VersionArgs:    []string{"--version"},
		PackageSource: []PackageSource{
			{OS: "windows", Manager: "winget", PackageID: "Python.Python.3.13"},
			{OS: "linux", Distribution: "debian", Manager: "apt-get", PackageID: "python3"},
			{OS: "linux", Distribution: "ubuntu", Manager: "apt-get", PackageID: "python3"},
		},
	},
	{
		ID:             "java",
		Name:           "Java runtime",
		MinimumVersion: "21",
		VersionArgs:    []string{"-version"},
		PackageSource: []PackageSource{
			{OS: "windows", Manager: "winget", PackageID: "Microsoft.OpenJDK.21"},
			{OS: "linux", Distribution: "debian", Manager: "apt-get", PackageID: "openjdk-21-jre"},
			{OS: "linux", Distribution: "ubuntu", Manager: "apt-get", PackageID: "openjdk-21-jre"},
		},
	},
	{
		ID:             "julia",
		Name:           "Julia",
		MinimumVersion: "1.0",
		VersionArgs:    []string{"--version"},
		PackageSource: []PackageSource{
			{OS: "windows", Manager: "winget", PackageID: "JuliaLang.Julia"},
			{OS: "linux", Distribution: "debian", Manager: "apt-get", PackageID: "julia"},
			{OS: "linux", Distribution: "ubuntu", Manager: "apt-get", PackageID: "julia"},
		},
	},
}

func All() []Implementation {
	result := make([]Implementation, len(implementations))
	copy(result, implementations)
	for index := range result {
		result[index].Artifacts = append([]Artifact(nil), result[index].Artifacts...)
		result[index].Downloads = append([]Download(nil), result[index].Downloads...)
		result[index].RuntimeArgs = append([]string(nil), result[index].RuntimeArgs...)
		result[index].RuntimeEnv = append([]string(nil), result[index].RuntimeEnv...)
		result[index].RuntimeFiles = append([]string(nil), result[index].RuntimeFiles...)
		result[index].RuntimeCommands = append([]RuntimeCommand(nil), result[index].RuntimeCommands...)
		for commandIndex := range result[index].RuntimeCommands {
			result[index].RuntimeCommands[commandIndex].Args = append([]string(nil), result[index].RuntimeCommands[commandIndex].Args...)
		}
	}
	return result
}

func Toolchains() []Toolchain {
	result := make([]Toolchain, len(toolchains))
	for index, toolchain := range toolchains {
		result[index] = toolchain
		result[index].VersionArgs = append([]string(nil), toolchain.VersionArgs...)
		result[index].PackageSource = append([]PackageSource(nil), toolchain.PackageSource...)
	}
	return result
}

func FindToolchain(id string) (Toolchain, bool) {
	for _, toolchain := range toolchains {
		if toolchain.ID == id {
			toolchain.VersionArgs = append([]string(nil), toolchain.VersionArgs...)
			toolchain.PackageSource = append([]PackageSource(nil), toolchain.PackageSource...)
			return toolchain, true
		}
	}
	return Toolchain{}, false
}

func (implementation Implementation) ManagedToolchain() (Toolchain, bool) {
	if implementation.ToolchainID == "" {
		return Toolchain{}, false
	}
	return FindToolchain(implementation.ToolchainID)
}

func Find(id string) (Implementation, bool) {
	for _, implementation := range implementations {
		if implementation.ID == id {
			return implementation, true
		}
	}
	return Implementation{}, false
}

func (implementation Implementation) ArtifactFor(target platform.Target) (Artifact, bool) {
	for _, artifact := range implementation.Artifacts {
		if artifact.OS == target.OS && artifact.Arch == target.Arch {
			return artifact, true
		}
	}
	return Artifact{}, false
}

func (implementation Implementation) DownloadsFor(target platform.Target) []Download {
	downloads := append([]Download(nil), implementation.Downloads...)
	if artifact, ok := implementation.ArtifactFor(target); ok {
		downloads = append(downloads, Download{
			Path:        artifact.Path,
			URL:         artifact.URL,
			GitBlobSHA1: artifact.GitBlobSHA1,
			Executable:  artifact.Executable,
		})
	}
	return downloads
}

func repositoryDownload(path, gitBlobSHA1 string, executable bool) Download {
	return Download{
		Path:        path,
		URL:         repositoryRawURL + path,
		GitBlobSHA1: gitBlobSHA1,
		Executable:  executable,
	}
}

func (implementation Implementation) RuntimeCommandsFor(target platform.Target) []RuntimeCommand {
	var commands []RuntimeCommand
	for _, command := range implementation.RuntimeCommands {
		if command.OS == target.OS {
			commands = append(commands, command)
		}
	}
	for _, command := range implementation.RuntimeCommands {
		if command.OS == "" {
			commands = append(commands, command)
		}
	}
	return commands
}
