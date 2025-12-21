package niri

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"hecate-shell/internal/config"
)

// Theme represents the theme.json structure (partial)
type Theme struct {
	Primary string `json:"primary"`
	Outline string `json:"outline"`
	Error   string `json:"error"`
	Shadow  string `json:"shadow"`
}

// UpdateNiriColors reads the theme.json and updates niri config colors
func UpdateNiriColors() error {
	// Read theme.json
	theme, err := readTheme()
	if err != nil {
		return fmt.Errorf("failed to read theme: %w", err)
	}

	// Read niri config
	niriConfigPath := filepath.Join(os.Getenv("HOME"), ".config/niri/config.kdl")
	file, err := os.Open(niriConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read niri config: %w", err)
	}

	var lines []string
	scanner := bufio.NewScanner(file)

	// Track which block we're in
	inFocusRing := false
	inBorder := false
	inShadow := false

	// Regex patterns for color lines
	activeColorRe := regexp.MustCompile(`^(\s*)active-color\s+"#[0-9a-fA-F]+"`)
	inactiveColorRe := regexp.MustCompile(`^(\s*)inactive-color\s+"#[0-9a-fA-F]+"`)
	urgentColorRe := regexp.MustCompile(`^(\s*)urgent-color\s+"#[0-9a-fA-F]+"`)
	shadowColorRe := regexp.MustCompile(`^(\s*)color\s+"#[0-9a-fA-F]+"`)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Track block entry/exit
		if strings.HasPrefix(trimmed, "focus-ring") && strings.Contains(trimmed, "{") {
			inFocusRing = true
		} else if strings.HasPrefix(trimmed, "border") && strings.Contains(trimmed, "{") {
			inBorder = true
		} else if strings.HasPrefix(trimmed, "shadow") && strings.Contains(trimmed, "{") {
			inShadow = true
		}

		// Check for block exit (closing brace on its own line within these blocks)
		if (inFocusRing || inBorder || inShadow) && trimmed == "}" {
			inFocusRing = false
			inBorder = false
			inShadow = false
		}

		// Replace colors in focus-ring or border blocks
		if inFocusRing || inBorder {
			indent := getIndent(line)
			if activeColorRe.MatchString(line) {
				line = fmt.Sprintf(`%sactive-color "%s"`, indent, theme.Primary)
			} else if inactiveColorRe.MatchString(line) {
				line = fmt.Sprintf(`%sinactive-color "%s"`, indent, theme.Outline)
			} else if urgentColorRe.MatchString(line) {
				line = fmt.Sprintf(`%surgent-color "%s"`, indent, theme.Error)
			}
		}

		// Replace shadow color (primary color with 80 alpha for transparency)
		if inShadow && shadowColorRe.MatchString(line) {
			indent := getIndent(line)
			line = fmt.Sprintf(`%scolor "%s80"`, indent, theme.Primary)
		}

		lines = append(lines, line)
	}
	file.Close()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading niri config: %w", err)
	}

	// Write back with trailing newline
	output := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(niriConfigPath, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write niri config: %w", err)
	}

	return nil
}

// readTheme reads the theme.json file
func readTheme() (*Theme, error) {
	shellDir, err := config.GetConfigDir()
	if err != nil {
		return nil, err
	}

	themePath := filepath.Join(shellDir, "theme.json")
	data, err := os.ReadFile(themePath)
	if err != nil {
		return nil, err
	}

	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, err
	}

	return &theme, nil
}

// getIndent extracts the leading whitespace from a line
func getIndent(line string) string {
	for i, c := range line {
		if c != ' ' && c != '\t' {
			return line[:i]
		}
	}
	return line
}
