// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector

import (
	httpclient "github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/client"
)

func (c *Client) SetHttpClient(httpClient httpclient.HttpClient) {
	c.httpClient = httpClient
}
