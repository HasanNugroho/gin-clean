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
	if message == "" {
		message = defaultMessage(code)
	}
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

var (
	ErrNotFound       = base("NOT_FOUND", http.StatusNotFound)
	ErrUnauthorized   = base("UNAUTHORIZED", http.StatusUnauthorized)
	ErrForbidden      = base("FORBIDDEN", http.StatusForbidden)
	ErrBadRequest     = base("BAD_REQUEST", http.StatusBadRequest)
	ErrInternalServer = base("INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
	ErrConflict       = base("CONFLICT", http.StatusConflict)
)

func base(code string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: defaultMessage(code),
		Status:  status,
	}
}

func Wrap(base *AppError, err error) *AppError {
	return &AppError{
		Code:    base.Code,
		Message: base.Message,
		Status:  base.Status,
		Err:     err,
	}
}

func (e *AppError) WithMessage(msg string) *AppError {
	e.Message = msg
	return e
}

func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

func (e *AppError) WithStatus(status int) *AppError {
	e.Status = status
	return e
}

func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

func Custom(code, message string, status int) *AppError {
	return New(code, message, status, nil)
}

func Is(target error, code string) bool {
	var appErr *AppError
	if errors.As(target, &appErr) {
		return appErr.Code == code
	}
	return false
}

func StatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Status
	}
	return http.StatusInternalServerError
}

func defaultMessage(code string) string {
	switch code {
	case "NOT_FOUND":
		return "Resource not found"
	case "UNAUTHORIZED":
		return "Unauthorized access"
	case "FORBIDDEN":
		return "Forbidden"
	case "BAD_REQUEST":
		return "Bad request"
	case "CONFLICT":
		return "Conflict"
	case "INTERNAL_SERVER_ERROR":
		return "Internal server error"
	default:
		return "Unexpected error"
	}
}
