package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/dbubel/passman/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetEngine returns gin engine with routes
func GetEngine(authHandler func(*gin.Context), db *sqlx.DB) *gin.Engine {
	var router *gin.Engine
	if os.Getenv("TEST") != "" {
		router = gin.New()
	} else {
		router = gin.Default()
	}

	credAPI := router.Group("/cred")
	credAPI.Use(authHandler)
	{
		credAPI.POST("/add", addCredential(db))
		credAPI.PUT("/update/:credId", updateCredential)
		credAPI.GET("/get", getCredential)
		credAPI.GET("/get/:credId", getCredential)
		credAPI.DELETE("/delete/:credId", deleteCredential)

		credAPI.DELETE("/user", deleteUser(db))
	}

	router.POST("/user", addUser(db))
	return router
}

func addUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/signupNewUser?key=AIzaSyBItfzjx74wXWCet-ARldNNpKIZVR1PQ5I%0A"
		req, _ := http.NewRequest("POST", url, c.Request.Body)
		res, err := http.DefaultClient.Do(req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var f models.FirebaseAuthResp
		json.NewDecoder(res.Body).Decode(&f)
		c.JSON(res.StatusCode, f)
		res.Body.Close()
	}
}

func deleteUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/deleteAccount?key=AIzaSyBItfzjx74wXWCet-ARldNNpKIZVR1PQ5I%0A"
		idToken, ok := c.Get("firebaseJSON")
		if !ok {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		payload := strings.NewReader(string(idToken.([]byte)))
		req, _ := http.NewRequest("POST", url, payload)
		res, err := http.DefaultClient.Do(req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var f models.FirebaseAuthResp
		json.NewDecoder(res.Body).Decode(&f)
		c.JSON(res.StatusCode, f)
		res.Body.Close()
	}
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
