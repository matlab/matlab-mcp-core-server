// Copyright 2026 The MathWorks, Inc.

package entities

type MCPClientInfo struct {
	Name       string
	Title      string
	WebsiteURL string
	Version    string
}

type GreetingInfo struct {
	ClientInfo        MCPClientInfo
	ShowMATLABDesktop bool
}
