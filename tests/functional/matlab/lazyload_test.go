// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/tests/testutils/mockmatlab"
	"github.com/matlab/matlab-mcp-server/tests/testutils/mockmatlab/mockruntime"
	"github.com/stretchr/testify/suite"
)

const expectedConnectTitleEvals = 3

type LazyLoadTestSuite struct {
	MockMATLABTestSuite
}

func TestLazyLoadTestSuite(t *testing.T) {
	suite.Run(t, new(LazyLoadTestSuite))
}

func (s *LazyLoadTestSuite) TestLazyLoad_MATLABStartsOnFirstToolCall() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil)
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	result, err := session.ListTools(ctx, nil)
	s.Require().NoError(err, "should list tools")
	s.NotEmpty(result.Tools, "should have tools registered")

	instanceEvents, err := session.ReadInstanceEvents()
	s.Require().NoError(err)
	s.Empty(instanceEvents, "MATLAB should not have started before a tool call")

	output, err := session.EvaluateCode(ctx, "disp('hello')", s.T().TempDir())
	s.Require().NoError(err, "first tool call should trigger MATLAB start")
	s.Contains(output, "disp('hello')")

	instanceEvents = s.WaitForInstanceEvents(session, 30*time.Second,
		func(events []mockruntime.InstanceEvents) bool {
			return len(events) == 1 && countTitleEvals(events[0]) == expectedConnectTitleEvals
		},
		"desktop connect should check the release to pick the title API, then read the title and write it back once each")

	s.Require().Len(instanceEvents, 1, "MATLAB should have started after first tool call")
	s.Equal("happy", instanceEvents[0].StartedMode())
	s.Equal(1, instanceEvents[0].CountEvent(mockruntime.EventStarted), "should have exactly one started event")
	s.True(instanceEvents[0].HasEvalMatching(isCdEval), "server should have sent a cd() eval to set the working directory")
	s.True(instanceEvents[0].HasEvalMatching(hasUserCode("disp('hello')")), "should have recorded the user eval")
	s.Equal(
		[]string{mockruntime.EventStarted, mockruntime.EventFeval, mockruntime.EventEval, mockruntime.EventEval},
		eventTypesExcludingTitle(instanceEvents[0]),
		"excluding the desktop connection-title evals, should have exactly: started, greeting feval, cd() eval, user eval",
	)
	s.Equal(expectedConnectTitleEvals, countTitleEvals(instanceEvents[0]),
		"desktop connect should check the release to pick the title API, then read the title and write it back once each")
}

func (s *LazyLoadTestSuite) TestEagerLoad_MATLABStartsOnSessionCreation() {
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil, "--initialize-matlab-on-startup")
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	instanceEvents := s.WaitForInstanceEvents(session, 30*time.Second,
		func(events []mockruntime.InstanceEvents) bool {
			return len(events) == 1 && events[0].HasEvent(mockruntime.EventStarted)
		},
		"MATLAB should have started before any tool call")

	s.Equal("happy", instanceEvents[0].StartedMode())
	s.False(instanceEvents[0].HasEvalMatching(isCdEval), "no cd() eval should have happened without a tool call")
	s.Zero(instanceEvents[0].CountEvent(mockruntime.EventEval),
		"no evals should happen before a tool call; the greeting and desktop connection title ops are fevals")
}

func (s *LazyLoadTestSuite) TestLazyLoad_SecondToolCallReusesSession() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil)
	s.Require().NoError(err)
	defer s.CleanupSession(session, true)

	output, err := session.EvaluateCode(ctx, "disp('first')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('first')")

	output, err = session.EvaluateCode(ctx, "disp('second')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('second')")

	instanceEvents, err := session.ReadInstanceEvents()
	s.Require().NoError(err)
	s.Require().Len(instanceEvents, 1, "only one mock MATLAB instance should have been created")
	s.True(instanceEvents[0].HasEvalMatching(hasUserCode("disp('first')")), "should have recorded first user eval")
	s.True(instanceEvents[0].HasEvalMatching(hasUserCode("disp('second')")), "should have recorded second user eval")
	s.Equal(2, instanceEvents[0].CountEvent(mockruntime.EventEval)-countCdEvals(instanceEvents[0]),
		"should have exactly two user evals (excluding cd evals)")
}

func (s *LazyLoadTestSuite) TestReconnection_AfterExit_MATLABRestartsOnNextToolCall() {
	ctx := s.T().Context()
	session, err := s.CreateSession(mockmatlab.HappyConfig(), nil)
	s.Require().NoError(err)
	defer s.CleanupSession(session, false)

	output, err := session.EvaluateCode(ctx, "disp('before exit')", s.T().TempDir())
	s.Require().NoError(err)
	s.Contains(output, "disp('before exit')")

	_, _ = session.EvaluateCode(ctx, "exit()", s.T().TempDir())

	output, err = session.EvaluateCode(ctx, "disp('after reconnect')", s.T().TempDir())
	s.Require().NoError(err, "tool call after exit() should succeed via reconnection")
	s.Contains(output, "disp('after reconnect')")

	instanceEvents, err := session.ReadInstanceEvents()
	s.Require().NoError(err)
	s.Require().Len(instanceEvents, 2, "should have started two different mock MATLAB instances")

	s.Equal(1, instanceEvents[0].ID)
	s.Equal("happy", instanceEvents[0].StartedMode())
	s.True(instanceEvents[0].HasEvalMatching(isCdEval), "first instance should have a cd() eval")
	s.True(instanceEvents[0].HasEvalMatching(hasUserCode("disp('before exit')")), "first instance should have eval before exit")
	s.True(instanceEvents[0].HasEvalMatching(hasUserCode("exit()")), "first instance should have exit eval")
	s.True(instanceEvents[0].HasEvent(mockruntime.EventExitRequested), "first instance should have recorded exit_requested")

	s.Equal(2, instanceEvents[1].ID)
	s.Equal("happy", instanceEvents[1].StartedMode())
	s.True(instanceEvents[1].HasEvalMatching(isCdEval), "second instance should have a cd() eval")
	s.True(instanceEvents[1].HasEvalMatching(hasUserCode("disp('after reconnect')")), "second instance should have recorded eval after reconnect")
	s.False(instanceEvents[1].HasEvent(mockruntime.EventExitRequested), "second instance should not have an exit event")
}

func isCdEval(code string) bool {
	code = strings.TrimPrefix(code, "feature('HotLinks',0);")
	return strings.HasPrefix(code, "cd('") && strings.HasSuffix(code, "')")
}

func hasUserCode(userCode string) func(string) bool {
	return func(code string) bool {
		return strings.HasSuffix(code, userCode)
	}
}

func countCdEvals(ie mockruntime.InstanceEvents) int {
	n := 0
	for _, e := range ie.Events {
		if e.Type == mockruntime.EventEval && isCdEval(e.Code) {
			n++
		}
	}
	return n
}

func isTitleEvent(e mockruntime.Event) bool {
	if e.Type != mockruntime.EventFeval {
		return false
	}
	return e.Function == "eval" || e.Function == "isMATLABReleaseOlderThan"
}

func countTitleEvals(ie mockruntime.InstanceEvents) int {
	n := 0
	for _, e := range ie.Events {
		if isTitleEvent(e) {
			n++
		}
	}
	return n
}

func eventTypesExcludingTitle(ie mockruntime.InstanceEvents) []string {
	var types []string
	for _, e := range ie.Events {
		if isTitleEvent(e) {
			continue
		}
		types = append(types, e.Type)
	}
	return types
}
