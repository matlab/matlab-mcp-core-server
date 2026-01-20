// Copyright 2025-2026 The MathWorks, Inc.

package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
)

type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewClientForSelfSignedTLSServer(certificatePEM []byte) (HttpClient, error) {
	caCertPool := x509.NewCertPool()

	if ok := caCertPool.AppendCertsFromPEM(certificatePEM); !ok {
		return nil, fmt.Errorf("failed to append certificate to pool")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    caCertPool,
		},
	}

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	return &http.Client{
		Transport: transport,
		Jar:       jar,
	}, nil
}

func (f *Factory) NewClientOverUDS(socketPath string) HttpClient {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", socketPath)
		},
	}

	return &http.Client{
		Transport: transport,
	}
}
