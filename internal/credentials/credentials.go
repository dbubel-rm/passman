package credentials

import (
	"github.com/jmoiron/sqlx"
)

func AddUserDB(dbConn *sqlx.DB, cred *Add, localID interface{}) error {
	_, err := dbConn.Exec(`INSERT INTO credentials
	(local_id, service_name, username, password)
	values (?, ?, ?, ?)`, localID, cred.ServiceName, cred.Username, cred.Password)
	return err
}
