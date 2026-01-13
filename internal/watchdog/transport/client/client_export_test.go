// Copyright 2025-2026 The MathWorks, Inc.

package client

import (
	"time"
)

func NewClient(
	osLayer OSLayer,
	httpClientFactory HTTPClientFactory,
	loggerFactory LoggerFactory,
) *Client {
	return newClient(
		osLayer,
		httpClientFactory,
		loggerFactory,
	)
}

func (c *Client) SetSocketWaitTimeout(socketWaitTimeout time.Duration) {
	c.socketWaitTimeout = socketWaitTimeout
}

func (c *Client) SetSocketRetryInterval(socketRetryInterval time.Duration) {
	c.socketRetryInterval = socketRetryInterval
}
