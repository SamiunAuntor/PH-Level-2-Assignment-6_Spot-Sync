package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	apperror "spotsync/apperror"
	"spotsync/models"
)

type ParkingZoneRepository interface {
	Create(ctx context.Context, zone *models.ParkingZone) error
	FindAll(ctx context.Context) ([]models.ParkingZone, error)
	FindByID(ctx context.Context, id int) (*models.ParkingZone, error)
	Update(ctx context.Context, id int, updates *models.ParkingZone) (*models.ParkingZone, error)
	Delete(ctx context.Context, id int) error
	CountActiveReservations(ctx context.Context, zoneID int) (int64, error)
}

type parkingZoneRepository struct {
	db *gorm.DB
}

func NewParkingZoneRepository(db *gorm.DB) ParkingZoneRepository {
	return &parkingZoneRepository{db: db}
}

func (r *parkingZoneRepository) Create(ctx context.Context, zone *models.ParkingZone) error {
	if err := r.db.WithContext(ctx).Create(zone).Error; err != nil {
		return apperror.Internal("Internal server error", err)
	}

	return nil
}

func (r *parkingZoneRepository) FindAll(ctx context.Context) ([]models.ParkingZone, error) {
	var zones []models.ParkingZone
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&zones).Error; err != nil {
		return nil, apperror.Internal("Internal server error", err)
	}

	return zones, nil
}

func (r *parkingZoneRepository) FindByID(ctx context.Context, id int) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	if err := r.db.WithContext(ctx).First(&zone, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("Parking zone not found", nil, err)
		}

		return nil, apperror.Internal("Internal server error", err)
	}

	return &zone, nil
}

func (r *parkingZoneRepository) Update(ctx context.Context, id int, updates *models.ParkingZone) (*models.ParkingZone, error) {
	var updatedZone models.ParkingZone

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("Parking zone not found", nil, err)
			}

			return apperror.Internal("Internal server error", err)
		}

		if updates.Name != "" {
			zone.Name = updates.Name
		}

		if updates.Type != "" {
			zone.Type = updates.Type
		}

		if updates.TotalCapacity != 0 {
			activeCount, err := r.countActiveReservationsTx(ctx, tx, zone.ID)
			if err != nil {
				return err
			}

			if updates.TotalCapacity < int(activeCount) {
				return apperror.Conflict("Total capacity cannot be less than active reservations", nil, nil)
			}

			zone.TotalCapacity = updates.TotalCapacity
		}

		if updates.PricePerHour != 0 {
			zone.PricePerHour = updates.PricePerHour
		}

		if err := tx.Save(&zone).Error; err != nil {
			return apperror.Internal("Internal server error", err)
		}

		updatedZone = zone
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &updatedZone, nil
}

func (r *parkingZoneRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		if err := tx.First(&zone, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("Parking zone not found", nil, err)
			}

			return apperror.Internal("Internal server error", err)
		}

		var reservationCount int64
		if err := tx.Model(&models.Reservation{}).Where("zone_id = ?", id).Count(&reservationCount).Error; err != nil {
			return apperror.Internal("Internal server error", err)
		}

		if reservationCount > 0 {
			return apperror.Conflict("Parking zone cannot be deleted because reservation history exists", nil, nil)
		}

		if err := tx.Delete(&zone).Error; err != nil {
			return apperror.Internal("Internal server error", err)
		}

		return nil
	})
}

func (r *parkingZoneRepository) CountActiveReservations(ctx context.Context, zoneID int) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, models.ReservationStatusActive).
		Count(&count).Error; err != nil {
		return 0, apperror.Internal("Internal server error", err)
	}

	return count, nil
}

func (r *parkingZoneRepository) countActiveReservationsTx(ctx context.Context, tx *gorm.DB, zoneID int) (int64, error) {
	var count int64
	if err := tx.WithContext(ctx).
		Model(&models.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, models.ReservationStatusActive).
		Count(&count).Error; err != nil {
		return 0, apperror.Internal("Internal server error", err)
	}

	return count, nil
}
