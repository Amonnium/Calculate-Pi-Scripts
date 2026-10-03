package artifact

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/catalog"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/platform"
)

func TestEnsureVerifiesAndReusesCachedDownload(t *testing.T) {
	contents := []byte("trusted artifact")
	var requestCount atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestCount.Add(1)
		writer.Write(contents)
	}))
	defer server.Close()

	cacheDir := t.TempDir()
	downloader := NewDownloaderWithClient(cacheDir, server.Client())
	file := catalog.Download{
		Path:        "Pi_in_C_Folder/Pi_in_C.exe",
		URL:         server.URL + "/artifact",
		GitBlobSHA1: gitBlobSHA1(contents),
		Executable:  true,
	}

	path, err := downloader.Ensure(context.Background(), file)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(contents) {
		t.Fatalf("downloaded contents = %q", got)
	}
	if _, err := downloader.Ensure(context.Background(), file); err != nil {
		t.Fatal(err)
	}
	if requestCount.Load() != 1 {
		t.Fatalf("request count after cache hit = %d, want 1", requestCount.Load())
	}

	if err := os.WriteFile(path, []byte("corrupted cache"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := downloader.Ensure(context.Background(), file); err != nil {
		t.Fatal(err)
	}
	got, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(contents) || requestCount.Load() != 2 {
		t.Fatalf("cache recovery contents = %q, requests = %d", got, requestCount.Load())
	}
}

func TestEnsureFollowsHTTPSRedirectAndRejectsDowngrade(t *testing.T) {
	plainServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("insecure"))
	}))
	defer plainServer.Close()

	contents := []byte("redirected securely")
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/redirect":
			http.Redirect(writer, request, "/artifact", http.StatusFound)
		case "/downgrade":
			http.Redirect(writer, request, plainServer.URL, http.StatusFound)
		default:
			writer.Write(contents)
		}
	}))
	defer server.Close()

	downloader := NewDownloaderWithClient(t.TempDir(), server.Client())
	file := catalog.Download{Path: "redirected.bin", URL: server.URL + "/redirect", GitBlobSHA1: gitBlobSHA1(contents)}
	if _, err := downloader.Ensure(context.Background(), file); err != nil {
		t.Fatalf("HTTPS redirect failed: %v", err)
	}

	file.Path = "downgraded.bin"
	file.URL = server.URL + "/downgrade"
	if _, err := downloader.Ensure(context.Background(), file); err == nil || !strings.Contains(err.Error(), "non-HTTPS redirect") {
		t.Fatalf("downgrade error = %v, want HTTPS redirect rejection", err)
	}
}

func TestEnsureReportsHTTPErrorAndDoesNotCacheResponse(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "missing", http.StatusNotFound)
	}))
	defer server.Close()

	cacheDir := t.TempDir()
	downloader := NewDownloaderWithClient(cacheDir, server.Client())
	file := catalog.Download{Path: "missing.bin", URL: server.URL + "/missing"}
	if _, err := downloader.Ensure(context.Background(), file); err == nil || !strings.Contains(err.Error(), "server returned 404 Not Found") {
		t.Fatalf("HTTP error = %v, want 404", err)
	}
	if _, err := os.Stat(filepath.Join(cacheDir, file.Path)); !os.IsNotExist(err) {
		t.Fatalf("failed response created a cache file: %v", err)
	}
}

func TestEnsureRejectsHTTPURLAndUnsafePath(t *testing.T) {
	downloader := NewDownloader(t.TempDir())
	if _, err := downloader.Ensure(context.Background(), catalog.Download{Path: "asset", URL: "http://example.test/asset"}); err == nil || !strings.Contains(err.Error(), "must use HTTPS") {
		t.Fatalf("HTTP URL error = %v", err)
	}
	if _, err := downloader.Ensure(context.Background(), catalog.Download{Path: "../outside", URL: "https://example.test/asset"}); err == nil || !strings.Contains(err.Error(), "escapes the cache directory") {
		t.Fatalf("unsafe path error = %v", err)
	}
}

func TestEnsureRejectsIntegrityMismatchWithoutFinalizing(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("unexpected bytes"))
	}))
	defer server.Close()

	cacheDir := t.TempDir()
	downloader := NewDownloaderWithClient(cacheDir, server.Client())
	file := catalog.Download{Path: "asset.bin", URL: server.URL + "/asset", GitBlobSHA1: gitBlobSHA1([]byte("expected bytes"))}
	if _, err := downloader.Ensure(context.Background(), file); err == nil || !strings.Contains(err.Error(), "SHA-1 mismatch") {
		t.Fatalf("integrity error = %v, want hash mismatch", err)
	}
	if _, err := os.Stat(filepath.Join(cacheDir, file.Path)); !os.IsNotExist(err) {
		t.Fatalf("invalid artifact was finalized: %v", err)
	}
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("temporary files remain after failed verification: %v", entries)
	}
}

func TestEnsureHonorsRequestTimeout(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		<-request.Context().Done()
	}))
	defer server.Close()

	client := server.Client()
	client.Timeout = 20 * time.Millisecond
	downloader := NewDownloaderWithClient(t.TempDir(), client)
	_, err := downloader.Ensure(context.Background(), catalog.Download{Path: "slow.bin", URL: server.URL + "/slow"})
	if err == nil || !strings.Contains(err.Error(), "Client.Timeout") {
		t.Fatalf("timeout error = %v, want client timeout", err)
	}
}

func TestCacheDirectoryIsPlatformAndRevisionSpecific(t *testing.T) {
	base, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	target := platform.Target{OS: "windows", Arch: "amd64"}
	got, err := CacheDirectory(target, "revision123")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(base, "CalculatePiScripts", "windows-amd64", "revision123")
	if got != want {
		t.Fatalf("CacheDirectory() = %q, want %q", got, want)
	}
	if _, err := CacheDirectory(target, "../outside"); err == nil {
		t.Fatal("invalid revision was accepted")
	}
}

func gitBlobSHA1(contents []byte) string {
	hasher := sha1.New()
	fmt.Fprintf(hasher, "blob %d\x00", len(contents))
	io.WriteString(hasher, string(contents))
	return hex.EncodeToString(hasher.Sum(nil))
}
