package service

import (
	"context"
	"strings"

	apperror "spotsync/apperror"
	"spotsync/dto"
	"spotsync/models"
	"spotsync/repository"
)

type ParkingZoneService interface {
	Create(ctx context.Context, request dto.CreateParkingZoneRequest) (*dto.ParkingZoneResponse, error)
	GetAll(ctx context.Context) ([]dto.ParkingZoneAvailabilityResponse, error)
	GetByID(ctx context.Context, id int) (*dto.ParkingZoneAvailabilityResponse, error)
	Update(ctx context.Context, id int, request dto.UpdateParkingZoneRequest) (*dto.ParkingZoneResponse, error)
	Delete(ctx context.Context, id int) error
}

type parkingZoneService struct {
	parkingZoneRepository repository.ParkingZoneRepository
}

func NewParkingZoneService(parkingZoneRepository repository.ParkingZoneRepository) ParkingZoneService {
	return &parkingZoneService{parkingZoneRepository: parkingZoneRepository}
}

func (s *parkingZoneService) Create(ctx context.Context, request dto.CreateParkingZoneRequest) (*dto.ParkingZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          strings.TrimSpace(request.Name),
		Type:          request.Type,
		TotalCapacity: request.TotalCapacity,
		PricePerHour:  request.PricePerHour,
	}

	if zone.Name == "" {
		return nil, apperror.BadRequest("Validation failed", map[string]string{
			"name": "Name is required",
		}, nil)
	}

	if err := s.parkingZoneRepository.Create(ctx, zone); err != nil {
		return nil, err
	}

	response := dto.ParkingZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}

	return &response, nil
}

func (s *parkingZoneService) GetAll(ctx context.Context) ([]dto.ParkingZoneAvailabilityResponse, error) {
	zones, err := s.parkingZoneRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ParkingZoneAvailabilityResponse, 0, len(zones))
	for _, zone := range zones {
		activeCount, err := s.parkingZoneRepository.CountActiveReservations(ctx, zone.ID)
		if err != nil {
			return nil, err
		}

		responses = append(responses, dto.ParkingZoneAvailabilityResponse{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: zone.TotalCapacity - int(activeCount),
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt,
		})
	}

	return responses, nil
}

func (s *parkingZoneService) GetByID(ctx context.Context, id int) (*dto.ParkingZoneAvailabilityResponse, error) {
	zone, err := s.parkingZoneRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	activeCount, err := s.parkingZoneRepository.CountActiveReservations(ctx, zone.ID)
	if err != nil {
		return nil, err
	}

	response := dto.ParkingZoneAvailabilityResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity - int(activeCount),
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}

	return &response, nil
}

func (s *parkingZoneService) Update(ctx context.Context, id int, request dto.UpdateParkingZoneRequest) (*dto.ParkingZoneResponse, error) {
	if request.Name == nil && request.Type == nil && request.TotalCapacity == nil && request.PricePerHour == nil {
		return nil, apperror.BadRequest("Validation failed", map[string]string{
			"body": "At least one field is required",
		}, nil)
	}

	updates := &models.ParkingZone{}
	if request.Name != nil {
		trimmedName := strings.TrimSpace(*request.Name)
		if trimmedName == "" {
			return nil, apperror.BadRequest("Validation failed", map[string]string{
				"name": "Name cannot be empty",
			}, nil)
		}

		updates.Name = trimmedName
	}

	if request.Type != nil {
		updates.Type = *request.Type
	}

	if request.TotalCapacity != nil {
		updates.TotalCapacity = *request.TotalCapacity
	}

	if request.PricePerHour != nil {
		updates.PricePerHour = *request.PricePerHour
	}

	zone, err := s.parkingZoneRepository.Update(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	response := dto.ParkingZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}

	return &response, nil
}

func (s *parkingZoneService) Delete(ctx context.Context, id int) error {
	return s.parkingZoneRepository.Delete(ctx, id)
}
