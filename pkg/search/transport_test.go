package search

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type SearchTransportTestSuite struct {
	suite.Suite
}

// TestNewTransport tests the NewTransport function.
func (s *SearchTransportTestSuite) TestNewTransport() {
	// test default transport
	s.Run("default", func() {
		t := NewTransport(TransportConfig{})
		s.Equal(http.DefaultTransport, t.httpClient.Transport)
		s.Equal("certsio/1.0.0", t.userAgent)
		s.Equal(defaultMaxRetries, t.maxRetries)
		s.Nil(t.retryBackoff)
		s.NotNil(t)
	})
	// test custom transport
	s.Run("custom", func() {
		config := TransportConfig{
			HTTPTransport: &mockTransport{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) { return &http.Response{Status: "MOCK"}, nil },
			},
			RetryBackoff: func(i int) time.Duration {
				return time.Duration(i) * time.Second
			},
			MaxRetries: 5,
		}
		t := NewTransport(config)
		s.Equal(config.HTTPTransport, t.httpClient.Transport)
	})
}

// TestTransportPOST tests the POST function with different statuses.
func (s *SearchTransportTestSuite) TestTransportPOST() {
	// Test Status OK on first try:
	s.Run("Status OK", func() {
		var (
			i       int
			numReqs = 1
		)
		config := TransportConfig{
			HTTPTransport: &mockTransport{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					i++
					return &http.Response{StatusCode: http.StatusOK}, nil
				},
			},
		}
		t := NewTransport(config)
		resp, err := t.POST("http://localhost", []byte("body"))
		s.Nil(err)
		s.Equal(http.StatusOK, resp.StatusCode)
		s.Equal(numReqs, i)
	})
	// Test retry from 429 -> 200:
	s.Run("Rate-limited 429->200", func() {
		var (
			i       int
			numReqs = 5
		)
		config := TransportConfig{
			HTTPTransport: &mockTransport{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					i++
					if i == numReqs {
						return &http.Response{StatusCode: http.StatusOK}, nil
					}
					return &http.Response{StatusCode: http.StatusTooManyRequests}, nil
				},
			},
			MaxRetries: 5,
		}
		t := NewTransport(config)
		resp, err := t.POST("http://localhost", []byte("body"))
		s.Nil(err)
		s.Equal(http.StatusOK, resp.StatusCode)
		s.Equal(numReqs, i)
	})
	// Test bad request:
	s.Run("Bad Request", func() {
		var (
			i       int
			numReqs = 1
		)
		config := TransportConfig{
			HTTPTransport: &mockTransport{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					i++
					return &http.Response{StatusCode: http.StatusBadRequest}, nil
				},
			},
		}
		t := NewTransport(config)
		resp, err := t.POST("http://localhost", []byte("body"))
		s.NotNil(err)
		s.Equal(http.StatusBadRequest, resp.StatusCode)
		s.Equal(numReqs, i)
	})

	s.Run("Too Many Retries", func() {
		var (
			i          int
			numReqs    = 5
			maxRetries = 4
		)
		config := TransportConfig{
			HTTPTransport: &mockTransport{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					i++
					if i == numReqs {
						return &http.Response{StatusCode: http.StatusOK}, nil
					}
					return &http.Response{StatusCode: http.StatusTooManyRequests}, nil
				},
			},
			MaxRetries: maxRetries,
		}
		t := NewTransport(config)
		_, err := t.POST("http://localhost", []byte("body"))
		s.NotNil(err)
		s.Equal(maxRetries, i)
	})

	s.Run("500", func() {
		var i int
		config := TransportConfig{
			HTTPTransport: &mockTransport{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					i++

					return &http.Response{StatusCode: http.StatusInternalServerError}, nil
				},
			},
		}
		t := NewTransport(config)
		_, err := t.POST("http://localhost", []byte("body"))
		s.NotNil(err)
	})
}

// TestRunTransportTestSuite runs the test suite.
func TestRunTransportTestSuite(t *testing.T) {
	suite.Run(t, new(SearchTransportTestSuite))
}

type mockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}
