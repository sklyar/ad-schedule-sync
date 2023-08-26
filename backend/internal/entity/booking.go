package entity

import "time"

type Booking struct {
	ID         uint32
	ClientName string
	BookingAt  time.Time
	VKPostID   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}
