package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/will-steinleitner/will-steinleitner.de/internal/common/server/middleware"
)

func AddRoute(
	mux *http.ServeMux,
	ctx context.Context,
	logger *slog.Logger,
) {
	defaultMiddle := middleware.NewDefault(logger)

	mux.Handle(http.MethodGet+" /{$}", defaultMiddle(handleRoot()))
	mux.Handle(http.MethodGet+" /healthz", defaultMiddle(handleHealthz()))
}
