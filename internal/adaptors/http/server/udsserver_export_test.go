// Copyright 2025-2026 The MathWorks, Inc.

package server

import "net/http"

func NewUDSServer(osLayer OSLayer, handlers map[string]http.HandlerFunc) *udsServer {
	return newUDSServer(osLayer, handlers)
}
