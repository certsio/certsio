package searchcmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/certsio/certsio/pkg/config"

	"github.com/certsio/certsio/pkg/certificate"
	"github.com/certsio/certsio/pkg/output"
	"github.com/certsio/certsio/pkg/search"
	"github.com/spf13/cobra"
)

// createSearchCommand creates a new search command for each searchable field.
func (s *Search) createSearchCommand(field search.Field) *cobra.Command {
	return &cobra.Command{
		Use:   field.String(),
		Short: fmt.Sprintf("Search for certificates by %s", field.String()),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Printf("usage: certsio search %s \"<value>\"\n", field.String())
				os.Exit(1)
			}
			cfgFile, err := cmd.Root().Flags().GetString("config")
			if err != nil {
				log.Fatalf("couldn't get config file: %v", err)
			}

			cfg, err := config.Get(cfgFile)
			if err != nil {
				if err := config.Write(cfgFile); err != nil {
					log.Fatalf("couldn't create config file: %v", err)
				}
				log.Fatalf("please update the config file with your API key")
			}

			s.runSearch(cfg, field, args)
		},
	}
}

// runSearch runs the search command.
func (s *Search) runSearch(cfg config.Config, field search.Field, args []string) {
	var (
		wg       sync.WaitGroup
		writerWg sync.WaitGroup
		file     *os.File
		err      error
	)

	// Create a new search client.
	if cfg.APIKey == "" || cfg.APIKey == "CHANGE_ME" {
		log.Fatalf("Please update the config file with your API key")
	}

	client := search.NewClient(cfg)
	client.WithMaxPages(s.opts.maxPages)

	resultChan := make(chan []certificate.Certificate)

	file = os.Stdout
	if s.opts.OutputFile != "" {
		file, err = os.OpenFile(s.opts.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			log.Fatalf("couldn't create output file: %v", err)
		}
	}
	defer file.Close()
	writer := output.NewWriter(file)
	writerWg.Add(1)
	go func() {
		defer writerWg.Done()
		for results := range resultChan {
			for _, res := range results {
				if err := writer.Write(res); err != nil {
					// TODO: handle error
					continue
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := client.StreamSearchResults(context.Background(), &search.Query{
			Field: field,
			Value: args[0],
			Page:  0,
		}, resultChan)
		// check errors
		if err != nil {
			// TODO: handle error
			return
		}
	}()

	wg.Wait()
	close(resultChan)
	writerWg.Wait()
}
