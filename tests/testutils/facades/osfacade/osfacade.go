// Copyright 2025 The MathWorks, Inc.

package osfacade

import "os"

// RealEnvironment implements Environment using the os package
type RealEnvironment struct{}

func (RealEnvironment) Getenv(key string) string {
	return os.Getenv(key)
}
