package certificate

import (
	"encoding/json"
	"time"
)

// Certificate hold the extracted information from an SSL certificate
type Certificate struct {
	Timestamp string `json:"@timestamp"` // @timestamp is needed for ES data-streams (https://www.elastic.co/guide/en/elasticsearch/reference/current/data-streams.html)
	// Server is the host or IP address of the server
	Server string `json:"server"`
	// Expired specifies whether the certificate has expired
	Expired bool `json:"expired"`
	// SelfSigned returns true if the certificate is self-signed
	SelfSigned bool `json:"self_signed"`
	// Revoked returns true if the certificate is revoked
	Revoked bool `json:"revoked"`
	// NotBefore is the not-before time for certificate
	NotBefore time.Time `json:"not_before"`
	// NotAfter is the not-after time for certificate
	NotAfter time.Time `json:"not_after"`
	// Names are all names identified in the certificate
	Names []string `json:"ssl_names"`
	// Org is the organization for the certificate
	SubjectOrg []string `json:"subject_org"`
	// Serial is the certificate serial number
	Serial string `json:"serial"`
	// IssuerNames is the issuer names for cert
	IssuerNames []string `json:"issuer_names"`
	// IssuerOrg is the organization for cert issuer
	IssuerOrg []string `json:"issuer_org"`
	// Emails is a list of Emails for the certificate
	Emails []string `json:"emails"`
	// FingerprintHash is the hashes for certificate
	FingerprintSha256Hash string `json:"fingerprint_sha256"`
	// ParentDomains is the list of parent domains for ssl_names in the certificate
	ParentDomains []string `json:"parent_domains"`
}

// Unmarshal parses a certificate from JSON-encoded data.
func (c *Certificate) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}

// Marshal returns the JSON encoding of Certificate.
func (c *Certificate) Marshal() ([]byte, error) {
	return json.Marshal(c)
}
