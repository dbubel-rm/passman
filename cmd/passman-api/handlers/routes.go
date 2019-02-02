package handlers

import (
	"log"
	"net/http"

	"github.com/dbubel/passman/internal/mid"
	"github.com/dbubel/passman/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(log *log.Logger) http.Handler {
	app := web.New(log, mid.RequestLogger)

	// Register health check endpoint. This route is not authenticated.
	check := Check{}
	app.Handle("GET", "/v1/health", check.Health)

	return app
}
