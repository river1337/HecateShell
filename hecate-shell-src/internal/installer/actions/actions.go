package actions

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"hecate-shell/internal/config"
)

// Dependencies for HecateShell
var Dependencies = []string{
	"quickshell-git",
	"cava",
	"pipewire",
	"wireplumber",
	"matugen-bin",
	"swww",
	"ttf-jetbrains-mono-nerd",
}

// DotfileMapping maps friendly names to their source/dest paths
var DotfileMapping = map[string]struct {
	Source string // Relative to HecateShell repo
	Dest   string // Relative to ~/.config/
}{
	"Niri":      {Source: "dotfiles/niri", Dest: "niri"},
	"Fish":      {Source: "dotfiles/fish", Dest: "fish"},
	"Kitty":     {Source: "dotfiles/kitty", Dest: "kitty"},
	"Micro":     {Source: "dotfiles/micro", Dest: "micro"},
	"Fastfetch": {Source: "dotfiles/fastfetch", Dest: "fastfetch"},
	"Neovim":    {Source: "dotfiles/nvim", Dest: "nvim"},
}

// TaskResult holds the result of an installation task
type TaskResult struct {
	Success bool
	Message string
	Error   error
}

// ProgressCallback is called with progress updates
type ProgressCallback func(current, total int, message string)

// InstallDependencies installs the required packages
func InstallDependencies(pkgManager string, progress ProgressCallback) TaskResult {
	// Build command based on package manager
	var cmd *exec.Cmd
	switch pkgManager {
	case "paru":
		cmd = exec.Command("paru", append([]string{"-S", "--noconfirm", "--needed"}, Dependencies...)...)
	case "yay":
		cmd = exec.Command("yay", append([]string{"-S", "--noconfirm", "--needed"}, Dependencies...)...)
	case "pacman":
		// Filter out AUR packages for pacman
		officialPkgs := []string{"cava", "pipewire", "wireplumber", "ttf-jetbrains-mono-nerd"}
		cmd = exec.Command("sudo", append([]string{"pacman", "-S", "--noconfirm", "--needed"}, officialPkgs...)...)
	default:
		return TaskResult{Success: false, Error: fmt.Errorf("unknown package manager: %s", pkgManager)}
	}

	if progress != nil {
		progress(0, 1, "Installing packages with "+pkgManager+"...")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return TaskResult{
			Success: false,
			Message: string(output),
			Error:   fmt.Errorf("failed to install dependencies: %w", err),
		}
	}

	if progress != nil {
		progress(1, 1, "Dependencies installed!")
	}

	return TaskResult{Success: true, Message: "Dependencies installed successfully"}
}

// InstallDotfile installs a single dotfile configuration
func InstallDotfile(name string, progress ProgressCallback) TaskResult {
	mapping, ok := DotfileMapping[name]
	if !ok {
		return TaskResult{Success: false, Error: fmt.Errorf("unknown dotfile: %s", name)}
	}

	if progress != nil {
		progress(0, 1, fmt.Sprintf("Installing %s config...", name))
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return TaskResult{Success: false, Error: err}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return TaskResult{Success: false, Error: err}
	}

	sourcePath := filepath.Join(configDir, mapping.Source)
	destPath := filepath.Join(homeDir, ".config", mapping.Dest)

	// Check if source exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return TaskResult{
			Success: false,
			Error:   fmt.Errorf("source dotfiles not found: %s", sourcePath),
		}
	}

	// Backup existing config if it exists
	if _, err := os.Stat(destPath); err == nil {
		backupPath := destPath + ".backup"
		// Remove old backup if exists
		os.RemoveAll(backupPath)
		if err := os.Rename(destPath, backupPath); err != nil {
			return TaskResult{
				Success: false,
				Error:   fmt.Errorf("failed to backup existing config: %w", err),
			}
		}
	}

	// Copy dotfiles
	if err := copyDir(sourcePath, destPath); err != nil {
		return TaskResult{Success: false, Error: fmt.Errorf("failed to copy dotfiles: %w", err)}
	}

	if progress != nil {
		progress(1, 1, fmt.Sprintf("%s config installed!", name))
	}

	return TaskResult{Success: true, Message: fmt.Sprintf("%s configuration installed", name)}
}

// InstallShell downloads and extracts HecateShell from the latest release
func InstallShell(force bool, progress ProgressCallback) TaskResult {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return TaskResult{Success: false, Error: err}
	}

	if progress != nil {
		progress(0, 3, "Downloading HecateShell...")
	}

	// Check if already installed
	if config.IsInstalled() {
		if !force {
			return TaskResult{
				Success: false,
				Error:   fmt.Errorf("HecateShell is already installed"),
			}
		}
		// Remove existing installation
		if err := os.RemoveAll(configDir); err != nil {
			return TaskResult{
				Success: false,
				Error:   fmt.Errorf("failed to remove existing installation: %w", err),
			}
		}
	}

	// Create config directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return TaskResult{
			Success: false,
			Error:   fmt.Errorf("failed to create config directory: %w", err),
		}
	}

	// Download the latest release archive
	// GitHub redirects /releases/latest/download/ to the actual release
	archiveURL := config.ReleaseURL
	resp, err := http.Get(archiveURL)
	if err != nil {
		return TaskResult{
			Success: false,
			Error:   fmt.Errorf("failed to download release: %w", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TaskResult{
			Success: false,
			Error:   fmt.Errorf("failed to download release: HTTP %d", resp.StatusCode),
		}
	}

	if progress != nil {
		progress(1, 3, "Extracting files...")
	}

	// Extract the tar.gz archive
	if err := extractTarGz(resp.Body, configDir); err != nil {
		return TaskResult{
			Success: false,
			Error:   fmt.Errorf("failed to extract archive: %w", err),
		}
	}

	if progress != nil {
		progress(2, 3, "Verifying installation...")
	}

	// Verify installation
	if !config.IsInstalled() {
		return TaskResult{
			Success: false,
			Error:   fmt.Errorf("installation verification failed"),
		}
	}

	if progress != nil {
		progress(3, 3, "HecateShell installed!")
	}

	return TaskResult{Success: true, Message: "HecateShell installed successfully"}
}

// InstallShellDev clones the full HecateShell repository (for development)
func InstallShellDev(force bool, progress ProgressCallback) TaskResult {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return TaskResult{Success: false, Error: err}
	}

	if progress != nil {
		progress(0, 2, "Cloning HecateShell repository...")
	}

	// Check if already installed
	if config.IsInstalled() {
		if !force {
			return TaskResult{
				Success: false,
				Error:   fmt.Errorf("HecateShell is already installed"),
			}
		}
		// Remove existing installation
		if err := os.RemoveAll(configDir); err != nil {
			return TaskResult{
				Success: false,
				Error:   fmt.Errorf("failed to remove existing installation: %w", err),
			}
		}
	}

	// Clone repository
	cmd := exec.Command("git", "clone", config.RepoURL, configDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return TaskResult{
			Success: false,
			Message: string(output),
			Error:   fmt.Errorf("failed to clone repository: %w", err),
		}
	}

	if progress != nil {
		progress(1, 2, "Verifying installation...")
	}

	// Verify installation
	if !config.IsInstalled() {
		return TaskResult{
			Success: false,
			Error:   fmt.Errorf("installation verification failed"),
		}
	}

	if progress != nil {
		progress(2, 2, "HecateShell installed (dev mode)!")
	}

	return TaskResult{Success: true, Message: "HecateShell installed successfully (dev mode)"}
}

// extractTarGz extracts a tar.gz archive to the destination directory
func extractTarGz(r io.Reader, dest string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Clean the path and make it relative to dest
		name := filepath.Clean(header.Name)
		// Remove leading ./ if present
		name = strings.TrimPrefix(name, "./")
		if name == "." || name == "" {
			continue
		}

		target := filepath.Join(dest, name)

		// Ensure the target is within the destination directory (security check)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dest)) {
			return fmt.Errorf("invalid file path in archive: %s", name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}

		case tar.TypeReg:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", target, err)
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("failed to write file %s: %w", target, err)
			}
			f.Close()

		case tar.TypeSymlink:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for symlink %s: %w", target, err)
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("failed to create symlink %s: %w", target, err)
			}
		}
	}

	return nil
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}

	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// CheckCommand checks if a command exists
func CheckCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// GetMissingDependencies returns a list of missing dependencies
func GetMissingDependencies() []string {
	var missing []string

	// Check for each dependency's binary
	binaryChecks := map[string]string{
		"quickshell-git":          "quickshell",
		"cava":                    "cava",
		"pipewire":                "pipewire",
		"wireplumber":             "wpctl",
		"matugen-bin":             "matugen",
		"swww":                    "swww",
		"ttf-jetbrains-mono-nerd": "", // Font, no binary to check
	}

	for pkg, binary := range binaryChecks {
		if binary != "" && !CheckCommand(binary) {
			missing = append(missing, pkg)
		}
	}

	return missing
}

// RunPostInstall runs post-installation tasks
func RunPostInstall(progress ProgressCallback) TaskResult {
	if progress != nil {
		progress(0, 1, "Finishing up...")
	}

	// Any post-install tasks can go here
	// For example, setting up default theme, etc.

	if progress != nil {
		progress(1, 1, "Done!")
	}

	return TaskResult{Success: true, Message: "Post-installation complete"}
}

// UpdateShell performs a git pull to update HecateShell
func UpdateShell(progress ProgressCallback) TaskResult {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return TaskResult{Success: false, Error: err}
	}

	if progress != nil {
		progress(0, 2, "Checking for updates...")
	}

	// Check if directory exists
	if !config.IsInstalled() {
		return TaskResult{
			Success: false,
			Error:   fmt.Errorf("HecateShell is not installed"),
		}
	}

	if progress != nil {
		progress(1, 2, "Pulling latest changes...")
	}

	// Git pull
	cmd := exec.Command("git", "-C", configDir, "pull", "origin", "main")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return TaskResult{
			Success: false,
			Message: string(output),
			Error:   fmt.Errorf("failed to pull updates: %w", err),
		}
	}

	if progress != nil {
		progress(2, 2, "Updated!")
	}

	return TaskResult{Success: true, Message: "HecateShell updated successfully"}
}

// BackupDotfile creates a backup of an existing dotfile config
func BackupDotfile(name string) (string, error) {
	mapping, ok := DotfileMapping[name]
	if !ok {
		return "", fmt.Errorf("unknown dotfile: %s", name)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	destPath := filepath.Join(homeDir, ".config", mapping.Dest)

	// Check if config exists
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return "", nil // No existing config to backup
	}

	// Create backup with timestamp
	timestamp := strings.Replace(strings.Replace(
		strings.Split(fmt.Sprintf("%v", time.Now()), ".")[0],
		" ", "_", -1), ":", "-", -1)
	backupPath := destPath + ".backup." + timestamp

	if err := copyDir(destPath, backupPath); err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}

	return backupPath, nil
}

// FormatDependencyList returns a formatted string of dependencies
func FormatDependencyList() string {
	var sb strings.Builder
	for i, dep := range Dependencies {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(dep)
	}
	return sb.String()
}
