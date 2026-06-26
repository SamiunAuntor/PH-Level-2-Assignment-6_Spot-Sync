package middleware

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	apperror "spotsync/apperror"
	"spotsync/service"
)

type contextKey string

const (
	userIDContextKey contextKey = "auth_user_id"
	roleContextKey   contextKey = "auth_role"
)

func JWTAuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	secret := []byte(jwtSecret)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := strings.TrimSpace(c.Request().Header.Get(echo.HeaderAuthorization))
			if header == "" {
				return apperror.Unauthorized("Missing or invalid token", nil, nil)
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
				return apperror.Unauthorized("Missing or invalid token", nil, nil)
			}

			claims := &service.AuthClaims{}
			token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, apperror.Unauthorized("Missing or invalid token", nil, nil)
				}

				return secret, nil
			})
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					return apperror.Unauthorized("Missing or invalid token", nil, err)
				}

				return apperror.Unauthorized("Missing or invalid token", nil, err)
			}

			if !token.Valid || claims.UserID == 0 || strings.TrimSpace(claims.Role) == "" {
				return apperror.Unauthorized("Missing or invalid token", nil, nil)
			}

			c.Set(string(userIDContextKey), claims.UserID)
			c.Set(string(roleContextKey), claims.Role)

			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (int, error) {
	value := c.Get(string(userIDContextKey))
	userID, ok := value.(int)
	if !ok || userID == 0 {
		return 0, apperror.Unauthorized("Missing or invalid token", nil, nil)
	}

	return userID, nil
}

func GetUserRole(c echo.Context) (string, error) {
	value := c.Get(string(roleContextKey))
	role, ok := value.(string)
	if !ok || strings.TrimSpace(role) == "" {
		return "", apperror.Unauthorized("Missing or invalid token", nil, nil)
	}

	return role, nil
}
