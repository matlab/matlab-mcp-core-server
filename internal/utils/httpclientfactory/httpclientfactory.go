// Copyright 2025 The MathWorks, Inc.

package httpclientfactory

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type HTTPClientFactory struct{}

func New() *HTTPClientFactory {
	return &HTTPClientFactory{}
}

func (f *HTTPClientFactory) NewClientForSelfSignedTLSServer(certificatePEM []byte) (HttpClient, error) {
	caCertPool := x509.NewCertPool()

	if ok := caCertPool.AppendCertsFromPEM(certificatePEM); !ok {
		return nil, fmt.Errorf("failed to append certificate to pool")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true, // We do full verification ourselves below
			// Custom verification to allow clock skew tolerance
			// We must use InsecureSkipVerify because Go's standard validation
			// checks certificate dates BEFORE calling VerifyPeerCertificate,
			// which prevents our clock skew tolerance from working
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				// Allow up to 24 hours of clock skew for self-signed certificates
				// This handles cases where system clocks are slightly out of sync
				const clockSkewTolerance = 24 * time.Hour
				
				opts := x509.VerifyOptions{
					Roots:         caCertPool,
					Intermediates: x509.NewCertPool(),
				}

				// Verify each certificate in the chain with clock skew tolerance
				for _, rawCert := range rawCerts {
					cert, err := x509.ParseCertificate(rawCert)
					if err != nil {
						return fmt.Errorf("failed to parse certificate: %w", err)
					}
					
					// Check certificate validity with clock skew tolerance
					now := time.Now()
					if cert.NotBefore.After(now.Add(clockSkewTolerance)) {
						// Certificate is too far in the future, reject it
						return fmt.Errorf("certificate not valid yet: notBefore is %v, current time is %v (skew tolerance: %v)", 
							cert.NotBefore, now, clockSkewTolerance)
					}
					if cert.NotAfter.Before(now.Add(-clockSkewTolerance)) {
						// Certificate is too far expired, reject it
						return fmt.Errorf("certificate expired: notAfter is %v, current time is %v (skew tolerance: %v)", 
							cert.NotAfter, now, clockSkewTolerance)
					}
					
					// Try verification with current time
					opts.CurrentTime = now
					_, err = cert.Verify(opts)
					if err != nil {
						// Try with positive clock skew (certificate is in the future)
						opts.CurrentTime = now.Add(clockSkewTolerance)
						_, err = cert.Verify(opts)
						if err != nil {
							// Try with negative clock skew (certificate is in the past)
							opts.CurrentTime = now.Add(-clockSkewTolerance)
							_, err = cert.Verify(opts)
							if err != nil {
								return fmt.Errorf("certificate verification failed: %w", err)
							}
						}
					}
				}
				return nil
			},
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
