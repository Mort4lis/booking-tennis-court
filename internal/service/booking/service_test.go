package booking_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Mort4lis/booking-tennis-court/internal/dto"
	"github.com/Mort4lis/booking-tennis-court/internal/service/booking"
)

func TestService_CreateBooking_NotImplemented(t *testing.T) {
	t.Parallel()

	svc := booking.NewService(nil)

	_, err := svc.CreateBooking(context.Background(), dto.CreateBookingRequest{})
	require.Error(t, err)
}
