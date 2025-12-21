package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"hecate-shell/internal/matugen"
	"hecate-shell/internal/niri"

	"github.com/spf13/cobra"
)

var wallpaperCmd = &cobra.Command{
	Use:   "wallpaper <path>",
	Short: "Set wallpaper",
	Long: `Set wallpaper using swww with optional theme generation.

Examples:
  hecate wallpaper /wallpaper.jpg
  hecate wallpaper /wallpaper.jpg --generate-theme`,
	Args: cobra.ExactArgs(1),
	RunE: runWallpaper,
}

func init() {
	rootCmd.AddCommand(wallpaperCmd)
	wallpaperCmd.Flags().BoolP("generate-theme", "g", false, "Generate theme colors from wallpaper")
	wallpaperCmd.Flags().StringP("transition", "t", "fade", "Transition effect (fade, wipe, grow, outer, center)")
	wallpaperCmd.Flags().IntP("duration", "d", 1, "Transition duration in seconds")
}

func runWallpaper(cmd *cobra.Command, args []string) error {
	wallpaperPath := args[0]
	generateTheme, _ := cmd.Flags().GetBool("generate-theme")

	// Get transition settings from flags or use defaults from config
	transition, _ := cmd.Flags().GetString("transition")
	duration, _ := cmd.Flags().GetInt("duration")

	// If not specified via flags, could read from config in the future
	// For now just use the flag defaults

	// Expand path
	if wallpaperPath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		wallpaperPath = filepath.Join(homeDir, wallpaperPath[1:])
	}

	// Check if file exists
	if _, err := os.Stat(wallpaperPath); err != nil {
		return fmt.Errorf("wallpaper file not found: %s", wallpaperPath)
	}

	// Set wallpaper with swww
	fmt.Printf("Setting wallpaper: %s\n", wallpaperPath)

	swwwCmd := exec.Command("swww", "img", wallpaperPath,
		"--transition-type", transition,
		"--transition-duration", fmt.Sprintf("%d", duration))
	swwwCmd.Stdout = os.Stdout
	swwwCmd.Stderr = os.Stderr

	if err := swwwCmd.Run(); err != nil {
		return fmt.Errorf("failed to set wallpaper with swww: %w (is swww daemon running?)", err)
	}

	fmt.Println("Wallpaper set successfully!")

	// Generate theme if flag is set
	if generateTheme {
		fmt.Println("\nGenerating theme from wallpaper colors...")

		// Run matugen with merged config (user's config + HecateShell template)
		if err := matugen.RunMatugen("image", wallpaperPath); err != nil {
			return fmt.Errorf("failed to generate theme with matugen: %w", err)
		}

		// Update niri config colors
		if err := niri.UpdateNiriColors(); err != nil {
			fmt.Printf("Warning: failed to update niri colors: %v\n", err)
		} else {
			fmt.Println("Niri colors updated!")
		}

		fmt.Println("Theme generated! Shell will auto-update within 1 second.")
	}

	return nil
}
