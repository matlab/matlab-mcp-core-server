// Copyright 2025-2026 The MathWorks, Inc.

package server

import (
	"context"
	"net/http"
)

type OSLayer interface {
	RemoveAll(name string) error
}

type HttpServer interface {
	Serve(socketPath string) error
	Shutdown(ctx context.Context) error
}

type Factory struct {
	osLayer OSLayer
}

func NewFactory(osLayer OSLayer) *Factory {
	return &Factory{
		osLayer: osLayer,
	}
}

func (f *Factory) NewServerOverUDS(handlers map[string]http.HandlerFunc) (HttpServer, error) {
	return newUDSServer(f.osLayer, handlers), nil
}
