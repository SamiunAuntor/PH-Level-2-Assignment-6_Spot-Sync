package models

import "time"

type ParkingZone struct {
	ID            int           `gorm:"primaryKey;autoIncrement;type:integer" json:"id"`
	Name          string        `gorm:"not null" json:"name"`
	Type          string        `gorm:"type:varchar(30);not null;check:type IN ('general','ev_charging','covered')" json:"type"`
	TotalCapacity int           `gorm:"not null;check:total_capacity > 0" json:"total_capacity"`
	PricePerHour  float64       `gorm:"type:numeric(10,2);not null;check:price_per_hour > 0" json:"price_per_hour"`
	Reservations  []Reservation `gorm:"foreignKey:ZoneID" json:"-"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
