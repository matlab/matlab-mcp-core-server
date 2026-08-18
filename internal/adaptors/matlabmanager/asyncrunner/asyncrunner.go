// Copyright 2026 The MathWorks, Inc.

package asyncrunner

import (
	"context"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/entities"
)

type Task struct {
	done chan struct{}
}

func Run(parent context.Context, logger entities.Logger, timeout time.Duration, label string, op func(context.Context) error) *Task {
	task := &Task{done: make(chan struct{})}
	go func() {
		defer close(task.done)
		ctx, cancel := context.WithTimeout(parent, timeout)
		defer cancel()
		if err := op(ctx); err != nil {
			logger.WithError(err).Warn(label)
		}
	}()
	return task
}

func (t *Task) Wait(ctx context.Context) bool {
	select {
	case <-t.done:
		return true
	case <-ctx.Done():
		return false
	}
}
