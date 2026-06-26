package service

import (
	"context"
	"strings"

	apperror "spotsync/apperror"
	"spotsync/dto"
	"spotsync/models"
	"spotsync/repository"
)

type ReservationService interface {
	Create(ctx context.Context, userID int, request dto.CreateReservationRequest) (*dto.ReservationResponse, error)
}

type reservationService struct {
	reservationRepository repository.ReservationRepository
	parkingZoneRepository repository.ParkingZoneRepository
}

func NewReservationService(
	reservationRepository repository.ReservationRepository,
	parkingZoneRepository repository.ParkingZoneRepository,
) ReservationService {
	return &reservationService{
		reservationRepository: reservationRepository,
		parkingZoneRepository: parkingZoneRepository,
	}
}

func (s *reservationService) Create(ctx context.Context, userID int, request dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	if _, err := s.parkingZoneRepository.FindByID(ctx, request.ZoneID); err != nil {
		return nil, err
	}

	licensePlate := strings.ToUpper(strings.TrimSpace(request.LicensePlate))
	if licensePlate == "" {
		return nil, apperror.BadRequest("Validation failed", map[string]string{
			"license_plate": "License plate is required",
		}, nil)
	}

	reservation := &models.Reservation{
		UserID:       userID,
		ZoneID:       request.ZoneID,
		LicensePlate: licensePlate,
		Status:       models.ReservationStatusActive,
	}

	if err := s.reservationRepository.Create(ctx, reservation); err != nil {
		return nil, err
	}

	response := dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}

	return &response, nil
}
