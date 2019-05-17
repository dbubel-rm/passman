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
	var app = web.New(log, mid.RequestLogger)

	check := Check{
		MasterDB: db.Database,
	}
	app.Router.NotFound = goAway(log)

	app.Handle("GET", "/health", check.Health)

	var apiKey = "AIzaSyBItfzjx74wXWCet-ARldNNpKIZVR1PQ5I%0A"
	var firebaseBaseURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty"
	firebase := Firebase{
		MasterDB:          db.Database,
		SigninURL:         fmt.Sprintf("%s/verifyPassword?key=%s", firebaseBaseURL, apiKey),
		CreateURL:         fmt.Sprintf("%s/signupNewUser?key=%s", firebaseBaseURL, apiKey),
		DeleteURL:         fmt.Sprintf("%s/deleteAccount?key=%s", firebaseBaseURL, apiKey),
		VerifyURL:         fmt.Sprintf("%s/getOobConfirmationCode?key=%s", firebaseBaseURL, apiKey),
		ChangePasswordURL: fmt.Sprintf("%s/setAccountInfo?key=%s", firebaseBaseURL, apiKey),
	}

	// TODO: update account password
	app.Handle("POST", "/v1/users", firebase.Create)
	app.Handle("POST", "/v1/users/verify", firebase.Verify)
	app.Handle("POST", "/v1/users/password", firebase.ChangePassword, auth)
	app.Handle("DELETE", "/v1/users", firebase.Delete, auth)
	app.Handle("GET", "/v1/signin", firebase.Signin)

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

	// app.Handle("OPTIONS", "/v1/*name", creds.cors)
	// TODO: get credentials by ID

	return app
}

func goAway(l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s -> %d -> %s -> %s", r.Method, r.ContentLength, r.RemoteAddr, r.RequestURI)
		w.WriteHeader(http.StatusNotFound)
	})
}
