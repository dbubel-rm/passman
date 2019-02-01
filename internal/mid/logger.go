package mid

import (
	"log"
	"net/http"

	"github.com/dbubel/passman/internal/platform/web"
	"github.com/julienschmidt/httprouter"
)

// RequestLogger writes some information about the request to the logs in
// the format: TraceID : (200) GET /foo -> IP ADDR (latency)
func RequestLogger(before web.Handler) web.Handler {

	// Wrap this handler around the next one provided.
	h := func(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {

		err := before(log, w, r, params)

		log.Printf("%s : (%d) : %s %s -> %s (%s)",

			r.Method, r.URL.Path,
		)

		// For consistency return the error we received.
		return err
	}

	return h
}
