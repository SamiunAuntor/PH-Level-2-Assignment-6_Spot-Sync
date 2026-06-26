package routes

import (
	"github.com/labstack/echo/v4"

	"spotsync/handler"
)

func RegisterAuthRoutes(e *echo.Echo, authHandler *handler.AuthHandler) {
	authGroup := e.Group("/api/v1/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
}
