package handlers

import (
	"log"
	"net/http"

	"github.com/dbubel/passman/internal/mid"
	"github.com/dbubel/passman/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(log *log.Logger) http.Handler {

	app := web.New(log, mid.ErrorHandler, mid.RequestLogger)

	// Register health check endpoint. This route is not authenticated.
	// check := Check{
	// 	MasterDB: masterDB,
	// }
	// app.Handle("GET", "/v1/health", check.Health)

	// Register user management and authentication endpoints.
	u := User{
	// MasterDB:       masterDB,
	// TokenGenerator: authenticator,
	}
	app.Handle("GET", "/v1/users", u.List)
	// app.Handle("POST", "/v1/users", u.Create, authmw.HasRole(auth.RoleAdmin), authmw.Authenticate)
	// app.Handle("GET", "/v1/users/:id", u.Retrieve, authmw.Authenticate)
	// app.Handle("PUT", "/v1/users/:id", u.Update, authmw.HasRole(auth.RoleAdmin), authmw.Authenticate)
	// app.Handle("DELETE", "/v1/users/:id", u.Delete, authmw.HasRole(auth.RoleAdmin), authmw.Authenticate)

	// This route is not authenticated
	app.Handle("GET", "/v1/users/token", u.Token)

	return app
}
