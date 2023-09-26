package output

import (
	"io"

	"github.com/certsio/certsio/pkg/certificate"
	jsoniter "github.com/json-iterator/go"
)

// CertificateWriter writes certificates to a destination.
type CertificateWriter interface {
	Write(cert certificate.Certificate) error
}

// writer is a writer for writing certificates.
type writer struct {
	writer io.Writer
}

// NewWriter creates a new writer for writing certificates.
func NewWriter(w io.Writer) *writer {
	return &writer{writer: w}
}

// Write writes the event to file and/or screen.
func (w *writer) Write(cert certificate.Certificate) error {
	return writeCertificate(w.writer, cert)
}

// Write writes the event to file and/or screen.
func writeCertificate(w io.Writer, cert certificate.Certificate) error {
	// encode the certificate
	encoder := jsoniter.NewEncoder(w)
	if err := encoder.Encode(&cert); err != nil {
		return err
	}

	return nil
}
