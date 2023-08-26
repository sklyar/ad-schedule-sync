package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/sklyar/go-transact"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/sklyar/ad-schedule-sync/backend/internal/config"
	bookingrepository "github.com/sklyar/ad-schedule-sync/backend/internal/repository/booking"
	"github.com/sklyar/ad-schedule-sync/backend/internal/server/http"
	bookingservice "github.com/sklyar/ad-schedule-sync/backend/internal/service/booking"

	"github.com/sklyar/go-transact/adapters/txstd"
)

const (
	defaultLogLevel = slog.LevelDebug
)

func main() {
	ctx := context.Background()
	cfg := config.Config{}

	logger, err := newLogger(cfg.Logging.Level)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	slog.SetDefault(logger)

	sqlDB, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	txManager, db, err := transact.NewManager(txstd.Wrap(sqlDB))
	if err != nil {
		panic(err)
	}

	bookingRepo := bookingrepository.NewRepository(db)
	bookingSrv := bookingservice.NewService(txManager, bookingRepo)

	server := http.NewServer(cfg.Server.HTTP, bookingSrv)
	slog.Info("starting server", slog.String("addr", cfg.Server.HTTP.Addr))
	if err := server.ListenAndServe(ctx); err != nil {
		panic(err)
	}
}

func newLogger(levelStr string) (*slog.Logger, error) {
	level := defaultLogLevel
	if len(levelStr) > 0 {
		if err := level.UnmarshalText([]byte(levelStr)); err != nil {
			return nil, fmt.Errorf("failed to parse log level: %w", err)
		}
	}
	opts := slog.HandlerOptions{Level: level}
	h := slog.NewTextHandler(os.Stdout, &opts)
	return slog.New(h), nil
}
