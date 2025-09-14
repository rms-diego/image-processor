package exception

import (
	"fmt"
)

type AppError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Trace   *error `json:"trace"`
}

type Error interface {
	Error() string
}

func (e *AppError) Error() string {
	if e.Trace != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Trace)
	}
	return e.Message
}

func New(message string, code int, err *error) Error {
	return &AppError{Message: message, Code: code, Trace: err}
}
