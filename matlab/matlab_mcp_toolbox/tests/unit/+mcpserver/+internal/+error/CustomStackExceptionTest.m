classdef CustomStackExceptionTest < matlab.unittest.TestCase
    %CustomStackExceptionTest Tests for mcpserver.internal.error.CustomStackException

    % Copyright 2026 The MathWorks, Inc.

    methods (Test)
        function testCustomStackException_PreservesIdentifier(testCase)
            % Arrange
            me = MException('my:custom:id', 'boom');
            stack = makeStack({'a.m', 'a', 1});

            % Act
            ex = mcpserver.internal.error.CustomStackException.fromStack(me, stack);

            % Assert
            testCase.verifyEqual(ex.identifier, 'my:custom:id', ...
                "Identifier should be copied from the wrapped exception");
        end

        function testCustomStackException_PreservesMessage_PlainText(testCase)
            % Arrange
            me = MException('test:err', 'something went wrong');
            stack = makeStack({'a.m', 'a', 1});

            % Act
            ex = mcpserver.internal.error.CustomStackException.fromStack(me, stack);

            % Assert
            testCase.verifyEqual(ex.message, 'something went wrong');
        end

        function testCustomStackException_PreservesMessage_WithFormatChars(testCase)
            % Arrange
            % Build an MException whose .message contains literal '%' and '\'.
            me = MException('test:fmt', '%s', 'literal %s and \n in message');
            stack = makeStack({'a.m', 'a', 1});

            % Act
            ex = mcpserver.internal.error.CustomStackException.fromStack(me, stack);

            % Assert
            testCase.verifyEqual(ex.message, me.message, ...
                "Format characters in the original message must not be re-interpreted");
        end

        function testCustomStackException_PreservesCauseChain(testCase)
            % Arrange
            cause1 = MException('cause:one', 'first cause');
            cause2 = MException('cause:two', 'second cause');
            me = MException('test:err', 'outer error');
            me = me.addCause(cause1);
            me = me.addCause(cause2);
            stack = makeStack({'a.m', 'a', 1});

            % Act
            ex = mcpserver.internal.error.CustomStackException.fromStack(me, stack);

            % Assert
            testCase.verifyLength(ex.cause, 2);
            testCase.verifyEqual(ex.cause{1}.identifier, 'cause:one');
            testCase.verifyEqual(ex.cause{2}.identifier, 'cause:two');
        end

        function testCustomStackException_AcceptsAnonymousFunctionFrame(testCase)
            % Anonymous function frames have empty file and line=0; the
            % validator must accept them, since MATLAB's own stacks contain
            % these for user code with anonymous functions.

            % Arrange
            me = MException('test:err', 'boom');
            stack = makeStack({'', '@(x)x(0)', 0});

            % Act & Assert
            testCase.verifyWarningFree(@() ...
                mcpserver.internal.error.CustomStackException.fromStack(me, stack), ...
                "Anonymous-function stack frames must be accepted");
        end

        function testCustomStackException_AcceptsEmptyStack(testCase)
            % Arrange
            me = MException('test:err', 'boom');
            stack = struct('file', {}, 'name', {}, 'line', {});

            % Act & Assert
            testCase.verifyWarningFree(@() ...
                mcpserver.internal.error.CustomStackException.fromStack(me, stack), ...
                "Empty stack should be accepted");
        end

        function testCustomStackException_RejectsMissingFields(testCase)
            % Arrange
            me = MException('test:err', 'boom');
            badStack = struct('file', {'a.m'}, 'name', {'a'});  % missing 'line'

            % Act & Assert
            testCase.verifyError(@() ...
                mcpserver.internal.error.CustomStackException.fromStack(me, badStack), ...
                "CustomStackException:CustomStackException", ...
                "Stack missing required fields should be rejected");
        end

        function testCustomStackException_SameScopeThrow_UsesCustomStack(testCase)
            % Arrange
            customStack = makeStack({'inner.m', 'inner', 3}, {'middle.m', 'middle', 2});

            % Act
            try
                throwAsCaller(mcpserver.internal.error.CustomStackException.fromStack( ...
                    MException('test:err', 'boom'), customStack));
            catch caught
                % Assert
                testCase.verifyEqual({caught.stack.name}, {'inner', 'middle'}, ...
                    "Custom stack should be reported when constructed and thrown in the same scope");
            end
        end

        function testCustomStackException_ReportRebuiltFromCustomStack(testCase)
            % Arrange
            customStack = makeStack({'inner.m', 'inner', 3}, {'middle.m', 'middle', 2});

            % Act
            try
                throwAsCaller(mcpserver.internal.error.CustomStackException.fromStack( ...
                    MException('test:err', 'something broke'), customStack));
            catch caught
                report = getReport(caught, 'extended', 'hyperlinks', 'off');

                % Assert
                testCase.verifySubstring(report, 'something broke');
                testCase.verifySubstring(report, 'inner');
                testCase.verifySubstring(report, 'middle');
            end
        end

        function testCustomStackException_NoOptions_UsesOriginalStack(testCase)
            % Arrange
            inner = MException('test:err', 'boom');
            try
                inner.throw();
            catch raised
                me = raised;
            end

            % Act
            ex = mcpserver.internal.error.CustomStackException.fromException(me);

            % Assert 
            testCase.verifyEqual(ex.identifier, 'test:err');
            testCase.verifyEqual(ex.message, 'boom');
        end

        function testCustomStackException_TrimBefore_RemovesNamedFrameAndBelow(testCase)
            % Build a real stack via a wrapper so the names are real.
            try
                trimBefore_helper(@() error('test:err', 'boom'));
            catch raised
                originalNames = {raised.stack.name};
            end

            try
                throwAsCaller(mcpserver.internal.error.CustomStackException.trimmedBefore( ...
                    raised, "trimBefore_helper"));
            catch caught
                trimmedNames = {caught.stack.name};

                % Assert
                testCase.verifyTrue(any(strcmp(originalNames, 'trimBefore_helper')), ...
                    "Sanity: original stack should contain the helper frame");
                testCase.verifyFalse(any(strcmp(trimmedNames, 'trimBefore_helper')), ...
                    "trimBefore_helper frame must be removed by TrimBefore");
            end
        end

        function testCustomStackException_TrimBefore_NoMatch_KeepsFullStack(testCase)
            % Arrange
            try
                error('test:err', 'boom');
            catch raised
                originalLen = numel(raised.stack);
            end

            % Act
            try
                throwAsCaller(mcpserver.internal.error.CustomStackException.trimmedBefore( ...
                    raised, "thisFrameDoesNotExist"));
            catch caught
                % Assert
                testCase.verifyEqual(numel(caught.stack), originalLen, ...
                    "Stack should be unchanged when TrimBefore name is not on the stack");
            end
        end

        function testCustomStackException_EmptyCustomStack_DoesNotEmitWarning(testCase)
            % Returning an empty stack from getStack triggers a MATLAB-emitted
            % warning. The class must avoid that by skipping the override
            % when the custom stack is empty.

            % Arrange (inline errors produce a 0x1 stack)
            try
                error('test:err', 'inline boom');
            catch raised
            end

            % Act
            try
                throwAsCaller(mcpserver.internal.error.CustomStackException.fromException(raised));
            catch caught
                testCase.verifyWarningFree(@() getReport(caught), ...
                    "getReport on a CustomStackException with empty stack must not warn");
            end
        end

    end

end

function stack = makeStack(varargin)
    %makeStack Build a column-vector stack struct from {file, name, line} cells
    stack = struct('file', {}, 'name', {}, 'line', {});
    for k = 1:nargin
        entry = varargin{k};
        stack(k, 1) = struct('file', entry{1}, 'name', entry{2}, 'line', entry{3});
    end
end

function trimBefore_helper(thunk)
    %trimBefore_helper Calls thunk so its frame appears on the error stack
    thunk();
end
