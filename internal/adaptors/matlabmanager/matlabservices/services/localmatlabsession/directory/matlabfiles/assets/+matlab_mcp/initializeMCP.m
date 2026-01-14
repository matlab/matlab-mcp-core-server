function initializeMCP()
    % initializeMCP intializes the MATLAB sessions for the MATLAB MCP Core Server.

    % Copyright 2025-2026 The MathWorks, Inc.

    connector.ensureServiceOn();

    sessionDir = getenv("MW_MCP_SESSION_DIR");
    securePortFile = fullfile(sessionDir, "connector.securePort");

    securePort = connector.securePort();

    % Record the port that connector is listening on so the MCP server can send messages to MATLAB
    securePortFileID = fopen(securePortFile, "w");
    closeSecurePortFile = onCleanup(@() fclose(securePortFileID));
    fprintf(securePortFileID, "%d", securePort);
end
