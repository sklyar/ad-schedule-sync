package booking

import (
	"context"
	"github.com/sklyar/ad-schedule-sync/backend/internal/entity"
	"github.com/sklyar/go-transact/txsql"
	"time"
)

type Repository struct {
	db txsql.DB
}

func NewRepository(db txsql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateBooking(ctx context.Context, booking entity.Booking) error {
	panic("implement me")
}

func (r *Repository) GetBookingsByDate(ctx context.Context, date time.Time) ([]entity.Booking, error) {
	panic("implement me")
}

func (r *Repository) MarkBookingAsCancelled(ctx context.Context, booking entity.Booking) error {
	panic("implement me")
}
