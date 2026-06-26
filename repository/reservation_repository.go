package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	apperror "spotsync/apperror"
	"spotsync/models"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation *models.Reservation) error
	FindByUserID(ctx context.Context, userID int) ([]models.Reservation, error)
	FindByID(ctx context.Context, id int) (*models.Reservation, error)
	UpdateStatus(ctx context.Context, id int, status string) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(ctx context.Context, reservation *models.Reservation) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, reservation.ZoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("Parking zone not found", nil, err)
			}

			return apperror.Internal("Internal server error", err)
		}

		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", reservation.ZoneID, models.ReservationStatusActive).
			Count(&activeCount).Error; err != nil {
			return apperror.Internal("Internal server error", err)
		}

		if activeCount >= int64(zone.TotalCapacity) {
			return apperror.Conflict("Zone is full", map[string]string{
				"zone_id": "No available spots in this zone",
			}, nil)
		}

		var duplicateCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("license_plate = ? AND status = ?", reservation.LicensePlate, models.ReservationStatusActive).
			Count(&duplicateCount).Error; err != nil {
			return apperror.Internal("Internal server error", err)
		}

		if duplicateCount > 0 {
			return apperror.Conflict("Duplicate active license plate", map[string]string{
				"license_plate": "License plate already has an active reservation",
			}, nil)
		}

		reservation.Status = models.ReservationStatusActive
		if err := tx.Create(reservation).Error; err != nil {
			if isDuplicateActiveLicensePlateError(err) {
				return apperror.Conflict("Duplicate active license plate", map[string]string{
					"license_plate": "License plate already has an active reservation",
				}, err)
			}

			return apperror.Internal("Internal server error", err)
		}

		return nil
	})
}

func (r *reservationRepository) FindByUserID(ctx context.Context, userID int) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.WithContext(ctx).
		Preload("Zone").
		Where("user_id = ?", userID).
		Order("id ASC").
		Find(&reservations).Error; err != nil {
		return nil, apperror.Internal("Internal server error", err)
	}

	return reservations, nil
}

func (r *reservationRepository) FindByID(ctx context.Context, id int) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := r.db.WithContext(ctx).First(&reservation, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("Reservation not found", nil, err)
		}

		return nil, apperror.Internal("Internal server error", err)
	}

	return &reservation, nil
}

func (r *reservationRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	result := r.db.WithContext(ctx).Model(&models.Reservation{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return apperror.Internal("Internal server error", result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NotFound("Reservation not found", nil, nil)
	}

	return nil
}

func isDuplicateActiveLicensePlateError(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate key") && strings.Contains(message, "idx_active_license_plate")
}
