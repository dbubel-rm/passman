package credentials

import (
	"github.com/jmoiron/sqlx"
)

func AddCredentialDB(dbConn *sqlx.DB, cred *Add, localID interface{}) error {
	_, err := dbConn.Exec(`INSERT INTO credentials
	(local_id, service_name, username, password)
	values (?, ?, ?, ?)`, localID, cred.ServiceName, cred.Username, cred.Password)
	return err
}

func GetCredentialDB(dbConn *sqlx.DB, serviceName string, localID interface{}) (*Credential, error) {
	var jason Credential

	err := dbConn.Get(&jason, `select username, password, service_name 
	from credentials 
	where service_name = ?
	and local_id = ?`, serviceName, localID)

	return &jason, err
}

func DeleteCredentialDB(dbConn *sqlx.DB, serviceName string, localID interface{}) error {

	_, err := dbConn.Exec(`delete from credentials 
	where service_name = ?
	and local_id = ?`, serviceName, localID)
	return err
}

func UpdateCredentialDB(dbConn *sqlx.DB, serviceName, password string, localID interface{}) error {

	_, err := dbConn.Exec(`update credentials 
	set password = ?
	where service_name = ?
	and local_id = ?`, password, serviceName, localID)
	return err
}

func GetServicesDB(dbConn *sqlx.DB, localID interface{}) (*[]Service, error) {
	var serviceList []Service
	err := dbConn.Select(&serviceList, `select service_name from credentials where local_id = ?`, localID)
	return &serviceList, err
}
