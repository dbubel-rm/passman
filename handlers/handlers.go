package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dbubel/passman/middleware"
	"github.com/dbubel/passman/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetEngine returns gin engine with routes
func GetEngine(authHandler func(*gin.Context), db *sqlx.DB) *gin.Engine {
	var router *gin.Engine

	router = gin.Default()
	credentialsAPI := router.Group("/credentials")
	credentialsAPI.Use(authHandler)
	{
		credentialsAPI.POST("/add", addCredential(db))
		// credentialsAPI.PUT("/update/:credId", updateCredential)
		// credentialsAPI.GET("/get", getCredential)
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
	log.Println("IN USER ADD")
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
	var f models.FirebaseAuthResp
	url := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/deleteAccount?key=AIzaSyBItfzjx74wXWCet-ARldNNpKIZVR1PQ5I%0A"

	idToken, ok := c.Get("firebaseJSON")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get firebaseJSON"})
		c.Abort()
		return
	}

	payload := strings.NewReader(string(idToken.([]byte)))
	req, _ := http.NewRequest("POST", url, payload)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	var fbResp interface{}
	json.NewDecoder(res.Body).Decode(&fbResp)
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		c.JSON(res.StatusCode, fbResp)
		c.Abort()
		return
	}

	err = json.NewDecoder(strings.NewReader(string(idToken.([]byte)))).Decode(&f)
	defer res.Body.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Set(models.LocalID, f.LocalID)
	c.Set(models.FirebaseResp, fbResp)
	c.Next()

}

func addCredential(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
		tx.Commit()
	}
}

func deleteCredential(c *gin.Context) {

}

func updateCredential(c *gin.Context) {

}

func getCredential(c *gin.Context) {

}
