package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	apperror "spotsync/apperror"
	"spotsync/dto"
	"spotsync/response"
	"spotsync/service"
)

type ParkingZoneHandler struct {
	parkingZoneService service.ParkingZoneService
}

func NewParkingZoneHandler(parkingZoneService service.ParkingZoneService) *ParkingZoneHandler {
	return &ParkingZoneHandler{parkingZoneService: parkingZoneService}
}

func (h *ParkingZoneHandler) Create(c echo.Context) error {
	var request dto.CreateParkingZoneRequest
	if err := c.Bind(&request); err != nil {
		return apperror.BadRequest("Invalid request body", nil, err)
	}

	if err := c.Validate(&request); err != nil {
		return err
	}

	zone, err := h.parkingZoneService.Create(c.Request().Context(), request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response.Success("Parking zone created successfully", zone))
}

func (h *ParkingZoneHandler) GetAll(c echo.Context) error {
	zones, err := h.parkingZoneService.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Success("Parking zones retrieved successfully", zones))
}

func (h *ParkingZoneHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperror.BadRequest("Invalid parking zone ID", nil, err)
	}

	zone, err := h.parkingZoneService.GetByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Success("Parking zone retrieved successfully", zone))
}

func (h *ParkingZoneHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperror.BadRequest("Invalid parking zone ID", nil, err)
	}

	var request dto.UpdateParkingZoneRequest
	if err := c.Bind(&request); err != nil {
		return apperror.BadRequest("Invalid request body", nil, err)
	}

	if err := c.Validate(&request); err != nil {
		return err
	}

	zone, err := h.parkingZoneService.Update(c.Request().Context(), id, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Success("Parking zone updated successfully", zone))
}

func (h *ParkingZoneHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperror.BadRequest("Invalid parking zone ID", nil, err)
	}

	if err := h.parkingZoneService.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Success("Parking zone deleted successfully", nil))
}
