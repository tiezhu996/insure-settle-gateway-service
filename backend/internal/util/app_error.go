package util

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel 错误。
var (
	ErrNotFound        = errors.New("not found")
	ErrConflict        = errors.New("conflict")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrValidation      = errors.New("validation error")
	ErrRateLimited     = errors.New("rate limited")
	ErrApiKeyInvalid   = errors.New("api key invalid")
	ErrClientDisabled  = errors.New("client disabled")
)

// AppError 业务错误。
type AppError struct {
	Code       int
	HTTPStatus int
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func NewAppError(code, status int, message string, err error) *AppError {
	return &AppError{Code: code, HTTPStatus: status, Message: message, Err: err}
}

func BadRequest(message string, err error) *AppError {
	return NewAppError(1000, http.StatusBadRequest, message, err)
}
func UnauthorizedError(message string, err error) *AppError {
	return NewAppError(1001, http.StatusUnauthorized, message, err)
}
func ForbiddenError(message string, err error) *AppError {
	return NewAppError(1002, http.StatusForbidden, message, err)
}
func NotFoundError(message string, err error) *AppError {
	return NewAppError(1003, http.StatusNotFound, message, err)
}
func ConflictError(message string, err error) *AppError {
	return NewAppError(1004, http.StatusConflict, message, err)
}
func InternalError(message string, err error) *AppError {
	return NewAppError(1007, http.StatusInternalServerError, message, err)
}

// MapNotFound 尝试把 not-found 语义的错误映射为 404；只识别 *AppError，
// 对错误链中被 %v 断链包装的 sentinel 无法识别，会返回 nil 导致上层按 500 处理。
func MapNotFound(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}
