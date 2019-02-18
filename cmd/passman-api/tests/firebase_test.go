package tests

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dbubel/passman/cmd/passman-api/handlers"
	"github.com/dbubel/passman/internal/platform/db"
	"github.com/dbubel/passman/internal/platform/web"
)

var a *web.App

const (
	Success = "\u2713"
	Failed  = "\u2717"
)

var d *db.MySQLDB
var l *log.Logger

func init() {
	l = log.New(ioutil.Discard, "PASSMAN : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	var err error
	d, err = db.New("root@tcp(db:3306)/passman")
	if err != nil {
		log.Fatalf("\t%s DB connection error: %s", Failed, err.Error())
	}
	log.Printf("\t%s DB connection OK", Success)
}

func TestUsers(t *testing.T) {

	// Test bad request
	a = handlers.API(l, d).(*web.App)
	r := httptest.NewRequest("POST", "/v1/users", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("\t%s Should receive a status code of 400. Received: %d", Failed, w.Code)
	}
	t.Logf("\t%s Should receive a status code of 400.", Success)

	// Test create
	a = handlers.API(l, d).(*web.App)
	r = httptest.NewRequest("POST", "/v1/users", strings.NewReader(`{"email":"dean@dean.com","password":"test123","returnSecureToken":true}`))
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("\t%s Should receive a status code of 200 for the response. Received: %d", Failed, w.Code)
	}
	t.Logf("\t%s Should receive a status code of 200.", Success)

	// Test signin
	a = handlers.API(l, d).(*web.App)
	r = httptest.NewRequest("GET", "/v1/signin", strings.NewReader(`{"email":"dean@dean.com","password":"test123","returnSecureToken":true}`))
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("\t%s Should receive a status code of 200 for the response. Received: %d", Failed, w.Code)
	}
	t.Logf("\t%s Should receive a status code of 200.", Success)

	// Test delete
	type response struct {
		IdToken string `json:"idToken"`
	}
	var s response

	json.NewDecoder(w.Body).Decode(&s)

	del, _ := json.Marshal(&s)

	a = handlers.API(l, d).(*web.App)
	r = httptest.NewRequest("DELETE", "/v1/users", strings.NewReader(string(del)))
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("\t%s Should receive a status code of 200 for the response. Received: %d", Failed, w.Code)
	}
	t.Logf("\t%s Should receive a status code of 200.", Success)
}
