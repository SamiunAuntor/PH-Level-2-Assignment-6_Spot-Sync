package middleware

import (
	"github.com/labstack/echo/v4"

	apperror "spotsync/apperror"
	"spotsync/models"
)

func AdminOnlyMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, err := GetUserRole(c)
			if err != nil {
				return err
			}

			if role != models.RoleAdmin {
				return apperror.Forbidden("Forbidden", nil, nil)
			}

			return next(c)
		}
	}
}
