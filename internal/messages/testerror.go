// Copyright 2026 The MathWorks, Inc.

package messages

// AnError is a test helper analogous to assert.AnError but satisfying messages.Error.
// Use it in tests when verifying error propagation paths that require messages.Error.
var AnError = &anError{} //nolint:gochecknoglobals // AnError is an error

type anError struct{}

func (e *anError) Error() string {
	return "AnError"
}

func (*anError) marker() {}
