package handlers

import (
	"log"
	"net/http"

	"github.com/dbubel/passman/internal/mid"
	"github.com/dbubel/passman/internal/platform/db"
	"github.com/dbubel/passman/internal/platform/web"
	"github.com/dbubel/passman/internal/user"
	"github.com/pkg/errors"
)

// API returns a handler for a set of routes.
func API(log *log.Logger, db *db.MySQLDB) http.Handler {
	app := web.New(log, mid.RequestLogger)

	// Register health check endpoint. This route is not authenticated.
	check := Check{
		MasterDB: db.Database,
	}
	app.Handle("GET", "/v1/health", check.Health)

	// users := User{
	// 	MasterDB: db.Database,
	// }
	f := Firebase{}

	app.Handle("POST", "/v1/users", f.Create)
	app.Handle("DELETE", "/v1/users", f.Delete)
	app.Handle("GET", "/v1/signin", f.Signin)

	creds := Credentials{
		MasterDB: db.Database,
	}

	app.Handle("POST", "/v1/credential", creds.Add)

	return app
}

func translate(err error) error {
	switch errors.Cause(err) {
	case user.ErrNotFound:
		return web.ErrNotFound
	case user.ErrInvalidID:
		return web.ErrInvalidID
	case user.ErrAuthenticationFailure:
		return web.ErrUnauthorized
	case user.ErrForbidden:
		return web.ErrForbidden
	}
	return err
}
