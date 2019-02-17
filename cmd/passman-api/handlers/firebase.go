package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dbubel/passman/internal/firebase"
	"github.com/dbubel/passman/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type Firebase struct {
}

// payload: {"email":"deans@dean.com","password":"deanb87", returnSecureToken:true}
func (f *Firebase) Signin(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {

	var newU firebase.Signin
	if err := web.Unmarshal(r.Body, &newU); err != nil {
		return errors.Cause(err)
	}

	url := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=AIzaSyBItfzjx74wXWCet-ARldNNpKIZVR1PQ5I%0A"

	req, err := http.NewRequest("POST", url, r.Body)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var ag interface{}
	json.NewDecoder(res.Body).Decode(&ag)
	web.Respond(log, w, ag, res.StatusCode)

	return nil
}

// Create a new user in the firebase account
// payload: {"email":"deans@dean.com","password":"deanb87"}
func (f *Firebase) Create(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	fmt.Println("creating")
	var newU firebase.Create
	if err := web.Unmarshal(r.Body, &newU); err != nil {
		return errors.Cause(err)
	}

	b, _ := json.Marshal(newU)

	url := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/signupNewUser?key=AIzaSyBItfzjx74wXWCet-ARldNNpKIZVR1PQ5I%0A"
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	fmt.Println("OK2")

	var ag interface{}
	json.NewDecoder(res.Body).Decode(&ag)
	web.Respond(log, w, ag, res.StatusCode)

	return nil
}

//payload `{"idToken":"Users JWT"}`
func (f *Firebase) Delete(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {

	var newU firebase.Delete
	if err := web.Unmarshal(r.Body, &newU); err != nil {
		return errors.Cause(err)
	}
	url := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/deleteAccount?key=AIzaSyBItfzjx74wXWCet-ARldNNpKIZVR1PQ5I%0A"
	req, _ := http.NewRequest("POST", url, r.Body)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	var ag interface{}
	json.NewDecoder(res.Body).Decode(&ag)
	web.Respond(log, w, ag, res.StatusCode)

	return nil
}
