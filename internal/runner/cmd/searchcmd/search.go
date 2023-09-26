package searchcmd

import (
	"github.com/certsio/certsio/pkg/search"
	"github.com/spf13/cobra"
)

type Options struct {
	maxPages   uint64
	OutputFile string
}

type Search struct {
	opts *Options
	cmd  *cobra.Command
}

// New instantiates the search command
func New(o *Options) *cobra.Command {
	// search command carries out all searches
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: `Search for certificates`,
	}

	searcher := &Search{
		opts: o,
		cmd:  searchCmd,
	}
	// add flags
	searcher.cmd.PersistentFlags().Uint64VarP(&searcher.opts.maxPages, "max-pages", "m", 0, "maximum number of pages to return (0 for all pages)")

	// add additional subcommands
	searcher.cmd.AddCommand(searcher.createSearchCommand(search.ByOrg))
	searcher.cmd.AddCommand(searcher.createSearchCommand(search.ByDomain))
	searcher.cmd.AddCommand(searcher.createSearchCommand(search.ByFingerprint))
	searcher.cmd.AddCommand(searcher.createSearchCommand(search.BySerial))
	searcher.cmd.AddCommand(searcher.createSearchCommand(search.ByEmails))
	searcher.cmd.AddCommand(searcher.createSearchCommand(search.ByCertNames))
	searcher.cmd.AddCommand(searcher.createSearchCommand(search.ByServer))

	return searcher.cmd
}
