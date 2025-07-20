package errors

import (
	"fmt"
	"net/http"
)

// AppError is your reusable error type with code & message.
type AppError struct {
	Code    int    // HTTP status code
	Message string // JSON message
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Predefined errors.
var (
	ErrUnauthorized      = &AppError{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrInvalidToken      = &AppError{Code: http.StatusUnauthorized, Message: "invalid or expired token"}
	ErrHeaderRequired    = &AppError{Code: http.StatusBadRequest, Message: "required header is missing"}
	ErrForbidden         = &AppError{Code: http.StatusForbidden, Message: "forbidden"}
	ErrBadRequest        = &AppError{Code: http.StatusBadRequest, Message: "bad request"}
	ErrInternal          = &AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrRateLimitExceeded = &AppError{Code: http.StatusTooManyRequests, Message: "rate limit exceeded"}
	ErrExtractorError    = &AppError{Code: http.StatusForbidden, Message: "error while extracting identifier"}
)
