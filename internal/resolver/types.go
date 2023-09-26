package resolver

import "github.com/certsio/certsio/pkg/certificate"

// Types of data result can return
type (
	// ResultType is the type of result returned
	ResultType int

	// Result is the resolution of a hostname tied to a certificate.
	Result struct {
		// Type is the type of result returned
		Type ResultType
		// Task is the resolution task performed.
		Task HostEntry
		// IPs is the resolved IP address for the host.
		IPs []string
		// Error is the error that occurred during resolution.
		Error error
	}

	// HostEntry defines a host with the source
	HostEntry struct {
		// Hostname is an SSL name found in a certificate.
		Host string
		// Source is the certificate that contains the hostname.
		Source certificate.Certificate
	}
)
