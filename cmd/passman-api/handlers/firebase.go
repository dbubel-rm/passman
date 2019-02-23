package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dbubel/passman/internal/firebase"
	"github.com/dbubel/passman/internal/platform/web"
	"github.com/julienschmidt/httprouter"
)

type Firebase struct {
	SigninURL string
	CreateURL string
	DeleteURL string
	VerifyURL string
}

// Signin - Get a valid JWT for the user
func (f *Firebase) Signin(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	var signinReq firebase.Signin

	if err := web.Unmarshal(r.Body, &signinReq); err != nil {
		return err
	}

	defer r.Body.Close()
	signinJSON, err := json.Marshal(signinReq)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", f.SigninURL, bytes.NewReader(signinJSON))

	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	var firebaseResp interface{}
	err = json.NewDecoder(res.Body).Decode(&firebaseResp)

	if err != nil {
		return err
	}

	web.Respond(log, w, firebaseResp, res.StatusCode)
	return nil
}

// Create - adds new user in the firebase account
func (f *Firebase) Create(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	var createReq firebase.Create
	if err := web.Unmarshal(r.Body, &createReq); err != nil {
		return err
	}

	defer r.Body.Close()
	createJSON, err := json.Marshal(createReq)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", f.CreateURL, bytes.NewReader(createJSON))

	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	var firebaseResp interface{}
	err = json.NewDecoder(res.Body).Decode(&firebaseResp)

	if err != nil {
		return err
	}

	web.Respond(log, w, firebaseResp, res.StatusCode)
	return nil
}

// Delete - removes a users firebase account
func (f *Firebase) Delete(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	var deleteReq firebase.Delete

	if err := web.Unmarshal(r.Body, &deleteReq); err != nil {
		return err
	}

	defer r.Body.Close()
	deleteJSON, err := json.Marshal(deleteReq)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", f.DeleteURL, bytes.NewReader(deleteJSON))

	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	var firebaseResp interface{}
	err = json.NewDecoder(res.Body).Decode(&firebaseResp)

	if err != nil {
		return err
	}

	web.Respond(log, w, firebaseResp, res.StatusCode)
	return nil
}

// Verify - send an email verification
func (f *Firebase) Verify(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	var verifyReq firebase.Verify

	if err := web.Unmarshal(r.Body, &verifyReq); err != nil {
		return err
	}

	defer r.Body.Close()
	verifyJSON, err := json.Marshal(verifyReq)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", f.VerifyURL, bytes.NewReader(verifyJSON))

	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	var firebaseResp interface{}
	err = json.NewDecoder(res.Body).Decode(&firebaseResp)

	if err != nil {
		return err
	}

	web.Respond(log, w, firebaseResp, res.StatusCode)
	return nil
}
