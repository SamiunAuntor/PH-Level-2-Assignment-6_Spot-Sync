package middleware

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	apperrors "spotsync/apperror"
	"spotsync/response"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		_ = c.JSON(appErr.StatusCode, response.Error(appErr.Message, appErr.Details))
		return
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		_ = c.JSON(http.StatusBadRequest, response.Error("Validation failed", apperrors.ValidationDetails(validationErrors)))
		return
	}

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		message, details := apperrors.FromHTTPError(httpErr)
		_ = c.JSON(httpErr.Code, response.Error(message, details))
		return
	}

	_ = c.JSON(http.StatusInternalServerError, response.Error("Internal server error", nil))
}
