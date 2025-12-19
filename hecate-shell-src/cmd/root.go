package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hecate",
	Short: "Hecate Shell - Modern Wayland shell for Niri",
	Long: `Hecate is a beautiful, customizable shell built with QuickShell.

Features:
  - Matugen color integration
  - Hot-reloading themes
  - Fully customizable`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be added here
}
