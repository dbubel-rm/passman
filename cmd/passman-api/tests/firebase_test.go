package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
	l = log.New(ioutil.Discard, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	var err error
	for i := 0; i < 20; i++ {
		d, err = db.New("root@tcp(db:3306)/passman")
		if err != nil {
			log.Printf("\t%s DB connection error: %s", Failed, err.Error())
			time.Sleep(time.Second)
			continue
		}
		log.Printf("\t%s DB connection OK", Success)
		break
	}

	// fixtureFiles, err := filepath.Glob("../../../*.sql")

	// fmt.Println(pwd)
	// fmt.Println("fixtures", fixtureFiles)

	// sql, err := ioutil.ReadFile(fixtureFiles[0])
	// fmt.Println(string(sql))
	// if err != nil {
	// 	panic(err)
	// }

	_, err = d.Database.Exec("truncate table credentials")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

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

	type response struct {
		IdToken string `json:"idToken"`
	}
	var s response

	json.NewDecoder(w.Body).Decode(&s)

	del, _ := json.Marshal(&s)

	// Test create credential
	a = handlers.API(l, d).(*web.App)
	r = httptest.NewRequest("POST", "/v1/credential", strings.NewReader(`{"username":"dean@dean.com","password":"test123","serviceName":"test_service"}`))
	tt := fmt.Sprintf("Bearer %s", s.IdToken)
	r.Header.Set("Authorization", tt)
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		b, _ := ioutil.ReadAll(w.Body)
		fmt.Println(string(b))
		t.Fatalf("\t%s Should receive a status code of 200 for the response. Received: %d", Failed, w.Code)
	}
	t.Logf("\t%s Should receive a status code of 200.", Success)

	// Test get credential
	a = handlers.API(l, d).(*web.App)
	r = httptest.NewRequest("GET", "/v1/credential/test_service", nil)
	r.Header.Set("Authorization", tt)
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)
	// b, _ := ioutil.ReadAll(w.Body)
	// fmt.Println(string(b))
	if w.Code != http.StatusOK {
		t.Logf("\t%s Should receive a status code of 200 for the response. Received: %d", Failed, w.Code)
	}
	t.Logf("\t%s Should receive a status code of 200.", Success)

	// Test delete account

	a = handlers.API(l, d).(*web.App)
	r = httptest.NewRequest("DELETE", "/v1/users", strings.NewReader(string(del)))
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("\t%s Should receive a status code of 200 for the response. Received: %d", Failed, w.Code)
	}
	t.Logf("\t%s Should receive a status code of 200.", Success)

}
