package main

import (
	"log"

	"github.com/dbubel/passman/db"
	"github.com/dbubel/passman/handlers"
	"github.com/dbubel/passman/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := db.GetDB()
	_ = db
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	engine := handlers.GetEngine(middleware.AuthUser, db)
	engine.Run()
}
