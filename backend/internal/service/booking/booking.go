package booking

import (
	"context"
	"fmt"
	"github.com/sklyar/ad-schedule-sync/backend/internal/entity"
	"github.com/sklyar/ad-schedule-sync/backend/internal/service"
	"github.com/sklyar/go-transact"
	"time"
)

type repository interface {
	CreateBooking(ctx context.Context, booking entity.Booking) error
	GetBookingsByDate(ctx context.Context, date time.Time) ([]entity.Booking, error)
	MarkBookingAsCancelled(ctx context.Context, booking entity.Booking) error
}

type bookingKey struct {
	clientName string
	bookingAt  time.Time
}

type bookingsMap map[bookingKey]entity.Booking

// Service is the booking service.
type Service struct {
	r repository

	txManager *transact.Manager
}

// NewService creates a new booking service.
func NewService(txManager *transact.Manager, r repository) *Service {
	return &Service{r: r, txManager: txManager}
}

func (s *Service) SyncBookings(ctx context.Context, data service.SyncBookingsData) error {
	return s.txManager.BeginFunc(ctx, func(ctx context.Context) error {
		return s.syncBookings(ctx, data)
	})
}

func (s *Service) syncBookings(ctx context.Context, data service.SyncBookingsData) error {
	existingBookings, err := s.getBookingsByDate(ctx, data.Date)
	if err != nil {
		return fmt.Errorf("failed to get bookings by date: %w", err)
	}

	if err := s.processNewBookings(ctx, existingBookings, data.Bookings); err != nil {
		return fmt.Errorf("failed to process new bookings: %w", err)
	}

	if err := s.markCancelledBookings(ctx, existingBookings, data.Bookings); err != nil {
		return fmt.Errorf("failed to mark cancelled bookings: %w", err)
	}

	return nil
}

func (s *Service) getBookingsByDate(ctx context.Context, date time.Time) (bookingsMap, error) {
	bookings, err := s.r.GetBookingsByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	bookingMap := make(bookingsMap)
	for _, b := range bookings {
		k := bookingKey{clientName: b.ClientName, bookingAt: b.BookingAt}
		bookingMap[k] = b
	}

	return bookingMap, nil
}

func (s *Service) processNewBookings(ctx context.Context, bookings bookingsMap, newBookings []service.SyncBookingData) error {
	for _, newBooking := range newBookings {
		key := bookingKey{clientName: newBooking.ClientName, bookingAt: newBooking.BookingAt}
		if _, exists := bookings[key]; exists {
			continue
		}

		booking := entity.Booking{
			ClientName: newBooking.ClientName,
			BookingAt:  newBooking.BookingAt,
		}
		if err := s.r.CreateBooking(ctx, booking); err != nil {
			return fmt.Errorf("failed to create booking: %w", err)
		}
	}

	return nil
}

func (s *Service) markCancelledBookings(ctx context.Context, bookings bookingsMap, newBookings []service.SyncBookingData) error {
	for _, booking := range bookings {
		var found bool
		for _, newBooking := range newBookings {
			if booking.ClientName == newBooking.ClientName && booking.BookingAt.Equal(newBooking.BookingAt) {
				found = true
				break
			}
		}

		if !found {
			if err := s.r.MarkBookingAsCancelled(ctx, booking); err != nil {
				return fmt.Errorf("failed to mark booking as cancelled: %w", err)
			}
		}
	}

	return nil
}
