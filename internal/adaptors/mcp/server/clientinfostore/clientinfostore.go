// Copyright 2026 The MathWorks, Inc.

package clientinfostore

import (
	"sync"

	"github.com/matlab/matlab-mcp-server/internal/entities"
)

type ClientInfoStore struct {
	mu         sync.RWMutex
	clientInfo entities.MCPClientInfo
}

func New() *ClientInfoStore {
	return &ClientInfoStore{}
}

func (s *ClientInfoStore) SetClientInfo(info entities.MCPClientInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clientInfo = info
}

func (s *ClientInfoStore) GetClientInfo() entities.MCPClientInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clientInfo
}
