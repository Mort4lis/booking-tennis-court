package dto

import (
	"time"
)

// CreateBookingRequest is the transport-level payload to create a booking.
type CreateBookingRequest struct {
	CourtID    int64
	TelegramID int64
	StartsAt   time.Time
	EndsAt     time.Time
}

// BookingResponse is the transport-level representation of a booking.
type BookingResponse struct {
	ID       int64
	CourtID  int64
	UserID   int64
	StartsAt time.Time
	EndsAt   time.Time
	Status   string
}
