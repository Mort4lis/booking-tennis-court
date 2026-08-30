package entity

import (
	"time"
)

// Booking is a reservation of a court for a time slot by a user.
type Booking struct {
	ID        int64
	CourtID   int64
	UserID    int64
	StartsAt  time.Time
	EndsAt    time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
