package routes

import (
	"github.com/labstack/echo/v4"

	"spotsync/config"
	"spotsync/handler"
	appmiddleware "spotsync/middleware"
)

func RegisterParkingZoneRoutes(e *echo.Echo, cfg config.Config, parkingZoneHandler *handler.ParkingZoneHandler) {
	zonesGroup := e.Group("/api/v1/zones")
	zonesGroup.GET("", parkingZoneHandler.GetAll)
	zonesGroup.GET("/:id", parkingZoneHandler.GetByID)

	adminGroup := e.Group("/api/v1/zones")
	adminGroup.Use(appmiddleware.JWTAuthMiddleware(cfg.JWTSecret))
	adminGroup.Use(appmiddleware.AdminOnlyMiddleware())
	adminGroup.POST("", parkingZoneHandler.Create)
	adminGroup.PATCH("/:id", parkingZoneHandler.Update)
	adminGroup.DELETE("/:id", parkingZoneHandler.Delete)
}
