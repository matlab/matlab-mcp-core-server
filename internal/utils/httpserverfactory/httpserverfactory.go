// Copyright 2025-2026 The MathWorks, Inc.

package httpserverfactory

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/socket"
)

const defaultReadHeaderTimeout = 10 * time.Second

type HttpServer interface {
	Serve(socketPath string) error
	Shutdown(ctx context.Context) error
}

type OSLayer interface {
	RemoveAll(name string) error
}

type HTTPServerFactory struct {
	osLayer OSLayer
}

func New(osLayer OSLayer) *HTTPServerFactory {
	return &HTTPServerFactory{
		osLayer: osLayer,
	}
}

func (f *HTTPServerFactory) NewServerOverUDS(handlers map[string]http.HandlerFunc) (HttpServer, error) {
	mux := http.NewServeMux()
	for pattern, handler := range handlers {
		mux.HandleFunc(pattern, handler)
	}

	return &udsServer{
		httpServer: &http.Server{
			Handler:           mux,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
		},
		osLayer: f.osLayer,

		lock: new(sync.Mutex),
	}, nil
}

type udsServer struct {
	httpServer *http.Server
	osLayer    OSLayer
	socketPath string

	lock *sync.Mutex
}

func (s *udsServer) Serve(socketPath string) error {
	// Socket path max length is 108 characters, but for safety using 105
	if len(socketPath) > 105 {
		return socket.ErrSocketPathTooLong
	}

	if err := s.osLayer.RemoveAll(socketPath); err != nil {
		return err
	}

	s.setSocketPath(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *udsServer) Shutdown(ctx context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	if s.socketPath == "" {
		return nil
	}

	if err := s.osLayer.RemoveAll(s.socketPath); err != nil {
		return err
	}

	return nil
}

func (s *udsServer) setSocketPath(socketPath string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.socketPath = socketPath
}
