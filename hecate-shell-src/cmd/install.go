package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"hecate-shell/internal/config"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install HecateShell configuration",
	Long: `Install HecateShell by cloning the repository to ~/.config/HecateShell.

This will:
  - Clone the HecateShell repository
  - Set up all configuration files
  - Prepare the shell for first run`,
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolP("force", "f", false, "Force reinstall (removes existing installation)")
}

func runInstall(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	// Check if already installed
	if config.IsInstalled() && !force {
		return fmt.Errorf("HecateShell is already installed. Use --force to reinstall")
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// If force, remove existing installation
	if force {
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

	fmt.Println("\nHecateShell installed successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Make sure you have dependencies installed:")
	fmt.Println("     - quickshell-git (for Niri support)")
	fmt.Println("     - cava (audio visualizer)")
	fmt.Println("     - pipewire + wireplumber (audio control)")
	fmt.Println("     - matugen (theme generation)")
	fmt.Println("     - swww (wallpaper daemon)")
	fmt.Println("\n  2. Start the shell:")
	fmt.Println("     hecate run")

	return nil
}
