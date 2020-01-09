package common

import (
	"fmt"
)

const errInfo = "Error code: %d, ErrorMessage: %s"

type GfsError struct {
	Message string
}

func (e *GfsError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf(e.Message)
}

func NewGfsError(message string) GfsError {
	return GfsError{Message: message}
}
