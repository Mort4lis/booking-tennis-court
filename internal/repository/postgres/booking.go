package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Mort4lis/booking-tennis-court/internal/entity"
)

// errNotImplemented is returned by every stubbed repository method.
var errNotImplemented = errors.New("not implemented")

// BookingRepository is the PostgreSQL-backed booking store. It holds the connection
// pool so real queries can be added against it.
type BookingRepository struct {
	pool *pgxpool.Pool
}

// NewBookingRepository builds a PostgreSQL-backed booking repository.
func NewBookingRepository(pool *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{pool: pool}
}

// Create persists a new booking and returns the stored row.
func (r *BookingRepository) Create(_ context.Context, _ entity.Booking) (entity.Booking, error) {
	return entity.Booking{}, errNotImplemented
}
