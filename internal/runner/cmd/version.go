package cmd

import "github.com/spf13/cobra"

const (
	// Version is the current version of the certsio CLI.
	Version = "0.1.0"
)

// versionCmd prints out the version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: `Print the current version`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("certsio version %s\n", Version)
	},
}
