package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/cybre/salesforge-assignment/internal/sequence"
	"github.com/cybre/salesforge-assignment/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SequenceService represents the service layer for sequences.
type SequenceService interface {
	CreateSequence(ctx context.Context, seq sequence.Sequence) error
	PatchSequence(ctx context.Context, patch sequence.SequencePatch) error
	GetSequence(ctx context.Context, id int) (sequence.Sequence, error)
	UpdateStep(ctx context.Context, step sequence.Step) error
	DeleteStep(ctx context.Context, id int) error
}

// Server contains the REST endpoints.
type Server struct {
	sequenceService SequenceService
}

// NewServer creates a new server.
func NewServer(sequenceService SequenceService) *Server {
	return &Server{
		sequenceService: sequenceService,
	}
}

// Start starts the HTTP server and closes it when the context is done.
func (s Server) Start(ctx context.Context, port string) error {
	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:       true,
		LogMethod:       true,
		LogURI:          true,
		LogError:        true,
		LogResponseSize: true,
		HandleError:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger := logging.FromContext(c.Request().Context())

			if v.Error == nil {
				logger.LogAttrs(ctx, slog.LevelInfo, "REQUEST",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.Int64("response_size", v.ResponseSize),
				)
			} else {
				logger.LogAttrs(ctx, slog.LevelError, "REQUEST_ERROR",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Use(middleware.Recover())

	s.RegisterRoutes(e)

	go func() {
		if err := e.Start(":" + port); err != nil {
			if err == http.ErrServerClosed {
				return
			}

			panic(err)
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}

// RegisterRoutes registers the REST endpoints.
func (s Server) RegisterRoutes(e *echo.Echo) {
	e.POST("/sequence", s.CreateSequence)
	e.PATCH("/sequence/:id", s.PatchSequence)
	e.GET("/sequence/:id", s.GetSequence)
	e.PUT("/step/:id", s.UpdateStep)
	e.DELETE("/step/:id", s.DeleteStep)
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}
