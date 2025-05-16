package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code, message string, status int, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

var (
	ErrNotFound       = New("NOT_FOUND", "Resource not found", http.StatusNotFound, nil)
	ErrUnauthorized   = New("UNAUTHORIZED", "Unauthorized access", http.StatusUnauthorized, nil)
	ErrForbidden      = New("FORBIDDEN", "Forbidden", http.StatusForbidden, nil)
	ErrBadRequest     = New("BAD_REQUEST", "Bad request", http.StatusBadRequest, nil)
	ErrInternalServer = New("INTERNAL_SERVER_ERROR", "Internal server error", http.StatusInternalServerError, nil)
	ErrConflict       = New("CONFLICT", "Conflict", http.StatusConflict, nil)
)

func Wrap(base *AppError, err error) *AppError {
	return &AppError{
		Code:    base.Code,
		Message: base.Message,
		Status:  base.Status,
		Err:     err,
	}
}

func Is(target error, code string) bool {
	var appErr *AppError
	if errors.As(target, &appErr) {
		return appErr.Code == code
	}
	return false
}

func StatusCode(err error) int {
	switch {
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
