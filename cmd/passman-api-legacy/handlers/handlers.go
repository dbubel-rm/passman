package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dbubel/passman/cmd/passman-api-legacy/middleware"
	"github.com/dbubel/passman/cmd/passman-api-legacy/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func init() {
	gin.SetMode(gin.DebugMode)
}

// GetEngine returns gin engine with routes
func GetEngine(authHandler func(*gin.Context), db *sqlx.DB) *gin.Engine {
	var router *gin.Engine

	router = gin.Default()
	privateAPI := router.Group("/")
	privateAPI.Use(authHandler)
	{
		privateAPI.POST("/credential", addCredentials, middleware.AddCredentialDB(db))
		privateAPI.GET("/credential/:serviceName", getCredential, middleware.GetCredentialDB(db))
		privateAPI.GET("/credentials", getCredential, middleware.GetCredentialsDB(db))
		privateAPI.DELETE("/user", deleteUser, middleware.DeleteUserDB(db))
		privateAPI.DELETE("/credential/:serviceName", deleteCredential, middleware.DeleteCredentialDB(db))
	}

	publicAPI := router.Group("/public")
	publicAPI.POST("/user", addUser)
	publicAPI.POST("/authUser", authUser)

	router.GET("/health", health)

	return router
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "v0.0.2"})
	return
}

func authUser(c *gin.Context) {
	// Make firebase call
	url := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=AIzaSyBItfzjx74wXWCet-ARldNNpKIZVR1PQ5I%0A"
	req, _ := http.NewRequest("POST", url, c.Request.Body)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.StatusCode != http.StatusOK {
		var fbResp interface{}
		json.NewDecoder(res.Body).Decode(&fbResp)
		defer res.Body.Close()
		c.JSON(res.StatusCode, fbResp)
		c.Abort()
		return
	}

	var fbJson models.FirebaseAuthResp
	err = json.NewDecoder(res.Body).Decode(&fbJson)
	if err != nil {
		log.Println(err.Error())
	}
	c.JSON(http.StatusOK, fbJson)
	return
}

func addUser(c *gin.Context) {
	// Make firebase call
	url := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/signupNewUser?key=AIzaSyBItfzjx74wXWCet-ARldNNpKIZVR1PQ5I%0A"
	req, _ := http.NewRequest("POST", url, c.Request.Body)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.StatusCode != http.StatusOK {
		var fbResp interface{}
		json.NewDecoder(res.Body).Decode(&fbResp)
		defer res.Body.Close()
		c.JSON(res.StatusCode, fbResp)
		c.Abort()
		return
	}

	var fbJson models.FirebaseAuthResp
	err = json.NewDecoder(res.Body).Decode(&fbJson)
	if err != nil {
		log.Println(err.Error())
	}
	c.JSON(http.StatusOK, fbJson)
	return
}

func deleteUser(c *gin.Context) {
	url := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/deleteAccount?key=AIzaSyBItfzjx74wXWCet-ARldNNpKIZVR1PQ5I%0A"

	idToken, ok := c.Get("jwt")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get firebaseJSON"})
		c.Abort()
		return
	}

	m := strings.NewReader(fmt.Sprintf(`{"idToken":"%s"}`, idToken))

	req, _ := http.NewRequest("POST", url, m)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	var fbResp interface{}
	json.NewDecoder(res.Body).Decode(&fbResp)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.JSON(res.StatusCode, fbResp)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, fbResp)
	return

}

func addCredentials(c *gin.Context) {
	var cred models.Credentials
	err := c.BindJSON(&cred)
	if err != nil {
		log.Println(err.Error())
	}
	c.Set("credentials", cred)
	c.Next()
}

func deleteCredential(c *gin.Context) {
	c.Next()
}

func updateCredential(c *gin.Context) {

}

func getCredential(c *gin.Context) {
	c.Next()
}
