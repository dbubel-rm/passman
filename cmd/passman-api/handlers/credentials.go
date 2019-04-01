package handlers

import (
	"log"

	"net/http"

	"github.com/dbubel/passman/internal/credentials"
	"github.com/dbubel/passman/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type Credentials struct {
	MasterDB *sqlx.DB
}

func (c *Credentials) add(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	var add credentials.Add

	if err := web.Unmarshal(r.Body, &add); err != nil {
		return err
	}
	defer r.Body.Close()

	err := credentials.AddCredentialDB(c.MasterDB, &add, r.Context().Value("localId"))

	if err != nil {
		web.RespondError(log, w, err, http.StatusInternalServerError)
	}

	web.Respond(log, w, nil, http.StatusOK)
	return nil
}

func (c *Credentials) get(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	credential, err := credentials.GetCredentialDB(c.MasterDB, params.ByName("serviceName"), r.Context().Value("localId"))

	if err != nil {
		web.RespondError(log, w, err, http.StatusInternalServerError)
	}

	web.Respond(log, w, credential, http.StatusOK)
	return nil
}

func (c *Credentials) delete(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	err := credentials.DeleteCredentialDB(c.MasterDB, params.ByName("serviceName"), r.Context().Value("localId"))

	if err != nil {
		web.RespondError(log, w, err, http.StatusInternalServerError)
	}

	web.Respond(log, w, nil, http.StatusOK)
	return nil
}

func (c *Credentials) update(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	var update credentials.Update

	if err := web.Unmarshal(r.Body, &update); err != nil {
		return err
	}
	defer r.Body.Close()

	err := credentials.UpdateCredentialDB(c.MasterDB, update.ServiceName, update.Password, r.Context().Value("localId"))

	if err != nil {
		web.RespondError(log, w, err, http.StatusInternalServerError)
	}

	web.Respond(log, w, nil, http.StatusOK)
	return nil
}

func (c *Credentials) services(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	services, err := credentials.GetServicesDB(c.MasterDB, r.Context().Value("localId"))

	if err != nil {
		web.RespondError(log, w, err, http.StatusInternalServerError)
	}

	web.Respond(log, w, services, http.StatusOK)
	return nil
}
