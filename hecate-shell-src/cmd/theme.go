package cmd

import (
	"fmt"

	"hecate-shell/internal/config"
	"hecate-shell/internal/hooks"
	"hecate-shell/internal/niri"

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
	Long: `Reload theme by updating niri colors and running post-theme hooks.

The shell will automatically hot-reload from theme.json.`,
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

	fmt.Println("Reloading theme...")

	// Update niri config colors
	if err := niri.UpdateNiriColors(); err != nil {
		fmt.Printf("Warning: failed to update niri colors: %v\n", err)
	} else {
		fmt.Println("Niri colors updated!")
	}

	// Run post-theme hooks (pywalfox, etc.)
	hooks.RunPostThemeHooks()

	return nil
}
