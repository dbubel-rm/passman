package mid

import (
	"log"
	"net/http"

	"runtime/debug"

	"github.com/dbubel/passman/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// ErrorHandler for catching and responding to errors.
func ErrorHandler(before web.Handler) web.Handler {

	// Create the handler that will be attached in the middleware chain.
	h := func(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {

		// In the event of a panic, we want to capture it here so we can send an
		// error down the stack.
		defer func() {
			if r := recover(); r != nil {

				// Indicate this request had an error.

				// Log the panic.
				log.Printf("%s : ERROR : Panic Caught : %s\n", r)

				// Respond with the error.
				web.RespondError(log, w, errors.New("unhandled"), http.StatusInternalServerError)

				// Print out the stack.
				log.Printf("%s : ERROR : Stacktrace\n%s\n", debug.Stack())
			}
		}()

		if err := before(log, w, r, params); err != nil {

			// Indicate this request had an error.

			// What is the root error.
			err = errors.Cause(err)

			if err != web.ErrNotFound {

				// Log the error.
				log.Printf("%s : ERROR : %v\n", err)
			}

			// Respond with the error.
			web.Error(log, w, err)

			// The error has been handled so we can stop propagating it.
			return nil
		}

		return nil
	}

	return h
}
