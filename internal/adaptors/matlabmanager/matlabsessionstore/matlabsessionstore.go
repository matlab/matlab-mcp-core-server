// Copyright 2025-2026 The MathWorks, Inc.

package matlabsessionstore

import (
	"context"
	"fmt"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"golang.org/x/sync/errgroup"
)

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type MATLABSessionClientWithCleanup interface {
	entities.MATLABSessionClient
	StopSession(ctx context.Context, sessionLogger entities.Logger) error
}

type LifecycleSignaler interface {
	AddShutdownFunction(shutdownFcn func() error)
}

type Store struct {
	l       *sync.RWMutex
	next    entities.SessionID
	clients map[entities.SessionID]MATLABSessionClientWithCleanup
}

func New(
	loggerFactory LoggerFactory,
	lifecycleSignaler LifecycleSignaler,
) *Store {
	store := &Store{
		l:       new(sync.RWMutex),
		next:    1,
		clients: map[entities.SessionID]MATLABSessionClientWithCleanup{},
	}

	lifecycleSignaler.AddShutdownFunction(func() error {
		store.l.Lock()
		defer store.l.Unlock()

		logger, err := loggerFactory.GetGlobalLogger()
		if err != nil {
			return err
		}

		wg := new(errgroup.Group)

		for sessionID, client := range store.clients {
			wg.Go(func() error {
				err := client.StopSession(context.Background(), logger)
				if err != nil {
					return fmt.Errorf("error stopping session %v: %w", sessionID, err)
				}
				return nil
			})
		}

		return wg.Wait()
	})

	return store
}

func (s *Store) Add(client MATLABSessionClientWithCleanup) entities.SessionID {
	s.l.Lock()
	defer s.l.Unlock()

	sessionID := s.next
	s.clients[sessionID] = client
	s.next++
	return entities.SessionID(sessionID)
}

func (s *Store) Get(sessionID entities.SessionID) (MATLABSessionClientWithCleanup, error) {
	s.l.RLock()
	defer s.l.RUnlock()

	client, exists := s.clients[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %v", sessionID)
	}

	return client, nil
}

func (s *Store) Remove(sessionID entities.SessionID) {
	s.l.Lock()
	defer s.l.Unlock()

	delete(s.clients, sessionID)
}
