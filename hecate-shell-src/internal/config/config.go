package config

import (
	"os"
	"path/filepath"
)

const (
	ConfigDirName = "HecateShell"
	RepoURL      = "https://github.com/river1337/HecateShell.git"
	VersionURL   = "https://raw.githubusercontent.com/river1337/HecateShell/refs/heads/main/version"
)

// GetConfigDir returns the path to the HecateShell config directory
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", ConfigDirName), nil
}

// IsInstalled checks if HecateShell is installed
func IsInstalled() bool {
	configDir, err := GetConfigDir()
	if err != nil {
		return false
	}
	
	// Check if directory exists and has shell.qml
	shellFile := filepath.Join(configDir, "shell.qml")
	_, err = os.Stat(shellFile)
	return err == nil
}

// GetVersionFile returns the path to the version file
func GetVersionFile() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "version"), nil
}

// GetShellFile returns the path to shell.qml
func GetShellFile() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "shell.qml"), nil
}

// GetThemeFile returns the path to theme.json
func GetThemeFile() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "theme.json"), nil
}
