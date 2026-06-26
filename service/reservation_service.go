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
	GetMyReservations(ctx context.Context, userID int) ([]dto.MyReservationResponse, error)
}

type reservationService struct {
	reservationRepository repository.ReservationRepository
}

func NewReservationService(
	reservationRepository repository.ReservationRepository,
) ReservationService {
	return &reservationService{
		reservationRepository: reservationRepository,
	}
}

func (s *reservationService) Create(ctx context.Context, userID int, request dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
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

func (s *reservationService) GetMyReservations(ctx context.Context, userID int) ([]dto.MyReservationResponse, error) {
	reservations, err := s.reservationRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MyReservationResponse, 0, len(reservations))
	for _, reservation := range reservations {
		responses = append(responses, dto.MyReservationResponse{
			ID:           reservation.ID,
			LicensePlate: reservation.LicensePlate,
			Status:       reservation.Status,
			Zone: dto.ReservationZoneSummary{
				ID:   reservation.Zone.ID,
				Name: reservation.Zone.Name,
				Type: reservation.Zone.Type,
			},
			CreatedAt: reservation.CreatedAt,
		})
	}

	return responses, nil
}
