package search

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultMaxRetries = 3
)

// TransportConfig is the configuration for the HTTP transport.
type TransportConfig struct {
	HTTPTransport http.RoundTripper
	RetryBackoff  func(attempt int) time.Duration
	MaxRetries    int
	ApiKey        string
}

// Transport is the HTTP transport for the certs.io API.
type Transport struct {
	httpClient *http.Client

	userAgent    string
	apiKey       string
	maxRetries   int
	retryBackoff func(attempt int) time.Duration
}

// NewTransport returns a new HTTP transport.
func NewTransport(cfg TransportConfig) *Transport {
	// set default max retries
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = defaultMaxRetries
	}

	client := &Transport{
		httpClient: &http.Client{
			Transport: http.DefaultTransport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// do not follow redirects.
				return http.ErrUseLastResponse
			},
			Timeout: 15 * time.Second,
		},
		maxRetries:   cfg.MaxRetries,
		retryBackoff: cfg.RetryBackoff,
		userAgent:    "certsio-client-go",
		apiKey:       cfg.ApiKey,
	}

	// set custom transport if provided
	if cfg.HTTPTransport != nil {
		client.httpClient.Transport = cfg.HTTPTransport
	}

	return client
}

// POST performs a POST request to the certs.io API with retries.
func (t *Transport) POST(url string, body []byte) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)
	for i := 0; i < t.maxRetries; i++ {
		resp, err = t.post(url, body)
		if err != nil {
			continue
		}

		// if the status code is 200, break out of the loop.
		if resp.StatusCode == http.StatusOK {
			break
		}

		if resp.StatusCode == http.StatusBadRequest {
			err = fmt.Errorf("bad search request")
			break
		}

		if resp.StatusCode == http.StatusUnauthorized {
			err = fmt.Errorf("bad api key")
			break
		}
		// if the status code is 429, retry with backoff.
		if resp.StatusCode == http.StatusTooManyRequests {
			err = fmt.Errorf("rate limit exceeded")
		} else {
			err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		// Delay the retry if a backoff function is configured
		if t.retryBackoff != nil {
			time.Sleep(t.retryBackoff(i + 1))
		}
	}

	return resp, err
}

// Post performs a POST request to the certs.io API.
func (t *Transport) post(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", t.userAgent)
	req.Header.Set("X-RapidAPI-Key", t.apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
