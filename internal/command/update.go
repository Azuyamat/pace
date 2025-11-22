package command

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	gear "github.com/azuyamat/gear/command"
	"github.com/azuyamat/pace/internal/version"
)

const githubAPI = "https://api.github.com/repos/azuyamat/pace/releases/latest"

type release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var updateCommand = gear.NewExecutableCommand("update", "Update pace to the latest version").
	Handler(updateHandler)

func init() {
	RootCommand.AddChild(updateCommand)
}

func updateHandler(ctx *gear.Context, args gear.ValidatedArgs) error {
	fmt.Println("Checking for updates...")

	rel, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := strings.TrimPrefix(rel.TagName, "v")
	if version.Version == latestVersion {
		fmt.Printf("Already on latest version: %s\n", version.Version)
		return nil
	}

	if version.Version != "dev" {
		fmt.Printf("Updating from %s to %s...\n", version.Version, latestVersion)
	}

	if cmd := getManagedUpdateCommand(); cmd != nil {
		return runManagedUpdate(cmd)
	}

	fmt.Println("Proceeding with self-update...")

	assetURL := findAssetURL(rel)
	if assetURL == "" {
		return fmt.Errorf("no compatible release found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	if err := replaceExecutable(assetURL); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Printf("\nSuccessfully updated to version %s\n", latestVersion)
	fmt.Println("Restart pace to use the new version")
	return nil
}

func getManagedUpdateCommand() *exec.Cmd {
	exe, err := os.Executable()
	if err != nil {
		return nil
	}

	path := filepath.Clean(exe)

	if runtime.GOOS == "windows" {
		if pf := os.Getenv("ProgramFiles"); pf != "" && strings.HasPrefix(path, filepath.Clean(pf)) {
			return exec.Command("winget", "upgrade", "Azuyamat.Pace")
		}
		if la := os.Getenv("LOCALAPPDATA"); la != "" && strings.Contains(path, filepath.Join(la, "Microsoft", "WinGet")) {
			return exec.Command("winget", "upgrade", "Azuyamat.Pace")
		}
	}

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		if home, err := os.UserHomeDir(); err == nil {
			goPath = filepath.Join(home, "go")
		}
	}

	if goPath != "" && strings.HasPrefix(path, filepath.Join(goPath, "bin")) {
		return exec.Command("go", "install", "github.com/azuyamat/pace/cmd/pace@latest")
	}

	if goBin := os.Getenv("GOBIN"); goBin != "" && strings.HasPrefix(path, filepath.Clean(goBin)) {
		return exec.Command("go", "install", "github.com/azuyamat/pace/cmd/pace@latest")
	}

	return nil
}

func runManagedUpdate(cmd *exec.Cmd) error {
	fmt.Printf("Running: %s\n", strings.Join(cmd.Args, " "))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("update command failed: %w", err)
	}

	fmt.Println("\nUpdate completed successfully!")
	return nil
}

func getLatestRelease() (*release, error) {
	resp, err := http.Get(githubAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}

	return &rel, nil
}

func findAssetURL(rel *release) string {
	pattern := fmt.Sprintf("pace-%s-%s", runtime.GOOS, runtime.GOARCH)
	for _, asset := range rel.Assets {
		if strings.Contains(asset.Name, pattern) {
			return asset.BrowserDownloadURL
		}
	}
	return ""
}

func replaceExecutable(url string) error {
	archive, err := downloadFile(url)
	if err != nil {
		return err
	}
	defer os.Remove(archive)

	binary, err := extractBinary(archive)
	if err != nil {
		return err
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	backup := exe + ".backup"
	if err := os.Rename(exe, backup); err != nil {
		return err
	}

	if err := os.WriteFile(exe, binary, 0755); err != nil {
		os.Rename(backup, exe)
		return err
	}

	os.Remove(backup)
	return nil
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmp, err := os.CreateTemp("", "pace-*.tmp")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	return tmp.Name(), nil
}

func extractBinary(archive string) ([]byte, error) {
	if strings.HasSuffix(archive, ".zip") {
		return extractFromZip(archive)
	}
	return extractFromTarGz(archive)
}

func extractFromZip(path string) ([]byte, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	name := getBinaryName()
	for _, f := range r.File {
		if filepath.Base(f.Name) == name {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}

	return nil, fmt.Errorf("binary not found in archive")
}

func extractFromTarGz(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	name := getBinaryName()

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if filepath.Base(header.Name) == name {
			return io.ReadAll(tr)
		}
	}

	return nil, fmt.Errorf("binary not found in archive")
}

func getBinaryName() string {
	if runtime.GOOS == "windows" {
		return "pace.exe"
	}
	return "pace"
}
