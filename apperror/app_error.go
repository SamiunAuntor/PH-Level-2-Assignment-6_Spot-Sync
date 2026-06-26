package apperror

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
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

func ValidationDetails(validationErrors validator.ValidationErrors) map[string]string {
	details := make(map[string]string, len(validationErrors))

	for _, fieldError := range validationErrors {
		field := toJSONFieldName(fieldError.Field())
		if _, exists := details[field]; exists {
			continue
		}

		details[field] = validationMessage(fieldError)
	}

	return details
}

func validationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return titleCase(fieldError.Field()) + " is required"
	case "email":
		return titleCase(fieldError.Field()) + " must be a valid email address"
	case "gt":
		switch fieldError.Field() {
		case "TotalCapacity":
			return "Total capacity must be greater than zero"
		case "PricePerHour":
			return "Price per hour must be greater than zero"
		}
	case "min":
		if fieldError.Field() == "Password" {
			return "Password must be at least 8 characters long"
		}
	case "max":
		if fieldError.Field() == "Password" {
			return "Password must not exceed 72 characters"
		}
	case "oneof":
		if fieldError.Field() == "Role" {
			return "Role must be either driver or admin"
		}
		if fieldError.Field() == "Type" {
			return "Type must be one of general, ev_charging, or covered"
		}
	}

	return titleCase(fieldError.Field()) + " is invalid"
}

func toJSONFieldName(name string) string {
	switch name {
	case "Name":
		return "name"
	case "Email":
		return "email"
	case "Password":
		return "password"
	case "Role":
		return "role"
	case "Type":
		return "type"
	case "TotalCapacity":
		return "total_capacity"
	case "PricePerHour":
		return "price_per_hour"
	default:
		return strings.ToLower(name)
	}
}

func titleCase(value string) string {
	if value == "" {
		return value
	}

	return strings.ToUpper(value[:1]) + strings.ToLower(value[1:])
}
