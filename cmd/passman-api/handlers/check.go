package handlers

import (
	"log"
	"net/http"

	"github.com/dbubel/passman/internal/platform/web"
	"github.com/julienschmidt/httprouter"
)

// Check provides support for orchestration health checks.
type Check struct {
}

// Health validates the service is healthy and ready to accept requests.
func (c *Check) Health(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	// ctx, span := trace.StartSpan(ctx, "handlers.Check.Health")
	// defer span.End()

	// dbConn := c.MasterDB.Copy()
	// defer dbConn.Close()

	// if err := dbConn.StatusCheck(ctx); err != nil {
	// 	return err
	// }

	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	web.Respond(log, w, status, http.StatusOK)
	return nil
}
