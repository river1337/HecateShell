package shell

import (
	"fmt"
	"os/exec"
)

// KillQuickShell kills all running QuickShell processes
func KillQuickShell() error {
	cmd := exec.Command("pkill", "-f", "quickshell")
	if err := cmd.Run(); err != nil {
		// Exit code 1 means no processes were found, which is fine
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil
		}
		return fmt.Errorf("failed to kill QuickShell: %w", err)
	}
	return nil
}
