package apperror

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AppError struct {
	StatusCode int
	Message    string
	Details    interface{}
	Err        error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func New(statusCode int, message string, details interface{}, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
		Err:        err,
	}
}

func BadRequest(message string, details interface{}, err error) *AppError {
	return New(http.StatusBadRequest, message, details, err)
}

func Unauthorized(message string, details interface{}, err error) *AppError {
	return New(http.StatusUnauthorized, message, details, err)
}

func Forbidden(message string, details interface{}, err error) *AppError {
	return New(http.StatusForbidden, message, details, err)
}

func NotFound(message string, details interface{}, err error) *AppError {
	return New(http.StatusNotFound, message, details, err)
}

func Conflict(message string, details interface{}, err error) *AppError {
	return New(http.StatusConflict, message, details, err)
}

func Internal(message string, err error) *AppError {
	return New(http.StatusInternalServerError, message, nil, err)
}

func FromHTTPError(httpErr *echo.HTTPError) (string, interface{}) {
	if httpErr == nil {
		return "Internal server error", nil
	}

	message := "Request failed"
	switch httpErr.Code {
	case http.StatusBadRequest:
		message = "Bad request"
	case http.StatusUnauthorized:
		message = "Unauthorized"
	case http.StatusForbidden:
		message = "Forbidden"
	case http.StatusNotFound:
		message = "Resource not found"
	case http.StatusMethodNotAllowed:
		message = "Method not allowed"
	default:
		if httpErr.Code >= 500 {
			message = "Internal server error"
		}
	}

	return message, nil
}
