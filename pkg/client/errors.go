package client

import (
	"fmt"
)

// ResponseError is the error returned from the client when Key Manager responds with an error or
// non-success HTTP status code. ResponseError gives
// access to the underlying errors and status code.
type ResponseError struct {
	// StatusCode is the HTTP status code.
	StatusCode int `json:"statusCode,omitempty" example:"404"`

	// ErrorCode is the key manager error code.
	ErrorCode string `json:"code,omitempty" example:"IR001"`

	// Errors are the underlying errors returned by Vault.
	Message string `json:"message" example:"error message"`
}

// ErrorResponse is the raw error type returned from the key manager
type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
	Code    string `json:"code,omitempty" example:"IR001"`
}

// Error returns a human-readable error string for the response error.
func (r *ResponseError) Error() string {
	return fmt.Sprintf("Error making API request.\nCode: %s. %s:\nStatus: %d.", r.ErrorCode, r.Message, r.StatusCode)
}
