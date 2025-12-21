package hooks

import (
	"fmt"
	"os/exec"
)

// RunPostThemeHooks runs commands that need to execute after theme generation
func RunPostThemeHooks() {
	// Update pywalfox (Firefox theming)
	if err := runPywalfoxUpdate(); err != nil {
		fmt.Printf("Warning: pywalfox update failed: %v\n", err)
	}
}

// runPywalfoxUpdate calls pywalfox update to apply new colors to Firefox
func runPywalfoxUpdate() error {
	// Check if pywalfox is installed
	if _, err := exec.LookPath("pywalfox"); err != nil {
		// pywalfox not installed, skip silently
		return nil
	}

	cmd := exec.Command("pywalfox", "update")
	return cmd.Run()
}
