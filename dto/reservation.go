package dto

import "time"

type CreateReservationRequest struct {
	ZoneID       int    `json:"zone_id" validate:"required,gt=0"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type ReservationResponse struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	ZoneID       int       `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ReservationZoneSummary struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type MyReservationResponse struct {
	ID           int                    `json:"id"`
	LicensePlate string                 `json:"license_plate"`
	Status       string                 `json:"status"`
	Zone         ReservationZoneSummary `json:"zone"`
	CreatedAt    time.Time              `json:"created_at"`
}

type ReservationUserSummary struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AdminReservationResponse struct {
	ID           int                    `json:"id"`
	UserID       int                    `json:"user_id"`
	ZoneID       int                    `json:"zone_id"`
	LicensePlate string                 `json:"license_plate"`
	Status       string                 `json:"status"`
	User         ReservationUserSummary `json:"user"`
	Zone         ReservationZoneSummary `json:"zone"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}
