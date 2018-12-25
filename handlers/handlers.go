package handlers

import (
	"github.com/gin-gonic/gin"
)

// GetEngine returns gin engine with routes
func GetEngine(authHandler func(*gin.Context)) *gin.Engine {
	router := gin.Default()
	router.Use(authHandler)
	router.POST("/cred", addCredential)
	router.DELETE("/cred/:credId", deleteCredential)
	router.PUT("/cred/:credId", updateCredential)
	router.GET("/creds", getCredential)
	router.GET("/cred/:credId", getCredential)
	return router
}

func addCredential(c *gin.Context) {

}
func deleteCredential(c *gin.Context) {

}
func updateCredential(c *gin.Context) {

}
func getCredential(c *gin.Context) {

}
