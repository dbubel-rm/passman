package main

import (
	"io"
	"log"
	"os"

	"github.com/dbubel/passman/cmd/passman-api-legacy/db"
	"github.com/dbubel/passman/cmd/passman-api-legacy/handlers"
	"github.com/dbubel/passman/cmd/passman-api-legacy/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mw := io.MultiWriter(os.Stdout)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	db, err := db.GetDB()
	if err != nil {
		log.Println(err.Error())
	}
	engine := handlers.GetEngine(middleware.AuthUser, db)
	log.Println("Passman running...")
	engine.Run()
}
