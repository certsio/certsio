package resolvecmd

import (
	"log"
	"os"

	"github.com/certsio/certsio/pkg/certresolve"
	"github.com/spf13/cobra"
)

type Config struct {
	inputFile string
}
type Command struct {
	cmd    *cobra.Command
	config *Config
}

func New() *Command {
	c := &Command{
		config: &Config{},
	}
	c.cmd = &cobra.Command{
		Use:   "resolve",
		Short: "Resolve the ssl_names within a certificate.",
		Long:  "Resolve the ssl_names within a certificate. Find potential origin bypasses or interesting certificates.",
		Run:   c.Run,
	}
	c.cmd.PersistentFlags().StringVarP(&c.config.inputFile, "input", "i", "", "input file containing TLS certificates")
	return c
}

// Command returns the cobra command.
func (c *Command) Command() *cobra.Command {
	return c.cmd
}

// Run executes the command.
func (c *Command) Run(cmd *cobra.Command, args []string) {
	var (
		file *os.File
		err  error
	)
	if c.config.inputFile == "" {
		log.Fatal("no input file specified.")
	}
	file, err = os.Open(c.config.inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	resolver, err := certresolve.New(&certresolve.Config{
		WorkerCount: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	resolver.Start(file)
}
