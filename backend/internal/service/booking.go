package service

import (
	"context"
	"time"
)

// Booking defines the interface for booking synchronization service.
type Booking interface {
	// SyncBookings synchronizes the bookings with the given data.
	SyncBookings(ctx context.Context, data SyncBookingsData) error
}

// SyncBookingData represents the booking data for a single booking.
type SyncBookingData struct {
	ClientName string
	BookingAt  time.Time
}

// SyncBookingsData represents the data needed for synchronizing
// the bookings for a particular date.
type SyncBookingsData struct {
	Date     time.Time
	Bookings []SyncBookingData
}
