package handlers

import (
	"log"
	"net/http"

	"github.com/dbubel/passman/internal/mid"
	"github.com/dbubel/passman/internal/platform/db"
	"github.com/dbubel/passman/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(log *log.Logger, db *db.MySQLDB) http.Handler {

	app := web.New(log, mid.RequestLogger)

	check := Check{
		MasterDB: db.Database,
	}
	app.Handle("GET", "/v1/health", check.Health)

	f := Firebase{}

	// TODO: update account password
	app.Handle("POST", "/v1/users", f.Create)
	app.Handle("DELETE", "/v1/users", f.Delete)
	app.Handle("GET", "/v1/signin", f.Signin)

	creds := Credentials{
		MasterDB: db.Database,
	}

	// TODO: make a time since password rotation field
	// store a new credential
	app.Handle("POST", "/v1/credential", creds.add, mid.AuthHandler)
	//get a credential
	app.Handle("GET", "/v1/credential/:serviceName", creds.get, mid.AuthHandler)
	// delete a credential
	app.Handle("DELETE", "/v1/credential/:serviceName", creds.delete, mid.AuthHandler)
	app.Handle("POST", "/v1/credential", creds.update, mid.AuthHandler)
	// TODO: get service names
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
