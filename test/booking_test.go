//go:build integration

package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/Mort4lis/booking-tennis-court/internal/dto"
	bookingrepo "github.com/Mort4lis/booking-tennis-court/internal/repository/postgres"
	bookingsvc "github.com/Mort4lis/booking-tennis-court/internal/service/booking"
	"github.com/Mort4lis/booking-tennis-court/test/testutil"
)

// bookingFixturesDir holds the YAML fixtures loaded before every test method
// in BookingSuite.
const bookingFixturesDir = "testdata/booking"

// BookingSuite exercises service/booking.Service wired against the real
// repository/postgres.BookingRepository — no mocks — so it proves the wiring
// (Service -> Repository -> PostgreSQL) actually works.
type BookingSuite struct {
	testutil.Suite

	service *bookingsvc.Service
}

func TestBookingSuite(t *testing.T) {
	// Safe to parallelize at the suite level: every suite gets its own
	// database (see SetupSuite), so two suites never touch the same rows.
	// Test methods within a single suite are not parallelized — they share
	// that one database and rely on SetupTest reloading fixtures between
	// them for isolation.
	t.Parallel()
	suite.Run(t, new(BookingSuite))
}

func (s *BookingSuite) SetupSuite() {
	s.SetupDatabase(context.Background(), pg.DSN, "booking_test")

	repo := bookingrepo.NewBookingRepository(s.Pool)
	s.service = bookingsvc.NewService(repo)
}

func (s *BookingSuite) SetupTest() {
	s.LoadFixtures(bookingFixturesDir)
}

func (s *BookingSuite) TestCreateBooking_StubNotImplemented() {
	ctx := context.Background()

	req := dto.CreateBookingRequest{
		CourtID:    1,
		TelegramID: 1001,
		StartsAt:   time.Now().Add(24 * time.Hour),
		EndsAt:     time.Now().Add(25 * time.Hour),
	}

	_, err := s.service.CreateBooking(ctx, req)
	s.ErrorContains(err, "not implemented")
}

func (s *BookingSuite) TestFixtures_SeedDataIsLoaded() {
	ctx := context.Background()

	var courtCount int

	err := s.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM courts").Scan(&courtCount)
	s.Require().NoError(err)
	s.Equal(2, courtCount)

	var userCount int

	err = s.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
	s.Require().NoError(err)
	s.Equal(2, userCount)
}
