package entity

import "time"

// User is a Telegram user who can create bookings.
type User struct {
	ID         int64
	TelegramID int64
	Username   string
	FullName   string
	CreatedAt  time.Time
}
