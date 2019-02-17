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

	localID := r.Context().Value("localId")

	return credentials.AddUserDB(c.MasterDB, &newU, localID)

}

// credentials, a := c.Get("credentials")
// u, _ := c.Get("localID") // should rename to local id

// if a {
// 	q := fmt.Sprintf(`INSERT INTO credentials
// 	(local_id, service_name, username, password)
// 	values ("%s", :service_name, :username,:password)`, u)
// 	_, err := db.NamedExec(q, credentials)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"status": "Credential added"})
// 	return
// }
// c.JSON(http.StatusInternalServerError, gin.H{"error": "Context parameters not present"})
// return
// }
