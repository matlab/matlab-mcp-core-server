classdef CustomStackException < MException
    %CustomStackException Wraps an MException to report a customized stack.
    %   Build via fromException / fromStack / trimmedBefore. The caller must
    %   throwAsCaller in the same scope as construction for the custom stack
    %   to take effect; afterwards it behaves as a normal MException.

    % Copyright 2026 The MathWorks, Inc.

    properties (Access = private)
        ConstructionStack;
        CustomStack;
        OriginalException;
    end

    methods (Static)
        function obj = fromException(me)
            arguments
                me(1, 1) MException
            end
            obj = mcpserver.internal.error.CustomStackException(me, me.stack);
        end

        function obj = fromStack(me, stack)
            arguments
                me(1, 1) MException
                stack(:, 1) struct {mustBeStackStruct}
            end
            obj = mcpserver.internal.error.CustomStackException(me, stack);
        end

        function obj = trimmedBefore(me, funcName)
            arguments
                me(1, 1) MException
                funcName(1, 1) string {mustBeNonzeroLengthTextScalar}
            end
            obj = mcpserver.internal.error.CustomStackException( ...
                me, trimStackBefore(me.stack, funcName));
        end
    end

    methods (Access = private)
        function obj = CustomStackException(me, customStack)
            % Private: use the static factories.

            % MException treats its message argument as an sprintf format
            % string, so escape %% and \\ to reproduce me.message verbatim.
            formattedMessage = replace(me.message, ["%", "\"], ["%%", "\\"]);
            obj@MException(me.identifier, formattedMessage, me.arguments{:});
            obj.type = me.type;
            for idx = 1:numel(me.cause)
                obj = obj.addCause(me.cause{idx});
            end
            if ~isempty(me.Correction)
                obj = obj.addCorrection(me.Correction);
            end

            % Depth 3 skips this constructor and the calling factory,
            % capturing the client scope the same-scope throwAsCaller sees.
            obj.ConstructionStack = dbstack(3, "-completenames");
            obj.CustomStack = customStack;
            obj.OriginalException = me;
        end
    end

    methods (Access = protected)
        function stack = getStack(obj)
            stack = getStack@MException(obj);

            % Only override for the same-scope throwAsCaller. Skip when
            % CustomStack is empty: an empty stack from getStack warns.
            if isequal(stack, obj.ConstructionStack) && ~isempty(obj.CustomStack)
                stack = obj.CustomStack;
            end
        end
    end

end

function stack = trimStackBefore(stack, callerName)
    if isempty(stack)
        return;
    end
    idx = find({stack.name} == callerName, 1);
    if isempty(idx)
        return;
    end
    stack(idx:end) = [];
end

function mustBeStackStruct(stack)

    fields = ["file", "name", "line"];
    for field = fields
        if ~isfield(stack, field)
            error("CustomStackException:CustomStackException", ...
                  "Specified stack must contain field: " + field);
        end
    end

    for idx = 1:numel(stack)
        mustBeStackFileValue(stack(idx).file);
        mustBeStackNameValue(stack(idx).name);
        mustBeStackLineValue(stack(idx).line);
    end
end

function mustBeStackFileValue(value)
    % File may be empty for anonymous functions and inline errors.
    mustBeA(value, "char");
end

function mustBeStackNameValue(value)
    mustBeA(value, "char");
    mustBeTextScalar(value);
    mustBeNonzeroLengthText(value);
end

function mustBeStackLineValue(value)
    mustBeA(value, "double");
    mustBeScalarOrEmpty(value);
    mustBeNonempty(value);
    mustBeNonnegative(value);
    mustBeInteger(value);
end

function mustBeNonzeroLengthTextScalar(value)
    mustBeTextScalar(value);
    mustBeNonzeroLengthText(value);
end
