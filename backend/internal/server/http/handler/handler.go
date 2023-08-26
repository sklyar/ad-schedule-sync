package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/sklyar/ad-schedule-sync/backend/internal/service"
)

type Handler struct {
	service service.Booking
}

func NewHandler(s service.Booking) Handler {
	return Handler{service: s}
}

type reqSyncBookingsData struct {
	Date     string `json:"date"`
	Bookings []struct {
		ClientName string `json:"client_name"`
		Time       string `json:"time"`
	} `json:"bookings"`
}

func (d reqSyncBookingsData) convert() (service.SyncBookingsData, error) {
	date, err := time.Parse("2006-01-02Z07:00", d.Date)
	if err != nil {
		return service.SyncBookingsData{}, fmt.Errorf("failed to parse date: %w", err)
	}

	bookings := make([]service.SyncBookingData, len(d.Bookings))
	for i, b := range d.Bookings {
		t, err := time.Parse("15:04", b.Time)
		if err != nil {
			return service.SyncBookingsData{}, fmt.Errorf("failed to parse booking time: %w", err)
		}

		bookings[i] = service.SyncBookingData{
			ClientName: b.ClientName,
			BookingAt:  time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, date.Location()),
		}
	}

	return service.SyncBookingsData{
		Date:     date,
		Bookings: bookings,
	}, nil
}

func (h Handler) SyncBookings(w http.ResponseWriter, r *http.Request) {
	req := reqSyncBookingsData{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := req.convert()
	if err != nil {
		slog.Error("failed to convert request data", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.SyncBookings(r.Context(), data); err != nil {
		slog.Error("failed to sync bookings", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
