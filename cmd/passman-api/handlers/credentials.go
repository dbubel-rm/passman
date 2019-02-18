package handlers

import (
	"log"
	"net/http"

	"github.com/dbubel/passman/internal/credentials"
	"github.com/dbubel/passman/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type Credentials struct {
	MasterDB *sqlx.DB
}

func (c *Credentials) Add(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	var newU credentials.Add
	if err := web.Unmarshal(r.Body, &newU); err != nil {
		return errors.Cause(err)
	}
	return credentials.AddCredentialDB(c.MasterDB, &newU, r.Context().Value("localId"))
}

func (c *Credentials) Get(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	cred, err := credentials.GetCredentialDB(c.MasterDB, params.ByName("serviceName"), r.Context().Value("localId"))
	if err != nil {
		web.RespondError(log, w, err, http.StatusInternalServerError)
	}
	web.Respond(log, w, cred, http.StatusOK)
	return nil
}
