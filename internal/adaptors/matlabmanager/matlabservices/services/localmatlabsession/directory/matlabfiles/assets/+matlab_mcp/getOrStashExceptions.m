function [returnedException] = getOrStashExceptions(exceptionMessage, resetFlag)
    % getOrStashExceptions will save the given exception in a persistent variable
    % to provide access later on.
    
    % This is largely a re-use of:
    % https://github.com/mathworks/jupyter-matlab-proxy/blob/057564dccb7de37f052e709f5380e3ece0b2c4a1/src/jupyter_matlab_kernel/matlab/%2Bjupyter/getOrStashExceptions.m#L1

    % Passing in an exceptionMessage will always stash it
    % When resetFlag is set, old exceptionMessage if any is returned, and is
    % cleared.

    % Copyright 2025-2026 The MathWorks, Inc.

    persistent stashedException

    % Initialize the persistent variable if it is not already done.
    if isempty(stashedException)
        stashedException = [];
    end

    % When resetFlag is set, return the saved exception as output and reset the
    % saved exception.
    if nargin == 2 && resetFlag == true
        returnedException = stashedException;
        stashedException = [];
        return
    end

    % If there are no input arguments, only return the saved exception
    if nargin == 0
        returnedException = stashedException;
        return
    end

    % Update the saved exception with the given exception.
    stashedException = exceptionMessage;
end