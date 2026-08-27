function basePath = basePathFromBaseUrl(baseUrl)
%basePathFromBaseUrl Extract the base path from a connector base URL

% Copyright 2026 The MathWorks, Inc.

    parsed = matlab.net.URI(baseUrl);
    basePath = string(parsed.EncodedPath);
end
