package certresolve

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/certsio/certsio/internal/resolver"
	"github.com/certsio/certsio/pkg/certificate"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

// Config defines the certificate resolver configuration.
type Config struct {
	WorkerCount int
}

// Resolver drives the resolution of certificate names.
type Resolver struct {
	config *Config
	pool   *resolver.Pool
}

// New creates a new certificate resolver.
func New(c *Config) (*Resolver, error) {
	client, err := resolver.New()
	if err != nil {
		return nil, fmt.Errorf("certresolve: %w", err)
	}

	return &Resolver{
		config: c,
		pool:   client.NewPool(c.WorkerCount),
	}, nil
}

func (r *Resolver) Start(in io.Reader) {
	// process the output
	var outWg sync.WaitGroup
	outWg.Add(1)
	go func() {
		defer outWg.Done()
		for result := range r.pool.Results {
			r.processResult(result)
		}
	}()

	reader := bufio.NewScanner(in)
	for reader.Scan() {
		var cert certificate.Certificate
		if err := jsoniter.Unmarshal(reader.Bytes(), &cert); err != nil {
			continue
		}
		r.resolveCertificateNames(cert)
	}

	close(r.pool.Tasks)
	outWg.Wait()
}

// resolveCertificateNames resolves the names of a certificate.
func (r *Resolver) resolveCertificateNames(cert certificate.Certificate) {
	for _, name := range cert.Names {
		r.pool.Tasks <- resolver.HostEntry{Host: name, Source: cert}
	}
}

func (r *Resolver) processResult(result resolver.Result) {
	switch result.Type {
	case resolver.Alive:
		var skip bool
		certIP := strings.Split(result.Task.Source.Server, ":")[0]
		// if certIP not in result.IPs, then print
		for _, ip := range result.IPs {
			if ip == certIP {
				skip = true
			}
		}
		if !skip {
			fields := logrus.Fields{}
			fields["host"] = result.Task.Host
			fields["resolved_ips"] = result.IPs
			fields["source_ip"] = certIP
			logrus.WithFields(fields).Info("Possible Origin Bypass")
		}
	default:
		fields := logrus.Fields{}
		fields["host"] = result.Task.Host
		fields["source"] = result.Task.Source.Server
		logrus.WithFields(fields).Warn("Possible Internal Host")
	}
}
