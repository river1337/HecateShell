package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"hecate-shell/internal/config"
	"hecate-shell/internal/shell"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the Hecate shell",
	Long: `Start the HecateShell using QuickShell.

This will launch QuickShell with the HecateShell configuration.`,
	RunE: runShell,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("reload", "r", false, "Kill existing instance before starting")
	runCmd.Flags().BoolP("debug", "d", false, "Run in foreground (without --daemonize)")
}

func runShell(cmd *cobra.Command, args []string) error {
	reload, _ := cmd.Flags().GetBool("reload")
	debug, _ := cmd.Flags().GetBool("debug")

	// Check if installed
	if !config.IsInstalled() {
		return fmt.Errorf("HecateShell is not installed. Run 'hecate install' first")
	}

	// Kill existing instance if reload flag is set
	if reload {
		fmt.Println("Stopping existing QuickShell instances...")
		if err := shell.KillQuickShell(); err != nil {
			fmt.Printf("Warning: %v\n", err)
		}
	}

	// Get config path
	shellPath, err := config.GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get shell config path: %w", err)
	}

	// Start QuickShell
	var quickshellCmd *exec.Cmd
	if debug {
		fmt.Println("Starting HecateShell in debug mode (foreground)...")
		quickshellCmd = exec.Command("quickshell", "--no-duplicate", "-c", shellPath)
	} else {
		fmt.Println("Starting HecateShell...")
		quickshellCmd = exec.Command("quickshell", "--no-duplicate", "--daemonize", "-c", shellPath)
	}
	quickshellCmd.Stdout = os.Stdout
	quickshellCmd.Stderr = os.Stderr

	// Handle signals for clean shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the process
	if err := quickshellCmd.Start(); err != nil {
		return fmt.Errorf("failed to start QuickShell: %w", err)
	}

	fmt.Printf("HecateShell started (PID: %d)\n", quickshellCmd.Process.Pid)

	// Wait for signal or process to exit
	go func() {
		<-sigChan
		fmt.Println("\nShutting down HecateShell...")
		quickshellCmd.Process.Signal(syscall.SIGTERM)
	}()

	// Wait for process to finish
	if err := quickshellCmd.Wait(); err != nil {
		return fmt.Errorf("QuickShell exited with error: %w", err)
	}

	return nil
}
