package models

import "time"

type Reservation struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	UserID       uint        `gorm:"not null;index" json:"user_id"`
	ZoneID       uint        `gorm:"not null;index" json:"zone_id"`
	LicensePlate string      `gorm:"type:varchar(15);not null;index:idx_active_license_plate,unique,where:status = 'active'" json:"license_plate"`
	Status       string      `gorm:"type:varchar(20);not null;default:active;index;check:status IN ('active','completed','cancelled')" json:"status"`
	User         User        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Zone         ParkingZone `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
