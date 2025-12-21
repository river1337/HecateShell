package matugen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"hecate-shell/internal/config"
)

// Inspired from https://github.com/AvengeMedia/DankMaterialShell/blob/master/core/internal/matugen/matugen.go

// RunMatugen executes matugen with a merged config (user config + HecateShell template)
func RunMatugen(sourceType string, sourcePath string) error {
	// Create temporary merged config
	configPath, cleanup, err := createMergedConfig()
	if err != nil {
		return fmt.Errorf("failed to create matugen config: %w", err)
	}
	defer cleanup()

	// Run matugen with the merged config
	var cmd *exec.Cmd
	if sourceType == "image" {
		cmd = exec.Command("matugen", "image", sourcePath, "-t", "scheme-fidelity", "-m", "dark", "-c", configPath, "--continue-on-error")
	} else if sourceType == "json" {
		cmd = exec.Command("matugen", "json", sourcePath, "-t", "scheme-fidelity", "-m", "dark", "-c", configPath, "--continue-on-error")
	} else {
		return fmt.Errorf("unknown source type: %s", sourceType)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// createMergedConfig creates a temporary config that merges user's matugen config with HecateShell's template
func createMergedConfig() (string, func(), error) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "hecate-matugen-*.toml")
	if err != nil {
		return "", nil, err
	}
	tmpPath := tmpFile.Name()

	cleanup := func() {
		os.Remove(tmpPath)
	}

	// Get user's matugen config (if it exists)
	userConfigPath := filepath.Join(os.Getenv("HOME"), ".config/matugen/config.toml")
	userConfig := ""
	if data, err := os.ReadFile(userConfigPath); err == nil {
		userConfig = string(data)
	}

	// Get HecateShell config dir
	shellDir, err := config.GetConfigDir()
	if err != nil {
		cleanup()
		return "", nil, err
	}

	// HecateShell template configs
	// Some templates adapted from DankMaterialShell (https://github.com/AvengeMedia/DankMaterialShell)
	homeDir := os.Getenv("HOME")
	hecateTemplates := fmt.Sprintf(`
# HecateShell core theme
[templates.hecate]
input_path = "%s/config/templates/hecate.json"
output_path = "%s/theme.json"

# Audio visualizer
[templates.hecate_cava]
input_path = "%s/config/templates/cava.ini"
output_path = "%s/.config/cava/config"

# Spotify
[templates.hecate_spicetify]
input_path = "%s/config/templates/spicetify.ini"
output_path = "%s/.config/spicetify/Themes/text/color.ini"

# Discord (Vencord)
[templates.hecate_discord]
input_path = "%s/config/templates/discord.css"
output_path = "%s/.config/Vencord/themes/sys24.css"

# Micro editor
[templates.hecate_micro]
input_path = "%s/config/templates/micro.micro"
output_path = "%s/.config/micro/colorschemes/matugen.micro"

# VSCode theme
[templates.hecate_vscode]
input_path = "%s/config/templates/vscode.json"
output_path = "%s/.vscode/extensions/hecate-theme/themes/hecate-dark.json"

# Niri compositor colors
[templates.hecate_niri]
input_path = "%s/config/templates/niri-colors.kdl"
output_path = "%s/.config/niri/hecate-colors.generated.kdl"

# Firefox (pywalfox)
[templates.hecate_pywalfox]
input_path = "%s/config/templates/pywalfox.json"
output_path = "%s/.cache/wal/colors.json"

# Kitty terminal
[templates.hecate_kitty]
input_path = "%s/config/templates/kitty.conf"
output_path = "%s/.config/kitty/hecate-colors.conf"

# Kitty tabs
[templates.hecate_kitty_tabs]
input_path = "%s/config/templates/kitty-tabs.conf"
output_path = "%s/.config/kitty/hecate-tabs.conf"

# Alacritty terminal
[templates.hecate_alacritty]
input_path = "%s/config/templates/alacritty.toml"
output_path = "%s/.config/alacritty/hecate-colors.toml"

# KDE color scheme
[templates.hecate_kde]
input_path = "%s/config/templates/kcolorscheme.colors"
output_path = "%s/.local/share/color-schemes/HecateShell.colors"

# Qt5ct/Qt6ct colors
[templates.hecate_qt]
input_path = "%s/config/templates/qt5ct-colors.conf"
output_path = "%s/.config/qt5ct/colors/HecateShell.conf"
`, shellDir, shellDir,
		shellDir, homeDir,
		shellDir, homeDir,
		shellDir, homeDir,
		shellDir, homeDir,
		shellDir, homeDir,
		shellDir, homeDir,
		shellDir, homeDir,
		shellDir, homeDir,
		shellDir, homeDir,
		shellDir, homeDir,
		shellDir, homeDir,
		shellDir, homeDir)

	// Merge configs: user config + HecateShell templates
	mergedConfig := userConfig + "\n" + hecateTemplates

	// Write to temp file
	if _, err := tmpFile.WriteString(mergedConfig); err != nil {
		cleanup()
		return "", nil, err
	}

	if err := tmpFile.Close(); err != nil {
		cleanup()
		return "", nil, err
	}

	return tmpPath, cleanup, nil
}
