package repository

import (
	"context"
	"strings"

	"gorm.io/gorm"

	apperror "spotsync/apperror"
	"spotsync/models"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation *models.Reservation) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(ctx context.Context, reservation *models.Reservation) error {
	if err := r.db.WithContext(ctx).Create(reservation).Error; err != nil {
		if isDuplicateActiveLicensePlateError(err) {
			return apperror.Conflict("Duplicate active license plate", map[string]string{
				"license_plate": "License plate already has an active reservation",
			}, err)
		}

		return apperror.Internal("Internal server error", err)
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
