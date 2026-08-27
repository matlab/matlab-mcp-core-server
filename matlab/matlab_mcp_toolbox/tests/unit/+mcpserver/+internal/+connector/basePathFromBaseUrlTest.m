classdef basePathFromBaseUrlTest < matlab.unittest.TestCase

    % Copyright 2026 The MathWorks, Inc.

    methods (Test)
        function testBasePathFromBaseUrl_WithPathPrefix(testCase)
            % Arrange
            baseUrl = "https://127.0.0.1:31415/user/testuser/session-123/";

            % Act
            basePath = mcpserver.internal.connector.basePathFromBaseUrl(baseUrl);

            % Assert
            testCase.verifyEqual(basePath, "/user/testuser/session-123/", ...
                "Base path should be the URL path component");
        end

        function testBasePathFromBaseUrl_NoPathPrefix(testCase)
            % Arrange
            baseUrl = "https://127.0.0.1:31415/";

            % Act
            basePath = mcpserver.internal.connector.basePathFromBaseUrl(baseUrl);

            % Assert
            testCase.verifyEqual(basePath, "/", ...
                "Base path should be the root path when there is no prefix");
        end

        function testBasePathFromBaseUrl_ReturnsString(testCase)
            % Arrange
            baseUrl = "https://127.0.0.1:31415/matlab/";

            % Act
            basePath = mcpserver.internal.connector.basePathFromBaseUrl(baseUrl);

            % Assert
            testCase.verifyClass(basePath, "string", ...
                "Base path should be returned as a string");
        end
    end

end
