package error

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Type       string
	Message    string
	StatusCode int
}

func Unimplemented(msg string) *APIError {
	return &APIError{
		Type:       "unimplemented",
		Message:    msg,
		StatusCode: http.StatusNotImplemented,
	}
}

func ResourceExhausted(msg string) *APIError {
	return &APIError{
		Type:       "resource_exhausted",
		Message:    msg,
		StatusCode: http.StatusConflict,
	}
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
