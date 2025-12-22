package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"hecate-shell/internal/hooks"
	"hecate-shell/internal/matugen"
	"hecate-shell/internal/niri"

	"github.com/spf13/cobra"
)

var wallpaperCmd = &cobra.Command{
	Use:   "wallpaper <path>",
	Short: "Set wallpaper",
	Long: `Set wallpaper with optional theme generation.

Examples:
  hecate wallpaper /wallpaper.jpg
  hecate wallpaper /wallpaper.jpg --generate-theme
  hecate wallpaper /wallpaper.jpg --transition fade --duration 2`,
	Args: cobra.ExactArgs(1),
	RunE: runWallpaper,
}

func init() {
	rootCmd.AddCommand(wallpaperCmd)
	wallpaperCmd.Flags().BoolP("generate-theme", "g", false, "Generate theme colors from wallpaper")
	wallpaperCmd.Flags().StringP("transition", "t", "", "Transition effect (only 'fade' is supported currently)")
	wallpaperCmd.Flags().IntP("duration", "d", 0, "Transition duration in milliseconds")
}

func runWallpaper(cmd *cobra.Command, args []string) error {
	wallpaperPath := args[0]
	generateTheme, _ := cmd.Flags().GetBool("generate-theme")
	transition, _ := cmd.Flags().GetString("transition")
	duration, _ := cmd.Flags().GetInt("duration")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Resolve wallpaper path (handle shortcuts, ~, relative, and absolute paths)
	absPath, err := resolveWallpaperPath(wallpaperPath, homeDir)
	if err != nil {
		return err
	}

	// Get config path
	configPath := filepath.Join(homeDir, ".config", "HecateShell", "config.json")

	// Read existing config
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config.json: %w", err)
	}

	// Parse config
	var config map[string]interface{}
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse config.json: %w", err)
	}

	// Update wallpaper section
	wallpaperConfig, ok := config["wallpaper"].(map[string]interface{})
	if !ok {
		wallpaperConfig = make(map[string]interface{})
		config["wallpaper"] = wallpaperConfig
	}

	wallpaperConfig["path"] = absPath
	if transition != "" {
		wallpaperConfig["transition"] = transition
	}
	if duration > 0 {
		wallpaperConfig["duration"] = duration
	}

	// Write updated config
	updatedData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write config.json: %w", err)
	}

	fmt.Printf("Wallpaper set: %s\n", absPath)

	// Generate theme if flag is set
	if generateTheme {
		fmt.Println("\nGenerating theme from wallpaper colors...")

		// Run matugen with merged config (user's config + HecateShell template)
		if err := matugen.RunMatugen("image", absPath); err != nil {
			return fmt.Errorf("failed to generate theme with matugen: %w", err)
		}

		// Update niri config colors
		if err := niri.UpdateNiriColors(); err != nil {
			fmt.Printf("Warning: failed to update niri colors: %v\n", err)
		} else {
			fmt.Println("Niri colors updated!")
		}

		// Run post-theme hooks (pywalfox, etc.)
		hooks.RunPostThemeHooks()

		fmt.Println("Theme generated! Shell will auto-update within 1 second.")
	}

	fmt.Println("Wallpaper will update within 1 second (hot-reload).")
	return nil
}

// resolveWallpaperPath resolves wallpaper path from shortcut name or file path
func resolveWallpaperPath(input, homeDir string) (string, error) {
	// If path contains "/" or starts with "~" or "./", treat as file path
	if filepath.IsAbs(input) || input[0] == '~' || input[0:2] == "./" || input[0:3] == "../" {
		// Expand ~ if present
		if len(input) > 0 && input[0] == '~' {
			input = filepath.Join(homeDir, input[1:])
		}

		// Make absolute
		absPath, err := filepath.Abs(input)
		if err != nil {
			return "", fmt.Errorf("failed to resolve path: %w", err)
		}

		// Check if file exists
		if _, err := os.Stat(absPath); err != nil {
			return "", fmt.Errorf("wallpaper file not found: %s", absPath)
		}

		return absPath, nil
	}

	// Treat as shortcut - look in ~/.config/HecateShell/wallpapers/
	shortcutDir := filepath.Join(homeDir, ".config", "HecateShell", "wallpapers")

	// Try common image extensions
	extensions := []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}
	for _, ext := range extensions {
		path := filepath.Join(shortcutDir, input+ext)
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Using wallpaper shortcut: %s -> %s\n", input, path)
			return path, nil
		}
	}

	return "", fmt.Errorf("wallpaper shortcut '%s' not found in %s (tried: %v)", input, shortcutDir, extensions)
}
