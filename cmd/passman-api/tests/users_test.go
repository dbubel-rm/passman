package tests

import (
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

var d *db.DB
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
func TestUsersCreate(t *testing.T) {

	a = handlers.API(l, d).(*web.App)
	r := httptest.NewRequest("POST", "/v1/users", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("\t%s Should receive a status code of 400 for the response. Received: %d", Failed, w.Code)
	}
	t.Logf("\t%s Should receive a status code of 400 for the response.", Success)
}
