package credentials

import (
	"github.com/jmoiron/sqlx"
)

func AddCredentialDB(dbConn *sqlx.DB, cred *Add, localID interface{}) error {
	_, err := dbConn.Exec(`INSERT INTO credentials (local_id, service_name, username, password) values (?, ?, ?, ?)`,
		localID, cred.ServiceName, cred.Username, cred.Password)
	return err
}

func GetCredentialDB(dbConn *sqlx.DB, serviceName string, localID interface{}) (*Credential, error) {
	var jason Credential
	err := dbConn.Get(&jason, `SELECT username, password, service_name FROM credentials WHERE service_name = ? AND local_id = ?`,
		serviceName, localID)
	return &jason, err
}

func DeleteCredentialDB(dbConn *sqlx.DB, serviceName string, localID interface{}) error {
	_, err := dbConn.Exec(`DELETE FROM credentials WHERE service_name = ? AND local_id = ?`,
		serviceName, localID)
	return err
}

func UpdateCredentialDB(dbConn *sqlx.DB, serviceName, password string, localID interface{}) error {
	_, err := dbConn.Exec(`UPDATE credentials SET password = ?, updated_at = NOW() WHERE service_name = ? AND local_id = ?`,
		password, serviceName, localID)
	return err
}

func GetServicesDB(dbConn *sqlx.DB, localID interface{}) (*[]Service, error) {
	var serviceList []Service
	err := dbConn.Select(&serviceList, `SELECT credential_id, service_name FROM credentials WHERE local_id = ?`, localID)
	return &serviceList, err
}

func DeleteAccountDB(dbConn *sqlx.DB, localID interface{}) error {
	_, err := dbConn.Exec(`DELETE FROM credentials WHERE local_id = ?`,
		localID)
	return err
}
