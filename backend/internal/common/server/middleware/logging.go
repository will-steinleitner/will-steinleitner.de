package middleware

import (
	"log/slog"
	"net/http"
)

func requestLogging(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Incomming request", slog.String("method", r.Method), slog.String("path", r.URL.Path))
	})
}
