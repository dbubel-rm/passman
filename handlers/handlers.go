package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dbubel/passman/middleware"
	"github.com/dbubel/passman/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// GetEngine returns gin engine with routes
func GetEngine(authHandler func(*gin.Context), db *sqlx.DB) *gin.Engine {
	var router *gin.Engine

	router = gin.New()
	credentialsAPI := router.Group("/credentials")
	credentialsAPI.Use(authHandler)
	{
		credentialsAPI.POST("/add", addCredentials, middleware.AddCredentialDB(db))
		// credentialsAPI.PUT("/update/:credId", updateCredential)
		credentialsAPI.GET("/get/:serviceName", getCredential, middleware.GetCredentialDB(db))
		// credentialsAPI.GET("/get/:credId", getCredential)
		// credentialsAPI.DELETE("/delete/:credId", deleteCredential)
	}

	userAPI := router.Group("/user")
	userAPI.Use(authHandler)
	{
		userAPI.DELETE("/delete", deleteUser, middleware.DeleteUserDB(db))
	}

	publicAPI := router.Group("/public")
	publicAPI.POST("/user", addUser, middleware.AddUserDB(db))

	return router
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
	c.Set(models.FbJSON, fbJson)
	c.Next()
}

func deleteUser(c *gin.Context) {
	// var f models.FirebaseAuthResp
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

	// c.Set(models.LocalID, f.LocalID)
	// c.Set(models.FirebaseResp, fbResp)
	c.Next()

}

func addCredentials(c *gin.Context) {
	var cred models.Credentials
	err := c.BindJSON(&cred)
	if err != nil {
		log.Println(err.Error())
	}
	c.Set("credentials", cred)
	c.Next()
	// userID, _ := c.Get("userID")
	// log.Println("CLAIMS", userID)

	// c.JSON(http.StatusAccepted, cred)

}

func deleteCredential(c *gin.Context) {

}

func updateCredential(c *gin.Context) {

}

func getCredential(c *gin.Context) {
	// var cred models.CredentialRequest
	// err := c.BindJSON(&cred)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// c.Set("credentials", cred)
	c.Next()
}
