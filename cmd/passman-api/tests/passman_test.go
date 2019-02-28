package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dbubel/passman/cmd/passman-api/handlers"
	"github.com/dbubel/passman/internal/mid"
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
var f web.Middleware

func init() {
	f = mid.FakeAuth
	l = log.New(ioutil.Discard, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	var err error
	var dsn = "root@tcp(db:3306)/passman"
	if os.Getenv("DB_HOST") != "" {
		dsn = os.Getenv("DB_HOST")
	}
	for i := 0; i < 20; i++ {
		d, err = db.New(dsn)
		if err != nil {
			log.Println(dsn)
			log.Printf("\t%s DB connection error: %s", Failed, err.Error())
			time.Sleep(time.Second)
			continue
		}
		log.Printf("\t%s DB connection OK", Success)
		break
	}

	fixtureFiles, err := filepath.Glob("../../../*.sql")

	fmt.Println("fixtures", fixtureFiles)

	sql, err := ioutil.ReadFile(fixtureFiles[0])
	fmt.Println(string(sql))
	if err != nil {
		panic(err)
	}

	_, err = d.Database.Exec("truncate table credentials")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

}

func TestPassman(t *testing.T) {

	// Test bad request

	a := handlers.API(l, d, f).(*web.App)
	r := httptest.NewRequest("POST", "/v1/users", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		b, _ := ioutil.ReadAll(w.Body)
		fmt.Println(string(b))
		t.Fatalf("\t%s Create user bad request failed.", Failed)
	}
	t.Logf("\t%s Create user bad request.", Success)

	// Test create user
	a = handlers.API(l, d, f).(*web.App)

	r = httptest.NewRequest("POST", "/v1/users", strings.NewReader(`{"email":"dean@dean.com","password":"test123","returnSecureToken":true}`))
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		b, _ := ioutil.ReadAll(w.Body)
		fmt.Println(string(b))
		t.Fatalf("\t%s Create user failed.", Failed)
	}
	t.Logf("\t%s Create user.", Success)

	type response struct {
		IdToken string `json:"idToken"`
	}

	var s response
	var del []byte
	var tt string

	// Test signin
	a = handlers.API(l, d, f).(*web.App)
	r = httptest.NewRequest("GET", "/v1/signin", strings.NewReader(`{"email":"dean@dean.com","password":"test123","returnSecureToken":true}`))
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("\t%s Signin failed.", Failed)
	}
	t.Logf("\t%s Signin.", Success)

	json.NewDecoder(w.Body).Decode(&s)

	del, _ = json.Marshal(&s)

	// Test create credential
	a = handlers.API(l, d, f).(*web.App)
	r = httptest.NewRequest("POST", "/v1/credential", strings.NewReader(`{"username":"dean@dean.com","password":"test123","serviceName":"test_service"}`))
	tt = fmt.Sprintf("Bearer %s", s.IdToken)
	r.Header.Set("Authorization", tt)
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		b, _ := ioutil.ReadAll(w.Body)
		fmt.Println(string(b))
		t.Fatalf("\t%s Create credential failed.", Failed)
	}
	t.Logf("\t%s Create credential.", Success)

	// Test get credential
	a = handlers.API(l, d, f).(*web.App)
	r = httptest.NewRequest("GET", "/v1/credential/test_service", nil)
	r.Header.Set("Authorization", tt)
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		b, _ := ioutil.ReadAll(w.Body)
		fmt.Println(string(b))
		t.Fatalf("\t%s Get credential failed.", Failed)
	}

	t.Logf("\t%s Get credential.", Success)

	// Test get services
	a = handlers.API(l, d, f).(*web.App)
	r = httptest.NewRequest("GET", "/v1/services", nil)
	r.Header.Set("Authorization", tt)
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)
	p, _ := ioutil.ReadAll(w.Body)
	fmt.Println(string(p))
	if w.Code != http.StatusOK {
		b, _ := ioutil.ReadAll(w.Body)
		fmt.Println(string(b))
		t.Fatalf("\t%s Get services failed.", Failed)
	}

	t.Logf("\t%s Get services.", Success)

	// Test update credential
	a = handlers.API(l, d, f).(*web.App)
	r = httptest.NewRequest("POST", "/v1/credential/update", strings.NewReader(`{"password":"test1235","serviceName":"test_service"}`))
	r.Header.Set("Authorization", tt)
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		b, _ := ioutil.ReadAll(w.Body)
		fmt.Println("HERE", w.Code, string(b))
		t.Fatalf("\t%s Update credential failed.", Failed)
	}

	a = handlers.API(l, d, f).(*web.App)
	r = httptest.NewRequest("GET", "/v1/credential/test_service", nil)
	r.Header.Set("Authorization", tt)
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)
	b, _ := ioutil.ReadAll(w.Body)
	if !strings.Contains(string(b), "test1235") {
		fmt.Println(string(b))
		fmt.Println("HERE", w.Code, string(b))
		t.Fatalf("\t%s Update credential failed.", Failed)
	}
	t.Logf("\t%s Update credential.", Success)

	// Test delete credential
	a = handlers.API(l, d, f).(*web.App)
	r = httptest.NewRequest("DELETE", "/v1/credential/test_service", nil)
	r.Header.Set("Authorization", tt)
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		b, _ := ioutil.ReadAll(w.Body)
		fmt.Println(string(b))
		t.Fatalf("\t%s Delete credential failed.", Failed)
	}

	t.Logf("\t%s Delete credential.", Success)

	// Test delete account
	a = handlers.API(l, d, f).(*web.App)
	r = httptest.NewRequest("DELETE", "/v1/users", strings.NewReader(string(del)))
	r.Header.Set("Authorization", tt)
	w = httptest.NewRecorder()

	a.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		b, _ := ioutil.ReadAll(w.Body)
		fmt.Println(string(b))
		t.Fatalf("\t%s Delete account failed.", Failed)
	}
	t.Logf("\t%s Delete account.", Success)

}
