package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apperror "spotsync/apperror"
	"spotsync/dto"
	"spotsync/response"
	"spotsync/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var request dto.RegisterRequest
	if err := c.Bind(&request); err != nil {
		return apperror.BadRequest("Invalid request body", nil, err)
	}

	if err := c.Validate(&request); err != nil {
		return err
	}

	user, err := h.authService.Register(c.Request().Context(), request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response.Success("User registered successfully", user))
}

func (h *AuthHandler) Login(c echo.Context) error {
	var request dto.LoginRequest
	if err := c.Bind(&request); err != nil {
		return apperror.BadRequest("Invalid request body", nil, err)
	}

	if err := c.Validate(&request); err != nil {
		return err
	}

	loginResponse, err := h.authService.Login(c.Request().Context(), request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Success("Login successful", loginResponse))
}
