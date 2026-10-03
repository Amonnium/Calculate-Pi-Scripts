package catalog

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/platform"
)

func TestAllContainsOnlySupportedImplementations(t *testing.T) {
	want := []string{"python", "c", "cpp", "rust", "go", "java", "julia"}
	items := All()
	got := make([]string, 0, len(items))
	for _, item := range items {
		got = append(got, item.ID)
		if item.SourcePath == "" {
			t.Errorf("implementation %q has no source path", item.ID)
		}
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("catalog IDs = %v, want %v", got, want)
	}
}

func TestArtifactForRequiresMatchingOSAndArchitecture(t *testing.T) {
	implementation, ok := Find("c")
	if !ok {
		t.Fatal("C implementation not found")
	}

	artifact, ok := implementation.ArtifactFor(platform.Target{OS: "windows", Arch: "amd64"})
	if !ok || artifact.Path != "Pi_in_C_Folder/Pi_in_C.exe" {
		t.Fatalf("Windows amd64 artifact = %#v, %t", artifact, ok)
	}
	if _, ok := implementation.ArtifactFor(platform.Target{OS: "linux", Arch: "amd64"}); ok {
		t.Fatal("Windows artifact matched a Linux target")
	}
}

func TestDownloadsForSelectsOnlyCurrentPlatformArtifacts(t *testing.T) {
	cImplementation, ok := Find("c")
	if !ok {
		t.Fatal("C implementation not found")
	}

	windowsDownloads := cImplementation.DownloadsFor(platform.Target{OS: "windows", Arch: "amd64"})
	if len(windowsDownloads) != 1 {
		t.Fatalf("Windows downloads = %d, want one executable", len(windowsDownloads))
	}
	if windowsDownloads[0].Path != "Pi_in_C_Folder/Pi_in_C.exe" || !windowsDownloads[0].Executable {
		t.Fatalf("Windows download = %#v", windowsDownloads[0])
	}
	if !strings.HasPrefix(windowsDownloads[0].URL, "https://raw.githubusercontent.com/") || windowsDownloads[0].GitBlobSHA1 == "" {
		t.Fatalf("download lacks HTTPS source or integrity hash: %#v", windowsDownloads[0])
	}

	linuxDownloads := cImplementation.DownloadsFor(platform.Target{OS: "linux", Arch: "amd64"})
	if len(linuxDownloads) != 0 {
		t.Fatalf("Linux downloads = %#v, want no Windows executable", linuxDownloads)
	}
}

func TestScriptImplementationHasPinnedHTTPSDownload(t *testing.T) {
	implementation, ok := Find("python")
	if !ok {
		t.Fatal("Python implementation not found")
	}

	downloads := implementation.DownloadsFor(platform.Target{OS: "linux", Arch: "amd64"})
	if len(downloads) != 1 {
		t.Fatalf("Python downloads = %#v, want its source script", downloads)
	}
	download := downloads[0]
	if download.Path != implementation.SourcePath || !strings.HasPrefix(download.URL, "https://") || download.GitBlobSHA1 == "" {
		t.Fatalf("Python download metadata = %#v", download)
	}
}

func TestManagedToolchainsContainOnlyRequiredRuntimes(t *testing.T) {
	toolchains := Toolchains()
	want := []string{"python", "java", "julia"}
	if len(toolchains) != len(want) {
		t.Fatalf("managed toolchains = %d, want %d", len(toolchains), len(want))
	}
	for index, spec := range toolchains {
		if spec.ID != want[index] {
			t.Errorf("toolchain %d = %q, want %q", index, spec.ID, want[index])
		}
	}
}
