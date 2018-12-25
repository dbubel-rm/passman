package main

import (
	"github.com/dbubel/passman/handlers"
	"github.com/dbubel/passman/middleware"
)

func main() {
	engine := handlers.GetEngine(middleware.AuthUser)
	router.POST("/someGet", middleware.AuthUser, hello)
	engine.Run()
}
