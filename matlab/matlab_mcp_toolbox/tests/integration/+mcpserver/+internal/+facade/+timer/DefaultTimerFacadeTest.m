classdef DefaultTimerFacadeTest < matlab.unittest.TestCase
    %DefaultTimerFacadeTest Functional tests for the real-timer facade

    % Copyright 2026 The MathWorks, Inc.

    methods (Access = private)
        function facade = createFacade(testCase)
            facade = mcpserver.internal.facade.timer.DefaultTimerFacade( ...
                0.25, 'fixedRate', 'drop', @(~, ~) []);
            testCase.addTeardown(@() delete(facade));
        end
    end

    methods (Test)
        function testDefaultTimerFacade_StartStop_HappyPath(testCase)
            facade = testCase.createFacade();

            testCase.verifyWarningFree(@() facade.Start());
            testCase.verifyWarningFree(@() facade.Stop());
        end

        function testDefaultTimerFacade_Stop_RawTimerDeleted_DoesNotError(testCase)
            % Reproduces the destructor warning: the underlying MATLAB timer
            % is torn down (as happens on global cleanup) while the facade
            % handle survives. Stop() must no-op rather than error.
            facade = testCase.createFacade();
            delete(timerfindall);

            testCase.verifyWarningFree(@() facade.Stop());
        end

        function testDefaultTimerFacade_Start_RawTimerDeleted_DoesNotError(testCase)
            facade = testCase.createFacade();
            delete(timerfindall);

            testCase.verifyWarningFree(@() facade.Start());
        end

        function testDefaultTimerFacade_Delete_RawTimerDeleted_DoesNotError(testCase)
            facade = mcpserver.internal.facade.timer.DefaultTimerFacade( ...
                0.25, 'fixedRate', 'drop', @(~, ~) []);
            delete(timerfindall);

            testCase.verifyWarningFree(@() delete(facade));
        end
    end

end
