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

	// HecateShell template config
	hecateTemplate := fmt.Sprintf(`
[templates.hecate]
input_path = "%s/config/templates/hecate.json"
output_path = "%s/theme.json"
`, shellDir, shellDir)

	// Merge configs: user config + HecateShell template
	mergedConfig := userConfig + "\n" + hecateTemplate

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
