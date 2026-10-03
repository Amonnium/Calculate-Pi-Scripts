package artifact

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/catalog"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/platform"
)

const (
	requestTimeout = 45 * time.Second
	maximumSize    = 128 << 20
)

type Downloader struct {
	cacheDir string
	client   *http.Client
	mutex    sync.Mutex
	locks    map[string]*sync.Mutex
}

func NewDownloader(cacheDir string) *Downloader {
	return NewDownloaderWithClient(cacheDir, &http.Client{Timeout: requestTimeout})
}

func NewDownloaderWithClient(cacheDir string, client *http.Client) *Downloader {
	if client == nil {
		client = &http.Client{}
	}
	configuredClient := *client
	if configuredClient.Timeout <= 0 {
		configuredClient.Timeout = requestTimeout
	}
	previousRedirectCheck := configuredClient.CheckRedirect
	configuredClient.CheckRedirect = func(request *http.Request, previous []*http.Request) error {
		if request.URL.Scheme != "https" {
			return fmt.Errorf("refusing non-HTTPS redirect to %q", request.URL.String())
		}
		if len(previous) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		if previousRedirectCheck != nil {
			return previousRedirectCheck(request, previous)
		}
		return nil
	}

	return &Downloader{
		cacheDir: cacheDir,
		client:   &configuredClient,
		locks:    make(map[string]*sync.Mutex),
	}
}

func CacheDirectory(target platform.Target, revision string) (string, error) {
	if revision == "" || revision == "." || revision == ".." || strings.ContainsAny(revision, `/\`) {
		return "", fmt.Errorf("invalid artifact revision %q", revision)
	}
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("locate user cache directory: %w", err)
	}
	return filepath.Join(base, "CalculatePiScripts", target.OS+"-"+target.Arch, revision), nil
}

func (downloader *Downloader) Ensure(ctx context.Context, file catalog.Download) (string, error) {
	destination, err := downloader.destination(file.Path)
	if err != nil {
		return "", err
	}
	if err := validateURL(file.URL); err != nil {
		return "", err
	}
	if file.GitBlobSHA1 != "" {
		if len(file.GitBlobSHA1) != sha1.Size*2 {
			return "", fmt.Errorf("invalid Git blob SHA-1 for %s", file.Path)
		}
		if _, err := hex.DecodeString(file.GitBlobSHA1); err != nil {
			return "", fmt.Errorf("invalid Git blob SHA-1 for %s: %w", file.Path, err)
		}
	}

	lock := downloader.pathLock(destination)
	lock.Lock()
	defer lock.Unlock()

	if validCachedFile(destination, file.GitBlobSHA1) {
		return destination, nil
	}

	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return "", fmt.Errorf("create artifact cache directory: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, file.URL, nil)
	if err != nil {
		return "", fmt.Errorf("create download request for %s: %w", file.Path, err)
	}
	response, err := downloader.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("download %s: %w", file.Path, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download %s: server returned %s", file.Path, response.Status)
	}
	if response.ContentLength > maximumSize {
		return "", fmt.Errorf("download %s exceeds the %d-byte size limit", file.Path, maximumSize)
	}

	temporary, err := os.CreateTemp(filepath.Dir(destination), ".pi-download-*")
	if err != nil {
		return "", fmt.Errorf("create temporary artifact for %s: %w", file.Path, err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	written, copyErr := io.Copy(temporary, io.LimitReader(response.Body, maximumSize+1))
	if copyErr != nil {
		temporary.Close()
		return "", fmt.Errorf("write temporary artifact for %s: %w", file.Path, copyErr)
	}
	if written > maximumSize {
		temporary.Close()
		return "", fmt.Errorf("download %s exceeds the %d-byte size limit", file.Path, maximumSize)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return "", fmt.Errorf("sync temporary artifact for %s: %w", file.Path, err)
	}
	mode := os.FileMode(0o644)
	if file.Executable {
		mode = 0o755
	}
	if err := temporary.Chmod(mode); err != nil {
		temporary.Close()
		return "", fmt.Errorf("set permissions for %s: %w", file.Path, err)
	}
	if err := temporary.Close(); err != nil {
		return "", fmt.Errorf("close temporary artifact for %s: %w", file.Path, err)
	}
	if file.GitBlobSHA1 != "" {
		if err := verifyGitBlobSHA1(temporaryPath, file.GitBlobSHA1); err != nil {
			return "", fmt.Errorf("verify downloaded artifact %s: %w", file.Path, err)
		}
	}
	if err := replaceFile(temporaryPath, destination); err != nil {
		return "", fmt.Errorf("finalize downloaded artifact %s: %w", file.Path, err)
	}
	return destination, nil
}

func (downloader *Downloader) destination(relativePath string) (string, error) {
	if downloader.cacheDir == "" {
		return "", errors.New("artifact cache directory is empty")
	}
	if filepath.IsAbs(relativePath) {
		return "", fmt.Errorf("artifact path must be relative: %q", relativePath)
	}
	clean := filepath.Clean(filepath.FromSlash(relativePath))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("artifact path escapes the cache directory: %q", relativePath)
	}
	root, err := filepath.Abs(downloader.cacheDir)
	if err != nil {
		return "", fmt.Errorf("resolve artifact cache directory: %w", err)
	}
	destination := filepath.Join(root, clean)
	relative, err := filepath.Rel(root, destination)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("artifact path escapes the cache directory: %q", relativePath)
	}
	return destination, nil
}

func (downloader *Downloader) pathLock(path string) *sync.Mutex {
	downloader.mutex.Lock()
	defer downloader.mutex.Unlock()
	lock, ok := downloader.locks[path]
	if !ok {
		lock = &sync.Mutex{}
		downloader.locks[path] = lock
	}
	return lock
}

func validateURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid artifact URL: %w", err)
	}
	if parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil {
		return fmt.Errorf("artifact URL must use HTTPS and include a valid host: %q", rawURL)
	}
	return nil
}

func validCachedFile(path, expectedSHA1 string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		return false
	}
	if expectedSHA1 == "" {
		return true
	}
	return verifyGitBlobSHA1(path, expectedSHA1) == nil
}

func verifyGitBlobSHA1(path, expected string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}
	hasher := sha1.New()
	if _, err := fmt.Fprintf(hasher, "blob %d\x00", info.Size()); err != nil {
		return err
	}
	if _, err := io.Copy(hasher, file); err != nil {
		return err
	}
	actual := hex.EncodeToString(hasher.Sum(nil))
	if !strings.EqualFold(actual, expected) {
		return fmt.Errorf("Git blob SHA-1 mismatch: got %s, want %s", actual, expected)
	}
	return nil
}

func replaceFile(source, destination string) error {
	if _, err := os.Lstat(destination); errors.Is(err, os.ErrNotExist) {
		return os.Rename(source, destination)
	} else if err != nil {
		return err
	}

	backup, err := os.CreateTemp(filepath.Dir(destination), ".pi-previous-*")
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
		restoreErr := os.Rename(backupPath, destination)
		if restoreErr != nil {
			return errors.Join(err, fmt.Errorf("restore previous cache file: %w", restoreErr))
		}
		return err
	}
	return os.Remove(backupPath)
}
