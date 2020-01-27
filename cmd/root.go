package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{}
)

// RootCommand ...
func RootCommand() *cobra.Command {
	rootCmd.AddCommand(&serveCmd)

	return rootCmd
}
