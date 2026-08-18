// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/asyncrunner"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	entitiesmocks "github.com/matlab/matlab-mcp-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// blockedGreetingTask returns a greeting task whose op parks until release is closed, so a test can
// observe the gate blocking. The timeout is long enough never to fire during a test.
func blockedGreetingTask(logger entities.Logger, release <-chan struct{}) *asyncrunner.Task {
	return asyncrunner.Run(context.Background(), logger, time.Hour, "greeting", func(context.Context) error {
		<-release
		return nil
	})
}

func TestGreetingSessionClient_Eval_BlocksUntilGreetingDone(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
		defer mockSessionClient.AssertExpectations(t)

		ctx := t.Context()
		evalRequest := entities.EvalRequest{Code: "1 + 1"}
		evalResponse := entities.EvalResponse{ConsoleOutput: "ans = 2"}

		greetingReleased := make(chan struct{})
		greetingTask := blockedGreetingTask(mockLogger, greetingReleased)

		mockSessionClient.EXPECT().
			Eval(ctx, mockLogger.AsMockArg(), evalRequest).
			Return(evalResponse, nil).
			Once()

		client := matlabmanager.NewGreetingSessionClient(mockSessionClient, greetingTask)

		// Act: call Eval in a goroutine so we can prove it is parked on the gate.
		type result struct {
			response entities.EvalResponse
			err      error
		}
		resultChan := make(chan result, 1)
		go func() {
			response, err := client.Eval(ctx, mockLogger, evalRequest)
			resultChan <- result{response, err}
		}()

		// Assert: with the greeting still in flight, Eval must not have returned.
		synctest.Wait()
		select {
		case <-resultChan:
			t.Fatal("Eval returned before the greeting completed")
		default:
		}

		// Release the greeting; the gate should now unblock and delegate to the underlying client.
		close(greetingReleased)
		got := <-resultChan
		require.NoError(t, got.err)
		assert.Equal(t, evalResponse, got.response)
	})
}

func TestGreetingSessionClient_FEval_BlocksUntilGreetingDone(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
		defer mockSessionClient.AssertExpectations(t)

		ctx := t.Context()
		fevalRequest := entities.FEvalRequest{Function: "plus", Arguments: []string{"1", "1"}, NumOutputs: 1}
		fevalResponse := entities.FEvalResponse{Outputs: []any{2.0}}

		greetingReleased := make(chan struct{})
		greetingTask := blockedGreetingTask(mockLogger, greetingReleased)

		mockSessionClient.EXPECT().
			FEval(ctx, mockLogger.AsMockArg(), fevalRequest).
			Return(fevalResponse, nil).
			Once()

		client := matlabmanager.NewGreetingSessionClient(mockSessionClient, greetingTask)

		type result struct {
			response entities.FEvalResponse
			err      error
		}
		resultChan := make(chan result, 1)
		go func() {
			response, err := client.FEval(ctx, mockLogger, fevalRequest)
			resultChan <- result{response, err}
		}()

		synctest.Wait()
		select {
		case <-resultChan:
			t.Fatal("FEval returned before the greeting completed")
		default:
		}

		close(greetingReleased)
		got := <-resultChan
		require.NoError(t, got.err)
		assert.Equal(t, fevalResponse, got.response)
	})
}

func TestGreetingSessionClient_EvalWithCapture_BlocksUntilGreetingDone(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
		defer mockSessionClient.AssertExpectations(t)

		ctx := t.Context()
		evalRequest := entities.EvalRequest{Code: "disp('hi')"}
		evalResponse := entities.EvalResponse{ConsoleOutput: "hi"}

		greetingReleased := make(chan struct{})
		greetingTask := blockedGreetingTask(mockLogger, greetingReleased)

		mockSessionClient.EXPECT().
			EvalWithCapture(ctx, mockLogger.AsMockArg(), evalRequest).
			Return(evalResponse, nil).
			Once()

		client := matlabmanager.NewGreetingSessionClient(mockSessionClient, greetingTask)

		type result struct {
			response entities.EvalResponse
			err      error
		}
		resultChan := make(chan result, 1)
		go func() {
			response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)
			resultChan <- result{response, err}
		}()

		synctest.Wait()
		select {
		case <-resultChan:
			t.Fatal("EvalWithCapture returned before the greeting completed")
		default:
		}

		close(greetingReleased)
		got := <-resultChan
		require.NoError(t, got.err)
		assert.Equal(t, evalResponse, got.response)
	})
}

func TestGreetingSessionClient_SecondCall_DoesNotWaitOnceGreetingDone(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
		defer mockSessionClient.AssertExpectations(t)

		ctx := t.Context()
		evalRequest := entities.EvalRequest{Code: "1 + 1"}
		evalResponse := entities.EvalResponse{ConsoleOutput: "ans = 2"}

		greetingReleased := make(chan struct{})
		greetingTask := blockedGreetingTask(mockLogger, greetingReleased)

		// Both calls delegate to the underlying client once the gate is open.
		mockSessionClient.EXPECT().
			Eval(ctx, mockLogger.AsMockArg(), evalRequest).
			Return(evalResponse, nil).
			Twice()

		client := matlabmanager.NewGreetingSessionClient(mockSessionClient, greetingTask)

		// Act: the first call parks on the gate while the greeting is still in flight.
		type result struct {
			response entities.EvalResponse
			err      error
		}
		resultChan := make(chan result, 1)
		go func() {
			response, err := client.Eval(ctx, mockLogger, evalRequest)
			resultChan <- result{response, err}
		}()

		synctest.Wait()
		select {
		case <-resultChan:
			t.Fatal("first Eval returned before the greeting completed")
		default:
		}

		// Release the greeting; the first call unblocks and the gate is now closed for good.
		close(greetingReleased)
		got := <-resultChan
		require.NoError(t, got.err)
		assert.Equal(t, evalResponse, got.response)

		// Act + Assert: a second call must sail straight through the already-closed gate.
		// Calling it directly means a regression that re-waits on the gate leaves every
		// goroutine durably blocked with nothing left to release it, which synctest fails
		// as a deadlock rather than letting the call return.
		response, err := client.Eval(ctx, mockLogger, evalRequest)
		require.NoError(t, err)
		assert.Equal(t, evalResponse, response)
	})
}

func TestGreetingSessionClient_Eval_ContextCancelled_DelegatesToUnderlying(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
		defer mockSessionClient.AssertExpectations(t)

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		evalRequest := entities.EvalRequest{Code: "1 + 1"}

		greetingReleased := make(chan struct{})
		defer close(greetingReleased)
		greetingTask := blockedGreetingTask(mockLogger, greetingReleased)

		// The gate falls through on cancellation and still delegates to the underlying
		// client with the cancelled context; the underlying returns the ctx error.
		mockSessionClient.EXPECT().
			Eval(ctx, mockLogger.AsMockArg(), evalRequest).
			Return(entities.EvalResponse{}, context.Canceled).
			Once()

		client := matlabmanager.NewGreetingSessionClient(mockSessionClient, greetingTask)

		type result struct {
			response entities.EvalResponse
			err      error
		}
		resultChan := make(chan result, 1)
		go func() {
			response, err := client.Eval(ctx, mockLogger, evalRequest)
			resultChan <- result{response, err}
		}()

		synctest.Wait()
		select {
		case <-resultChan:
			t.Fatal("Eval returned before the context was cancelled")
		default:
		}

		// Cancelling the context releases the gate; Eval then delegates to the underlying.
		cancel()
		got := <-resultChan

		// Assert
		require.ErrorIs(t, got.err, context.Canceled)
	})
}

func TestGreetingSessionClient_Ping_DoesNotGate(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
		defer mockSessionClient.AssertExpectations(t)

		ctx := t.Context()

		greetingReleased := make(chan struct{})
		defer close(greetingReleased)
		greetingTask := blockedGreetingTask(mockLogger, greetingReleased)

		mockSessionClient.EXPECT().
			Ping(ctx, mockLogger.AsMockArg()).
			Return(entities.PingResponse{IsAlive: true}).
			Once()

		client := matlabmanager.NewGreetingSessionClient(mockSessionClient, greetingTask)

		// Act: Ping with the greeting still held open.
		response := client.Ping(ctx, mockLogger)

		// Assert: it returned without waiting on the greeting.
		assert.True(t, response.IsAlive)
	})
}

func TestGreetingSessionClient_Eval_DelegatesWhenGreetingErrored(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
		defer mockSessionClient.AssertExpectations(t)

		ctx := t.Context()
		evalRequest := entities.EvalRequest{Code: "1 + 1"}
		evalResponse := entities.EvalResponse{ConsoleOutput: "ans = 2"}

		// A greeting task whose op fails still completes, so the gate must open.
		greetingTask := asyncrunner.Run(t.Context(), mockLogger, time.Hour, "greeting", func(context.Context) error {
			return assert.AnError
		})

		mockSessionClient.EXPECT().
			Eval(ctx, mockLogger.AsMockArg(), evalRequest).
			Return(evalResponse, nil).
			Once()

		client := matlabmanager.NewGreetingSessionClient(mockSessionClient, greetingTask)

		// Act: the greeting failed, but the gate still releases and Eval delegates.
		response, err := client.Eval(ctx, mockLogger, evalRequest)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, evalResponse, response)
	})
}
