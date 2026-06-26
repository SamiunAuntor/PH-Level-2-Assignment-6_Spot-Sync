package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	apperror "spotsync/apperror"
	"spotsync/dto"
	appmiddleware "spotsync/middleware"
	"spotsync/response"
	"spotsync/service"
)

type ReservationHandler struct {
	reservationService service.ReservationService
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

func (h *ReservationHandler) Create(c echo.Context) error {
	userID, err := appmiddleware.GetUserID(c)
	if err != nil {
		return err
	}

	var request dto.CreateReservationRequest
	if err := c.Bind(&request); err != nil {
		return apperror.BadRequest("Invalid request body", nil, err)
	}

	if err := c.Validate(&request); err != nil {
		return err
	}

	reservation, err := h.reservationService.Create(c.Request().Context(), userID, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response.Success("Reservation confirmed successfully", reservation))
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, err := appmiddleware.GetUserID(c)
	if err != nil {
		return err
	}

	reservations, err := h.reservationService.GetMyReservations(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Success("My reservations retrieved successfully", reservations))
}

func (h *ReservationHandler) GetAll(c echo.Context) error {
	reservations, err := h.reservationService.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Success("All reservations retrieved successfully", reservations))
}

func (h *ReservationHandler) Cancel(c echo.Context) error {
	userID, err := appmiddleware.GetUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperror.BadRequest("Invalid reservation ID", nil, err)
	}

	if err := h.reservationService.Cancel(c.Request().Context(), id, userID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Success("Reservation cancelled successfully", nil))
}
