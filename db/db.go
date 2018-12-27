package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func GetDB() *sqlx.DB {
	db, err := sqlx.Connect("mysql", "root:@/passman")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("DB ok")
	return db
}
