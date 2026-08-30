package booking

import (
	"context"
	"errors"

	"github.com/Mort4lis/booking-tennis-court/internal/dto"
	"github.com/Mort4lis/booking-tennis-court/internal/entity"
)

// errNotImplemented is returned by every stubbed use-case.
var errNotImplemented = errors.New("not implemented")

// Repository is the persistence port the service depends on.
type Repository interface {
	Create(ctx context.Context, b entity.Booking) (entity.Booking, error)
}

// Service orchestrates booking use-cases on top of a Repository.
type Service struct {
	repo Repository
}

// NewService builds the booking service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateBooking is the single use-case stub. The wiring
// (transport -> service -> repository) is in place; the validation, user
// resolution, and slot checks are the extension point.
func (s *Service) CreateBooking(_ context.Context, _ dto.CreateBookingRequest) (dto.BookingResponse, error) {
	return dto.BookingResponse{}, errNotImplemented
}
