package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"hecate-shell/internal/config"
	"hecate-shell/internal/installer"
	"hecate-shell/internal/vscode"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install or Update HecateShell configuration",
	Long: `Install or Update HecateShell with an interactive installer.

This will guide you through:
  - Installing dependencies (optional)
  - Installing dotfiles (optional)
  - Setting up the HecateShell QuickShell configuration

If HecateShell is already installed, it will offer to update instead.`,
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolP("force", "f", false, "Force reinstall (removes existing installation)")
	installCmd.Flags().Bool("legacy", false, "Use legacy (non-interactive) installer")
}

func runInstall(cmd *cobra.Command, args []string) error {
	legacy, _ := cmd.Flags().GetBool("legacy")

	if legacy {
		return runLegacyInstall(cmd, args)
	}

	// Run the interactive TUI installer
	return installer.Run()
}

// runLegacyInstall is the old simple installer (kept for scripting/automation)
func runLegacyInstall(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	// Check if already installed
	if config.IsInstalled() {
		if !force {
			// Prompt user if they want to update
			fmt.Println("HecateShell is already installed.")
			fmt.Print("Would you like to update instead? [y/N]: ")

			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response == "y" || response == "yes" {
				return runLegacyUpdate()
			}

			fmt.Println("Use --force to reinstall.")
			return nil
		}
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// If force, remove existing installation
	if force && config.IsInstalled() {
		fmt.Println("Removing existing installation...")
		if err := os.RemoveAll(configDir); err != nil {
			return fmt.Errorf("failed to remove existing installation: %w", err)
		}
	}

	// Clone repository
	fmt.Println("Installing HecateShell...")
	fmt.Printf("   Cloning from: %s\n", config.RepoURL)
	fmt.Printf("   Installing to: %s\n", configDir)

	gitCmd := exec.Command("git", "clone", config.RepoURL, configDir)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Install VSCode theme extension if VSCode is present
	if err := vscode.InstallThemeExtension(); err != nil {
		fmt.Printf("Warning: failed to install VSCode theme: %v\n", err)
	} else {
		fmt.Println("VSCode theme extension installed!")
	}

	fmt.Println("\nHecateShell installed successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Make sure you have dependencies installed:")
	fmt.Println("     - quickshell-git (for Niri support)")
	fmt.Println("     - cava (audio visualizer)")
	fmt.Println("     - pipewire + wireplumber (audio control)")
	fmt.Println("     - matugen (theme generation)")
	fmt.Println("     - swww (wallpaper daemon)")
	fmt.Println("     - ttf-jetbrains-mono-nerd (font)")
	fmt.Println("\n  2. Start the shell:")
	fmt.Println("     hecate run")

	return nil
}

// runLegacyUpdate performs a git pull update
func runLegacyUpdate() error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	fmt.Println("Updating HecateShell...")

	// Pull latest changes
	gitCmd := exec.Command("git", "-C", configDir, "pull", "origin", "main")
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("failed to pull updates: %w", err)
	}

	// Ensure VSCode theme extension is installed
	if err := vscode.InstallThemeExtension(); err != nil {
		fmt.Printf("Warning: failed to install VSCode theme: %v\n", err)
	}

	fmt.Println("\nUpdated successfully!")
	fmt.Println("\nIf the shell is running, reload it with:")
	fmt.Println("  hecate run --reload")

	return nil
}
