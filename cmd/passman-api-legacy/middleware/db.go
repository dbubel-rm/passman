package middleware

import (
	"fmt"
	"net/http"

	"github.com/dbubel/passman/cmd/passman-api-legacy/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// func AddUserDB(db *sqlx.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		localID, exists := c.Get(models.FbJSON)
// 		if exists {
// 			_, err := db.NamedExec(`INSERT INTO users (local_id, email) VALUES (:local_id, :email)`,
// 				map[string]interface{}{
// 					"local_id": localID.(models.FirebaseAuthResp).LocalID,
// 					"email":    localID.(models.FirebaseAuthResp).Email,
// 				})
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}
// 			c.JSON(http.StatusOK, localID.(models.FirebaseAuthResp))
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "fbJson not present"})
// 		return
// 	}
// }

func DeleteUserDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		localID, localIDExist := c.Get("localID")

		if localIDExist {
			_, err := db.NamedExec(`DELETE FROM credentials where local_id = :local_id`,
				map[string]interface{}{
					"local_id": localID,
				})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "User deleted"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Context parameters not present"})
		return
	}
}

func AddCredentialDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		credentials, a := c.Get("credentials")
		u, _ := c.Get("localID") // should rename to local id

		if a {

			_, err := db.Exec(`INSERT INTO credentials 
			(local_id, service_name, username, password) 
			values ("%s", :service_name, :username,:password)`, credentials)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "Credential added"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Context parameters not present"})
		return
	}
}

func GetCredentialDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		u, _ := c.Get("localID") // should rename to local id
		serviceName := c.Param("serviceName")

		fmt.Println(u, serviceName)

		if serviceName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": ""})
			return
		}

		jason := models.Credentials{}
		q := `select username, password, service_name 
		from credentials 
		where service_name = ?
		and local_id = ?`
		err := db.Get(&jason, q, serviceName, u)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, jason)
		return

	}
}

func DeleteCredentialDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		u, _ := c.Get("localID") // should rename to local id
		serviceName := c.Param("serviceName")

		fmt.Println(u, serviceName)

		if serviceName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": ""})
			return
		}

		q := `delete from credentials 
		where service_name = :service_name
		and local_id = :local_id`

		_, err := db.NamedExec(q, map[string]interface{}{
			"service_name": serviceName,
			"local_id":     u,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Credential deleted"})
		return

	}
}

func GetCredentialsDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		u, _ := c.Get("localID") // should rename to local id
		// serviceName := c.Param("serviceName")

		// fmt.Println(u, serviceName)

		// if serviceName == "" {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": ""})
		// 	return
		// }

		jason := []models.Credentials{}
		q := `select username, password, service_name 
		from credentials 
		where local_id = ?`
		err := db.Select(&jason, q, u)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, jason)
		return

	}
}
