package handlers

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/jmoiron/sqlx"
)

type User struct {
	MasterDB *sqlx.DB
}

// Create inserts a new user into the system.
func (u *User) Create(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {

	// var newU user.NewUser
	// if err := web.Unmarshal(r.Body, &newU); err != nil {
	// 	return errors.Wrap(err, "Error unmarshaling")
	// }

	// usr, err := user.Create(u.MasterDB, &newU)
	// if err = translate(err); err != nil {
	// 	return errors.Wrapf(err, "User: %+v", &usr)
	// }

	// web.Respond(log, w, usr, http.StatusCreated)
	return nil
}

// func (u *User) Create(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
// 	fmt.Println("CREATE")

// 	var newU user.NewUser
// 	if err := web.Unmarshal(r.Body, &newU); err != nil {
// 		return errors.Wrap(err, "Error unmarshaling")
// 	}

// 	usr, err := user.Create(u.MasterDB, &newU)
// 	if err = translate(err); err != nil {
// 		return errors.Wrapf(err, "User: %+v", &usr)
// 	}

// 	web.Respond(log, w, usr, http.StatusCreated)
// 	return nil
// }
