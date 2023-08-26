package booking

import (
	"context"
	"testing"
	"time"

	"github.com/sklyar/ad-schedule-sync/backend/internal/service"

	"github.com/stretchr/testify/assert"

	"github.com/sklyar/go-transact/txtest"

	"github.com/sklyar/go-transact"
	"github.com/sklyar/go-transact/txsql"

	"github.com/sklyar/ad-schedule-sync/backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateBooking(ctx context.Context, booking entity.Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func (m *mockRepository) GetBookingsByDate(ctx context.Context, date time.Time) ([]entity.Booking, error) {
	args := m.Called(ctx, date)
	if args.Get(0).([]entity.Booking) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.Booking), args.Error(1)
}

func (m *mockRepository) MarkBookingAsCancelled(ctx context.Context, booking entity.Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func TestService_SyncBookings(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type mocks struct {
		db *txtest.DB
		tx *txtest.Tx

		repository *mockRepository
	}

	defaultTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		in      service.SyncBookingsData
		setup   func(*mocks)
		wantErr bool
	}{
		{
			name: "all new bookings",
			in: service.SyncBookingsData{
				Date: makeDate("2020-01-01+05:00"),
				Bookings: []service.SyncBookingData{
					{ClientName: "client1", BookingAt: makeDateTime("2020-01-01 12:00+05:00")},
					{ClientName: "client2", BookingAt: makeDateTime("2020-01-01 14:00+05:00")},
				},
			},
			setup: func(m *mocks) {
				ctx := txtest.WithContext(ctx)
				m.db.On("Begin", ctx, (*txsql.TxOptions)(nil)).Return(m.tx, nil)

				m.repository.
					On("GetBookingsByDate", mock.Anything, makeDate("2020-01-01+05:00")).
					Return([]entity.Booking{}, nil)

				m.repository.
					On("CreateBooking", mock.Anything, entity.Booking{
						ClientName: "client1",
						BookingAt:  makeDateTime("2020-01-01 12:00+05:00"),
						CreatedAt:  defaultTime,
						UpdatedAt:  defaultTime,
					}).Return(nil)

				m.repository.
					On("CreateBooking", mock.Anything, entity.Booking{
						ClientName: "client2",
						BookingAt:  makeDateTime("2020-01-01 14:00+05:00"),
						CreatedAt:  defaultTime,
						UpdatedAt:  defaultTime,
					}).Return(nil)

				m.tx.On("Commit", ctx).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := mocks{
				db:         txtest.NewDB(t),
				tx:         txtest.NewTx(t),
				repository: &mockRepository{},
			}
			if tt.setup != nil {
				tt.setup(&mocks)
			}

			adapter := func(_ transact.TransactionStore) (txsql.DB, error) { return mocks.db, nil }
			txManager, _, err := transact.NewManager(adapter)
			assert.NoError(t, err)

			svc := NewService(txManager, mocks.repository)
			svc.timeNow = func() time.Time { return defaultTime }

			err = svc.SyncBookings(context.Background(), tt.in)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func makeDate(s string) time.Time {
	t, err := time.Parse("2006-01-02Z07:00", s)
	if err != nil {
		panic(err)
	}

	return t
}

func makeTime(s string) time.Time {
	t, err := time.Parse("15:04", s)
	if err != nil {
		panic(err)
	}

	return t
}

func makeDateTime(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04+05:00", s)
	if err != nil {
		panic(err)
	}

	return t
}
