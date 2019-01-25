package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func GetDB() (*sqlx.DB, error) {
	log.Println("starting DB")
	db, err := sqlx.Connect("mysql", "passman:wtfthispasswordNeedsLonger29@tcp(production-database.cfneifgjtyib.us-east-1.rds.amazonaws.com:3306)/passman")
	if err != nil {
		log.Println(err)
	}
	return db, err
}
