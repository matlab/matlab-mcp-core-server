classdef evaluateWithCaptureErrorReportingTest < matlab.unittest.TestCase
    %evaluateWithCaptureErrorReportingTest Integration tests for error reporting from evaluateWithCapture
    %   These tests run user code through evaluateWithCapture end-to-end
    %   and assert the exception surfaced to the caller has the right
    %   identifier, message, stack, and produces no MATLAB-emitted warnings.

    % Copyright 2026 The MathWorks, Inc.

    properties (Access = private)
        ScratchDir(1, 1) string
    end

    methods (TestMethodSetup)
        function createScratch(testCase)
            testCase.ScratchDir = string(fullfile(tempdir, ...
                "evaluateWithCaptureErrorReportingTest_" + matlab.lang.internal.uuid));
            mkdir(testCase.ScratchDir);
            addpath(testCase.ScratchDir);
            testCase.addTeardown(@() rmdir(testCase.ScratchDir, 's'));
            testCase.addTeardown(@() rmpath(testCase.ScratchDir));
        end
    end

    methods (Test)
        function testUserFunctionError_PreservesIdentifierAndCleanMessage(testCase)
            % Arrange
            testCase.writeUserFunc("userFunc", ...
                "function userFunc()", ...
                "    error('user:err', 'something went wrong inside userFunc');", ...
                "end");

            % Act & Assert
            try
                mcpserver.internal.evaluateWithCapture("userFunc()", ...
                    ShowFigureWindows=false, DisplayDuringExecution=false);
                testCase.verifyFail("Should have thrown");
            catch ME
                testCase.verifyEqual(ME.identifier, 'user:err');
                testCase.verifyEqual(ME.message, ...
                    'something went wrong inside userFunc');
            end
        end

        function testUserFunctionError_StackContainsUserFrameNotToolboxFrame(testCase)
            % Arrange
            testCase.writeUserFunc("userFunc", ...
                "function userFunc()", ...
                "    error('user:err', 'boom');", ...
                "end");

            % Act
            try
                mcpserver.internal.evaluateWithCapture("userFunc()", ...
                    ShowFigureWindows=false, DisplayDuringExecution=false);
                testCase.verifyFail("Should have thrown");
            catch ME
                names = string({ME.stack.name});

                % Assert
                testCase.verifyTrue(any(names == "userFunc"), ...
                    "User function frame should appear on the reported stack");
                testCase.verifyFalse(any(names == "evaluateWithCapture"), ...
                    "Toolbox internal frame should not leak to user-visible stack");
            end
        end

        function testUserFunctionError_ReportContainsUserFileAndLine(testCase)
            % Arrange
            testCase.writeUserFunc("userFunc", ...
                "function userFunc()", ...
                "    error('user:err', 'boom');", ...
                "end");

            % Act
            try
                mcpserver.internal.evaluateWithCapture("userFunc()", ...
                    ShowFigureWindows=false, DisplayDuringExecution=false);
                testCase.verifyFail("Should have thrown");
            catch ME
                report = getReport(ME, 'extended', 'hyperlinks', 'off');

                % Assert
                testCase.verifySubstring(report, 'userFunc', ...
                    "Report should mention the failing user function");
                testCase.verifySubstring(report, 'line 2', ...
                    "Report should mention the failing line");
            end
        end

        function testInlineError_DoesNotEmitWarning(testCase)
            % Inline errors produce a 0x1 stack. Returning that empty stack
            % from the CustomStackException override would trigger a
            % MATLAB-emitted "did not return a valid stack" warning. The
            % implementation must avoid that.

            % Act & Assert
            try
                testCase.verifyWarningFree(@() ...
                    mcpserver.internal.evaluateWithCapture( ...
                        "error('test:e', 'inline boom')", ...
                        ShowFigureWindows=false, DisplayDuringExecution=false), ...
                    "Inline error must not trigger a MATLAB-emitted warning");
            catch ME
                testCase.verifyEqual(ME.identifier, 'test:e');
                testCase.verifyEqual(ME.message, 'inline boom');
            end
        end

        function testErrorWithFormatChars_PreservesLiteralMessage(testCase)
            % Act
            try
                mcpserver.internal.evaluateWithCapture( ...
                    "error('test:fmt', '%s', 'literal %s and \n must survive')", ...
                    ShowFigureWindows=false, DisplayDuringExecution=false);
                testCase.verifyFail("Should have thrown");
            catch ME
                % Assert
                testCase.verifyEqual(ME.message, 'literal %s and \n must survive');
            end
        end

        function testAnonymousFunctionError_DoesNotRejectStackFrames(testCase)
            % Anonymous-function user code yields stack frames with empty
            % file and line=0. The implementation must accept these
            % without throwing a validator error.

            % Act & Assert
            try
                mcpserver.internal.evaluateWithCapture( ...
                    "f = @(x) x(0); g = @() f([1 2 3]); g();", ...
                    ShowFigureWindows=false, DisplayDuringExecution=false);
                testCase.verifyFail("Should have thrown");
            catch ME
                % The original error has identifier 'MATLAB:badsubscript'.
                % If our wrapper had thrown, the identifier would be different.
                testCase.verifyEqual(ME.identifier, 'MATLAB:badsubscript', ...
                    "User error should pass through untouched");
            end
        end
    end

    methods (Access = private)
        function writeUserFunc(testCase, name, varargin)
            funcPath = fullfile(testCase.ScratchDir, name + ".m");
            fid = fopen(funcPath, 'w');
            cleanup = onCleanup(@() fclose(fid));
            for k = 1:numel(varargin)
                fprintf(fid, '%s\n', varargin{k});
            end
            rehash;
        end
    end

end
