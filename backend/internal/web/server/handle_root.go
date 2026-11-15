package server

import (
	"net/http"

	"github.com/will-steinleitner/will-steinleitner.de/internal/common/server"
)

func handleRoot() http.Handler {
	type RootResponse struct {
		Message string `json:"message"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			resp := RootResponse{
				Message: "Welcome to will-steinleitner.de",
			}

			if err := server.Encode(w, r, http.StatusOK, resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		},
	)
}
