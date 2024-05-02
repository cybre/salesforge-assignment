package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/cybre/salesforge-assignment/internal/config"
	"github.com/cybre/salesforge-assignment/internal/database"
	"github.com/cybre/salesforge-assignment/internal/sequence"
	"github.com/cybre/salesforge-assignment/internal/transport/http"
	"github.com/cybre/salesforge-assignment/pkg/logging"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	config := config.LoadConfig()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("service", "server")
	slog.SetDefault(logger)

	ctx = logging.WithLogger(ctx, logger)

	database, err := database.NewPostgresDB(config.Database)
	if err != nil {
		log.Fatalf(err.Error())
	}

	sequenceRepo := sequence.NewPostgresRepository(database)
	sequenceService := sequence.NewService(sequenceRepo)
	server := http.NewServer(sequenceService)

	if err := server.Start(ctx, config.Port); err != nil {
		log.Fatalf(err.Error())
	}
}
