// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type MySQLDB struct {
	Database *sqlx.DB
}

func New(url string) (*MySQLDB, error) {
	// db, err := sqlx.Connect("mysql", "passman:wtfthispasswordNeedsLonger29@tcp(production-database.cfneifgjtyib.us-east-1.rds.amazonaws.com:3306)/passman")
	mysql, err := sqlx.Connect("mysql", url)

	if err != nil {
		return nil, errors.Wrap(err, "Error connecting to DB when application started")
	}
	mysql.SetMaxIdleConns(10)
	mysql.DB.SetMaxIdleConns(10)
	mysql.DB.SetMaxOpenConns(10)

	db := MySQLDB{
		Database: mysql,
	}

	return &db, nil
}

// Close closes a DB value being used with MongoDB.
func (db *MySQLDB) Close() {
	db.Database.Close()
}
