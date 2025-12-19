package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"hecate-shell/internal/config"

	"github.com/spf13/cobra"
)

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Theme management commands",
	Long:  `Manage HecateShell themes, including reloading and generating from sources.`,
}

var themeReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload the current theme",
	Long: `Regenerate colors from theme.json using matugen.

The shell will automatically hot-reload the new colors.`,
	RunE: runThemeReload,
}

func init() {
	rootCmd.AddCommand(themeCmd)
	themeCmd.AddCommand(themeReloadCmd)
}

func runThemeReload(cmd *cobra.Command, args []string) error {
	// Check if installed
	if !config.IsInstalled() {
		return fmt.Errorf("HecateShell is not installed. Run 'hecate install' first")
	}

	themeFile, err := config.GetThemeFile()
	if err != nil {
		return fmt.Errorf("failed to get theme file path: %w", err)
	}

	fmt.Println("Reloading theme from theme.json...")

	// Run matugen
	matugenCmd := exec.Command("matugen", "json", themeFile, "-m", "dark", "--continue-on-error")
	matugenCmd.Stdout = os.Stdout
	matugenCmd.Stderr = os.Stderr

	if err := matugenCmd.Run(); err != nil {
		return fmt.Errorf("failed to run matugen: %w", err)
	}

	fmt.Println("Theme reloaded! Shell will auto-update within 1 second.")
	return nil
}
