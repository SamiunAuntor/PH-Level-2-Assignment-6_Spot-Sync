package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"spotsync/config"
	"spotsync/database"
	"spotsync/models"
	apperrors "spotsync/apperror"
	appmiddleware "spotsync/middleware"
	"spotsync/validator"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.AutoMigrate(db, &models.User{}, &models.ParkingZone{}, &models.Reservation{}); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Validator = validator.New()
	e.HTTPErrorHandler = appmiddleware.HTTPErrorHandler

	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: config.ParseAllowedOrigins(cfg.CORSAllowedOrigins),
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions},
	}))

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- e.Start(":" + cfg.Port)
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	case <-shutdownCtx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
}
