package search

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/certsio/certsio/pkg/config"

	"github.com/cenkalti/backoff"

	"github.com/certsio/certsio/pkg/certificate"
)

const _baseURL = "https://certs-io1.p.rapidapi.com/certificates"

type Config struct {
	maxPages uint64
	baseURL  string
}

// Client is the API search client for certs.io
type Client struct {
	transport *Transport
	config    *Config
}

// NewClient returns a new API client.
func NewClient(cfg config.Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = _baseURL
	}
	retryBackoff := backoff.NewExponentialBackOff()
	return &Client{
		transport: NewTransport(TransportConfig{
			ApiKey: cfg.APIKey,
			RetryBackoff: func(i int) time.Duration {
				if i == 1 {
					retryBackoff.Reset()
				}
				return retryBackoff.NextBackOff()
			},
		}),
		config: &Config{
			maxPages: 0,
			baseURL:  cfg.BaseURL,
		},
	}
}

// WithBaseURL sets the base URL for API requests.
func (c *Client) WithBaseURL(baseURL string) *Client {
	c.config.baseURL = baseURL
	return c
}

// WithMaxPages sets the maximum number of pages to return.
func (c *Client) WithMaxPages(maxPages uint64) *Client {
	c.config.maxPages = maxPages
	return c
}

// WithTimeout sets the timeout for API requests.
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.transport.httpClient.Timeout = timeout
	return c
}

// WithRetries sets the maximum number of retries for API requests.
func (c *Client) WithRetries(retries int) *Client {
	c.transport.maxRetries = retries
	return c
}

// Query is the search query for the certs.io API.
type Query struct {
	Field Field  `json:"field"`
	Value string `json:"term"`
	Page  uint64 `json:"page,omitempty"`
}

// Response is the response from the certs.io API.
type Response struct {
	Total        uint64                    `json:"total_certificates"`
	Pages        uint64                    `json:"total_pages"`
	CurrentPage  uint64                    `json:"page"`
	Certificates []certificate.Certificate `json:"certificates"`
}

// Search performs a search and wraps the searchWithCtx method.
func (c *Client) Search(ctx context.Context, query *Query) ([]certificate.Certificate, error) {
	return c.searchWithCtx(ctx, query)
}

// searchWithCtx performs a search with a context
func (c *Client) searchWithCtx(ctx context.Context, query *Query) ([]certificate.Certificate, error) {
	var results []certificate.Certificate
	resultChan := make(chan []certificate.Certificate)
	errChan := make(chan error) // Create a channel to receive errors

	go func() {
		defer close(resultChan)
		err := c.search(ctx, query, resultChan)
		if err != nil {
			errChan <- err // Send the error to the error channel
			return
		}
	}()

	for certificates := range resultChan {
		results = append(results, certificates...)
	}

	// Check if there is an error in the error channel
	select {
	case err := <-errChan:
		return nil, err // Return the error
	default:
		return results, nil // No error, return the results
	}
}

// StreamSearchResults performs a search with a context
func (c *Client) StreamSearchResults(ctx context.Context, query *Query, resultChan chan<- []certificate.Certificate) error {
	return c.search(ctx, query, resultChan)
}

// search performs a search with a context and streams the certificates to a results channel.
func (c *Client) search(ctx context.Context, query *Query, resultChan chan<- []certificate.Certificate) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			body, err := json.Marshal(query)
			if err != nil {
				return fmt.Errorf("api client: %w", err)
			}

			resp, err := c.transport.POST(c.config.baseURL, body)
			if err != nil {
				return fmt.Errorf("api client: %w", err)
			}

			var result Response
			// decode the response
			if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
				resp.Body.Close()
				return fmt.Errorf("api client: %w", err)
			}
			resp.Body.Close()

			// send the results to the result channel
			resultChan <- result.Certificates

			// check if we've reached the last page or the maximum number of pages.
			if result.Pages == result.CurrentPage || (c.config.maxPages > 0 && result.CurrentPage == c.config.maxPages-1) {
				return nil
			}
			query.Page++
		}
	}
}
