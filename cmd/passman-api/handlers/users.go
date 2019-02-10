package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/dbubel/passman/internal/platform/web"
	"github.com/dbubel/passman/internal/user"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type User struct {
	MasterDB *sqlx.DB
}

// Create inserts a new user into the system.
func (u *User) Create(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	fmt.Println("CREATE")

	var newU user.NewUser
	if err := web.Unmarshal(r.Body, &newU); err != nil {
		return errors.Wrap(err, "Error unmarshaling")
	}

	usr, err := user.Create(u.MasterDB, &newU)
	if err = translate(err); err != nil {
		return errors.Wrapf(err, "User: %+v", &usr)
	}

	web.Respond(log, w, usr, http.StatusCreated)
	return nil
}

func translate(err error) error {
	switch errors.Cause(err) {
	case user.ErrNotFound:
		return web.ErrNotFound
	case user.ErrInvalidID:
		return web.ErrInvalidID
	case user.ErrAuthenticationFailure:
		return web.ErrUnauthorized
	case user.ErrForbidden:
		return web.ErrForbidden
	}
	return err
}
