// Copyright 2026 The MathWorks, Inc.

package asyncrunner_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/asyncrunner"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testLabel = "async task failed"

func TestRun_HappyPath(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		ctx := t.Context()
		mockLogger := testutils.NewInspectableLogger()
		ran := false

		// Act
		task := asyncrunner.Run(ctx, mockLogger, time.Second, testLabel, func(context.Context) error {
			ran = true
			return nil
		})
		completed := task.Wait(ctx)

		// Assert
		assert.True(t, completed, "Wait should report completion")
		assert.True(t, ran, "op should have run")
		assert.Empty(t, mockLogger.WarnLogs(), "a successful op should not log")
	})
}

func TestRun_OpError_LogsWarnWithLabel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		ctx := t.Context()
		mockLogger := testutils.NewInspectableLogger()

		// Act
		task := asyncrunner.Run(ctx, mockLogger, time.Second, testLabel, func(context.Context) error {
			return assert.AnError
		})
		completed := task.Wait(ctx)

		// Assert
		assert.True(t, completed)
		fields, logged := mockLogger.WarnLogs()[testLabel]
		require.True(t, logged, "op error should be logged at Warn under the label")
		assert.Equal(t, assert.AnError, fields["error"], "the op error should be attached to the log")
	})
}

func TestRun_Timeout_CancelsOpAndLogs(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		ctx := t.Context()
		mockLogger := testutils.NewInspectableLogger()

		// Act: the op observes its context and returns the timeout error once the budget elapses.
		task := asyncrunner.Run(ctx, mockLogger, time.Second, testLabel, func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		})
		completed := task.Wait(ctx)

		// Assert
		assert.True(t, completed)
		fields, logged := mockLogger.WarnLogs()[testLabel]
		require.True(t, logged, "a timed-out op should be logged at Warn")
		assert.ErrorIs(t, fields["error"].(error), context.DeadlineExceeded)
	})
}

func TestTask_Wait_ReturnsFalseOnContextCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		ctx := t.Context()
		mockLogger := testutils.NewInspectableLogger()
		release := make(chan struct{})
		task := asyncrunner.Run(ctx, mockLogger, time.Minute, testLabel, func(context.Context) error {
			<-release
			return nil
		})

		waitCtx, cancel := context.WithCancel(ctx)

		// Act: cancel the wait context while the op is still in flight.
		cancel()
		completed := task.Wait(waitCtx)

		// Assert
		assert.False(t, completed, "Wait should report non-completion when its context is cancelled")

		// Drain the still-running task before the bubble ends.
		close(release)
		synctest.Wait()
	})
}
