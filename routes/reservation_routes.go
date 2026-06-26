package routes

import (
	"github.com/labstack/echo/v4"

	"spotsync/config"
	"spotsync/handler"
	appmiddleware "spotsync/middleware"
)

func RegisterReservationRoutes(e *echo.Echo, cfg config.Config, reservationHandler *handler.ReservationHandler) {
	reservationsGroup := e.Group("/api/v1/reservations")
	reservationsGroup.Use(appmiddleware.JWTAuthMiddleware(cfg.JWTSecret))
	reservationsGroup.POST("", reservationHandler.Create)
	reservationsGroup.GET("/my-reservations", reservationHandler.GetMyReservations)
	reservationsGroup.DELETE("/:id", reservationHandler.Cancel)

	adminReservationsGroup := e.Group("/api/v1/reservations")
	adminReservationsGroup.Use(appmiddleware.JWTAuthMiddleware(cfg.JWTSecret))
	adminReservationsGroup.Use(appmiddleware.AdminOnlyMiddleware())
	adminReservationsGroup.GET("", reservationHandler.GetAll)
}
