package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"hecate-shell/internal/config"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update HecateShell to the latest version",
	Long: `Check for updates and update HecateShell from the GitHub repository.

This will:
  - Check the remote version
  - Pull latest changes if a new version is available
  - Preserve your theme.json and custom configurations`,
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("force", "f", false, "Force update even if already on latest version")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	// Check if installed
	if !config.IsInstalled() {
		return fmt.Errorf("HecateShell is not installed. Run 'hecate install' first")
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Get current version
	currentVersion, err := getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	// Get remote version
	remoteVersion, err := getRemoteVersion()
	if err != nil {
		return fmt.Errorf("failed to get remote version: %w", err)
	}

	fmt.Printf("Current version: %s\n", currentVersion)
	fmt.Printf("Remote version:  %s\n", remoteVersion)

	// Check if update needed
	if currentVersion == remoteVersion && !force {
		fmt.Println("Already on the latest version!")
		return nil
	}

	if force {
		fmt.Println("\nForce updating...")
	} else {
		fmt.Println("\nUpdate available! Updating...")
	}

	// Pull latest changes
	gitCmd := exec.Command("git", "-C", configDir, "pull", "origin", "main")
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("failed to pull updates: %w", err)
	}

	// Get new version after pull
	newVersion, err := getCurrentVersion()
	if err != nil {
		newVersion = "unknown"
	}

	fmt.Printf("\nUpdated successfully to version %s!\n", newVersion)
	fmt.Println("\nIf the shell is running, reload it with:")
	fmt.Println("  hecate run --reload")

	return nil
}

func getCurrentVersion() (string, error) {
	versionFile, err := config.GetVersionFile()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(versionFile)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func getRemoteVersion() (string, error) {
	resp, err := http.Get(config.VersionURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch remote version: HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}
