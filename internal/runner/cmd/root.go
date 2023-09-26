package cmd

import (
	"github.com/certsio/certsio/internal/runner/cmd/searchcmd"
	"github.com/spf13/cobra"
)

type Config struct {
	searchOpts searchcmd.Options
	configFile string
}

var rootConfig = &Config{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "certsio",
	Short: `command line client for certs.io`,
}

// init initializes the root command.
func init() {
	rootCmd.PersistentFlags().StringVarP(&rootConfig.configFile, "config", "c", "", "config file (default is $HOME/.certsio.toml)")
	rootCmd.PersistentFlags().StringVarP(&rootConfig.searchOpts.OutputFile, "output", "o", "", "output file (default is stdout)")

	rootCmd.AddCommand(searchcmd.New(&rootConfig.searchOpts))
	rootCmd.AddCommand(versionCmd)
	//rootCmd.AddCommand(resolvecmd.New().Command())
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
