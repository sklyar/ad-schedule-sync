package handler_test

import (
	"context"
	"io"
	"log/slog"
	httpstd "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sklyar/ad-schedule-sync/backend/internal/service"

	"github.com/stretchr/testify/mock"

	"github.com/sklyar/ad-schedule-sync/backend/internal/config"
	"github.com/sklyar/ad-schedule-sync/backend/internal/server/http"
)

type mockBookingService struct {
	mock.Mock
}

func (m *mockBookingService) SyncBookings(ctx context.Context, data service.SyncBookingsData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func TestHandler_SyncBookings(t *testing.T) {
	t.Parallel()

	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	secretKey := "secret"

	tests := []struct {
		name string

		setup func(t *testing.T, svc *mockBookingService)

		secretKey string
		body      string

		wantStatus int
	}{
		{
			name: "valid request",
			setup: func(t *testing.T, svc *mockBookingService) {
				t.Helper()

				date, err := time.Parse("2006-01-02Z07:00", "2020-01-01+05:00")
				require.NoError(t, err)

				t1 := time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, date.Location())
				t2 := time.Date(date.Year(), date.Month(), date.Day(), 14, 0, 0, 0, date.Location())

				data := service.SyncBookingsData{
					Date: date,
					Bookings: []service.SyncBookingData{
						{ClientName: "Alice", BookingAt: t1},
						{ClientName: "Bob", BookingAt: t2},
					},
				}

				svc.On("SyncBookings", mock.Anything, data).Return(nil)
			},
			secretKey: secretKey,
			body: `{
				"date": "2020-01-01+05:00",
				"bookings": [
					{"client_name": "Alice", "time": "12:00"},
					{"client_name": "Bob", "time": "14:00"}
				]
			}`,
			wantStatus: httpstd.StatusOK,
		},
		{
			name: "valid request with empty bookings",
			setup: func(t *testing.T, svc *mockBookingService) {
				t.Helper()

				date, err := time.Parse("2006-01-02Z07:00", "2020-01-01+05:00")
				require.NoError(t, err)

				data := service.SyncBookingsData{
					Date:     date,
					Bookings: []service.SyncBookingData{},
				}

				svc.On("SyncBookings", mock.Anything, data).Return(nil)
			},
			secretKey: secretKey,
			body: `{
				"date": "2020-01-01+05:00",
				"bookings": []
			}`,
			wantStatus: httpstd.StatusOK,
		},
		{
			name:      "wrong secret key",
			secretKey: "wrong",
			body: `{
				"date": "2020-01-01+05:00",
				"bookings": []
			}`,
			wantStatus: httpstd.StatusUnauthorized,
		},
		{
			name:       "invalid body",
			secretKey:  secretKey,
			body:       `{`,
			wantStatus: httpstd.StatusBadRequest,
		},
		{
			name:      "bad date",
			secretKey: secretKey,
			body: `{
				"date": "2020-01-01",
				"bookings": []
			}`,
			wantStatus: httpstd.StatusBadRequest,
		},
		{
			name:      "bad time",
			secretKey: secretKey,
			body: `{
				"date": "2020-01-01+05:00",
				"bookings": [
					{"client_name": "Alice", "time": "12:00:00"}
				]
			}`,
			wantStatus: httpstd.StatusBadRequest,
		},
		{
			name: "internal error",
			setup: func(t *testing.T, svc *mockBookingService) {
				t.Helper()

				svc.On("SyncBookings", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			secretKey: secretKey,
			body: `{
				"date": "2020-01-01+05:00",
				"bookings": []
			}`,
			wantStatus: httpstd.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := &mockBookingService{}
			if tt.setup != nil {
				tt.setup(t, svc)
			}

			cfg := config.ServerHTTP{SecretKey: secretKey}
			server := http.NewServer(cfg, svc)

			req, err := httpstd.NewRequest(httpstd.MethodPost, "/sync", strings.NewReader(tt.body))
			require.NoError(t, err)
			req.Header.Set("X-Secret-Key", tt.secretKey)

			rr := httptest.NewRecorder()
			server.Router.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}
