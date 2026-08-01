// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector

import "fmt"

type matlabError struct {
	message string
}

func newMATLABError(message string) matlabError {
	return matlabError{
		message: message,
	}
}

func (e matlabError) Error() string {
	return fmt.Sprintf("matlab error: %v", e.message)
}
