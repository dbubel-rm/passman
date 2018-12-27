package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func GetDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "root:@/passman")
	if err != nil {
		log.Fatalln(err)
	}
	return db, err
}
