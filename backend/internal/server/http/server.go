package http

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sklyar/ad-schedule-sync/backend/internal/config"
	"github.com/sklyar/ad-schedule-sync/backend/internal/server/http/handler"
	"github.com/sklyar/ad-schedule-sync/backend/internal/service"
)

type Server struct {
	Router *mux.Router

	cfg config.ServerHTTP
}

func NewServer(cfg config.ServerHTTP, bookingService service.Booking) *Server {
	h := handler.NewHandler(bookingService)

	r := mux.NewRouter()
	r.Use(checkSecretKeyMiddleware(cfg.SecretKey))
	r.HandleFunc("/sync", h.SyncBookings).Methods(http.MethodPost)

	return &Server{
		Router: r,
		cfg:    cfg,
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: s.Router,
	}

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(ctx)
	}()

	return server.ListenAndServe()
}
