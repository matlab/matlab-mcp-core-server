// Copyright 2026 The MathWorks, Inc.

package embeddedconnector_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_MessageServiceEndpoint_NormalisesBasePath(t *testing.T) {
	const host = "localhost"
	const port = "31515"

	testCases := []struct {
		name     string
		basePath string
		expected string
	}{
		{
			name:     "Empty",
			basePath: "",
			expected: "https://localhost:31515/messageservice/json/secure",
		},
		{
			name:     "BareSlash",
			basePath: "/",
			expected: "https://localhost:31515/messageservice/json/secure",
		},
		{
			name:     "LeadingAndTrailingSlash",
			basePath: "/matlab/",
			expected: "https://localhost:31515/matlab/messageservice/json/secure",
		},
		{
			name:     "NoLeadingSlash",
			basePath: "matlab/",
			expected: "https://localhost:31515/matlab/messageservice/json/secure",
		},
		{
			name:     "NoTrailingSlash",
			basePath: "/user/foo",
			expected: "https://localhost:31515/user/foo/messageservice/json/secure",
		},
		{
			name:     "NestedNoSlashes",
			basePath: "user/foo/session-id",
			expected: "https://localhost:31515/user/foo/session-id/messageservice/json/secure",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &embeddedconnector.Client{}
			client.SetHost(host)
			client.SetPort(port)
			client.SetBasePath(tc.basePath)

			endpoint, err := client.MessageServiceEndpoint("secure")

			require.NoError(t, err)
			assert.Equal(t, tc.expected, endpoint)
		})
	}
}

func TestClient_MessageServiceEndpoint_Channel(t *testing.T) {
	client := &embeddedconnector.Client{}
	client.SetHost("localhost")
	client.SetPort("31515")
	client.SetBasePath("/matlab/")

	endpoint, err := client.MessageServiceEndpoint("state")

	require.NoError(t, err)
	assert.Equal(t, "https://localhost:31515/matlab/messageservice/json/state", endpoint)
}
