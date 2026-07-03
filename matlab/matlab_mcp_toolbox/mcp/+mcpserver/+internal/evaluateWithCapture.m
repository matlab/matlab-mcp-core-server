function results = evaluateWithCapture(code, options)
    %evaluateWithCapture Evaluate MATLAB code with output capture and figure control

    % Copyright 2026 The MathWorks, Inc.

    arguments
        code(1, 1) string
        options.ShowFigureWindows(1, 1) logical = true
        options.DisplayCapturedEvents(1, 1) logical = false
        options.DisplayDuringExecution(1, 1) logical = true
        options.FigureVisibilityController(1, 1) mcpserver.internal.display.FigureVisibilityController = mcpserver.internal.display.DefaultFigureVisibilityController()
        options.OutputCapture(1, 1) mcpserver.internal.capture.OutputCapture = mcpserver.internal.capture.DefaultOutputCapture()
        options.CapturedEventProvider(1, 1) mcpserver.internal.capture.CapturedEventProvider = mcpserver.internal.capture.DefaultCapturedEventProvider()
        options.EventDisplayer(1, 1) mcpserver.internal.display.console.EventDisplayer
        options.EventPoller(1, 1) mcpserver.internal.capture.EventPoller
        options.OutputBuilder(1, 1) mcpserver.internal.output.OutputBuilder = mcpserver.internal.output.OutputBuilder()
    end

    capture = options.OutputCapture;

    % Optional arguments not supplied by the caller are absent from the options struct
    if ~isfield(options, 'EventDisplayer')
        options.EventDisplayer = mcpserver.internal.display.console.DefaultEventDisplayer();
    end

    if options.DisplayDuringExecution && ~isfield(options, 'EventPoller')
        options.EventPoller = mcpserver.internal.capture.DefaultEventPoller( ...
            EventSource=options.CapturedEventProvider, Displayer=options.EventDisplayer);
    end

    if options.ShowFigureWindows
        options.FigureVisibilityController.show();
    else
        options.FigureVisibilityController.hide();
    end
    restoreVisibility = onCleanup(@() options.FigureVisibilityController.restore());

    displayer = options.EventDisplayer;

    % EventPoller extends CapturedEventProvider, so it serves as the canonical event source when active
    hasPoller = isfield(options, 'EventPoller');
    if hasPoller
        poller = options.EventPoller;
        eventSource = poller;
    else
        eventSource = options.CapturedEventProvider;
    end

    capture.enable();
    restoreCapture = onCleanup(@() capture.disable());

    if hasPoller
        poller.start();
        restorePoller = onCleanup(@() stopAndFlush(poller));
    end

    userError = [];
    try
        evalin('base', code);
    catch ME
        userError = ME;
    end

    if hasPoller
        delete(restorePoller);
    end

    drawnow;

    delete(restoreCapture);

    rawEvents = eventSource.getEvents();

    if options.DisplayCapturedEvents && ~hasPoller
        displayer.displayEvents(rawEvents);
    end

    results = options.OutputBuilder.build(rawEvents);

    if ~isempty(userError)
        % Construct + throwAsCaller must be in the same function scope so
        % CustomStackException honors the trimmed stack.
        throwAsCaller(mcpserver.internal.error.CustomStackException.trimmedBefore( ...
            userError, "evaluateWithCapture"));
    end
end

function stopAndFlush(poller)
    poller.stop();
    poller.flush();
end
