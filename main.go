package main

import (
	"io"
	"log"
	"os"

	"github.com/dbubel/passman/db"
	"github.com/dbubel/passman/handlers"
	"github.com/dbubel/passman/middleware"
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
