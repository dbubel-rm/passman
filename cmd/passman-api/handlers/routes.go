package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dbubel/passman/internal/mid"
	"github.com/dbubel/passman/internal/platform/db"
	"github.com/dbubel/passman/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(log *log.Logger, db *db.MySQLDB, auth web.Middleware) http.Handler {

	app := web.New(log, mid.RequestLogger)

	check := Check{
		MasterDB: db.Database,
	}
	app.Handle("GET", "/v1/health", check.Health)

	apiKey := "AIzaSyBItfzjx74wXWCet-ARldNNpKIZVR1PQ5I%0A"
	f := Firebase{
		SigninURL: fmt.Sprintf("https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=%s", apiKey),
	}

	// TODO: update account password
	app.Handle("POST", "/v1/users", f.Create)
	app.Handle("POST", "/v1/users/verify", f.Verify)
	app.Handle("DELETE", "/v1/users", f.Delete)
	app.Handle("GET", "/v1/signin", f.Signin)

	creds := Credentials{
		MasterDB: db.Database,
	}

	// TODO: make a time since password rotation field
	// store a new credential

	app.Handle("POST", "/v1/credential", creds.add, auth)
	app.Handle("GET", "/v1/credential/:serviceName", creds.get, auth)
	app.Handle("DELETE", "/v1/credential/:serviceName", creds.delete, auth)
	app.Handle("POST", "/v1/credential/update", creds.update, auth)
	app.Handle("GET", "/v1/services", creds.services, auth)
	// TODO: get credentials by ID

	return app
}

// func translate(err error) error {
// 	switch errors.Cause(err) {
// 	case user.ErrNotFound:
// 		return web.ErrNotFound
// 	case user.ErrInvalidID:
// 		return web.ErrInvalidID
// 	case user.ErrAuthenticationFailure:
// 		return web.ErrUnauthorized
// 	case user.ErrForbidden:
// 		return web.ErrForbidden
// 	}
// 	return err
// }
