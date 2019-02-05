package handlers

import (
	"log"
	"net/http"

	"github.com/dbubel/passman/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

// Check provides support for orchestration health checks.
type Check struct {
	MasterDB *sqlx.DB
}

// Health validates the service is healthy and ready to accept requests.
func (c *Check) Health(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	err := c.MasterDB.Ping()
	if err != nil {
		return err
	}
	web.Respond(log, w, status, http.StatusOK)
	return nil
}
