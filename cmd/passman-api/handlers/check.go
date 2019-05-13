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

var (
	Build      string
	GitHash    string
	BuildDate  string
	InstanceID string
)

// Health validates the service is healthy and ready to accept requests.
func (c *Check) Health(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {

	status := struct {
		DBStatus   string `json:"dbStatus"`
		Version    string `json:"version"`
		GitHash    string `json:"gitHash"`
		BuildDate  string `json:"buildDate"`
		InstanceID string `json:"instanceId"`
	}{
		DBStatus:   "ok",
		Version:    Build,
		GitHash:    GitHash,
		BuildDate:  BuildDate,
		InstanceID: InstanceID,
	}

	err := c.MasterDB.Ping()
	if err != nil {
		status.DBStatus = err.Error()
	}

	web.Respond(log, w, status, http.StatusOK)
	return nil
}
