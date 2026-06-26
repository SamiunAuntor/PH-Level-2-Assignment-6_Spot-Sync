package models

import "time"

type User struct {
	ID           int           `gorm:"primaryKey;autoIncrement;type:integer" json:"id"`
	Name         string        `gorm:"not null" json:"name"`
	Email        string        `gorm:"not null;uniqueIndex" json:"email"`
	Password     string        `gorm:"not null" json:"-"`
	Role         string        `gorm:"type:varchar(20);not null;default:driver;check:role IN ('driver','admin')" json:"role"`
	Reservations []Reservation `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}
