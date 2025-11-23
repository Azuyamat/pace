package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	gear "github.com/azuyamat/gear/command"
	"github.com/azuyamat/pace/internal/logger"
)

const (
	installScriptUnix    = "https://raw.githubusercontent.com/azuyamat/pace/master/install.sh"
	installScriptWindows = "https://raw.githubusercontent.com/azuyamat/pace/master/install.ps1"
)

var updateCommand = gear.NewExecutableCommand("update", "Update pace to the latest version").
	Handler(updateHandler)

func init() {
	RootCommand.AddChild(updateCommand)
}

func updateHandler(ctx *gear.Context, args gear.ValidatedArgs) error {
	logger.Info("Checking for updates...")

	if cmd := getManagedUpdateCommand(); cmd != nil {
		return runManagedUpdate(cmd)
	}

	logger.Info("Running installation script to update...")
	return runInstallScript()
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
	logger.Task("Running: %s", strings.Join(cmd.Args, " "))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("update command failed: %w", err)
	}

	logger.Success("Update completed successfully!")
	return nil
}

func runInstallScript() error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		script := fmt.Sprintf("irm %s | iex", installScriptWindows)
		cmd = exec.Command("powershell", "-NoProfile", "-Command", script)
	} else {
		if !commandExists("curl") {
			return fmt.Errorf("curl is required for updates but not found in PATH")
		}
		script := fmt.Sprintf("curl -sSL %s | sh", installScriptUnix)
		cmd = exec.Command("sh", "-c", script)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installation script failed: %w", err)
	}

	logger.Success("Update completed successfully!")
	return nil
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
