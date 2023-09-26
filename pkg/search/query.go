package search

// Field is a search field.
type Field string

const (
	// ByDomain searches for certificates by domain.
	ByDomain Field = "domain"
	// ByIP searches for certificates by IP address:port.
	ByServer Field = "server"
	// ByFingerprint searches for certificates by SHA256 fingerprint.
	ByFingerprint Field = "fingerprint_sha256"
	// ByEmail searches for certificates by email address.
	ByEmails Field = "emails"
	// ByIssuer searches for certificates by organization.
	ByOrg Field = "org"
	// BySerial searches for certificates by serial number.
	BySerial Field = "serial"
	// ByCertNames searches for certificates by common name or subject alternative name.
	ByCertNames Field = "ssl_names"
)

// String returns the string representation of a Field.
func (f Field) String() string {
	return string(f)
}
