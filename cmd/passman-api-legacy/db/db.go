package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// for testing docker run --name some-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:latest

func GetDB() (*sqlx.DB, error) {
	log.Println("starting DB")
	endPoint := os.Getenv("MYSQL_ENDPOINT")
	if endPoint == "" {
		endPoint = "root:my-secret-pw@tcp(127.0.0.1:3306)/passman"
	}
	endPoint = endPoint
	db, err := sqlx.Connect("mysql", endPoint)
	// db, err := sqlx.Connect("mysql", "passman:wtfthispasswordNeedsLonger29@tcp(production-database.cfneifgjtyib.us-east-1.rds.amazonaws.com:3306)/passman")
	if err != nil {
		log.Println(err)
	}
	return db, err
}
