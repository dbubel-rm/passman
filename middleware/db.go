package middleware

import (
	"fmt"
	"net/http"

	"github.com/dbubel/passman/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AddUserDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		localID, exists := c.Get(models.FbJSON)
		if exists {
			_, err := db.NamedExec(`INSERT INTO users (local_id, email) VALUES (:local_id, :email)`,
				map[string]interface{}{
					"local_id": localID.(models.FirebaseAuthResp).LocalID,
					"email":    localID.(models.FirebaseAuthResp).Email,
				})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, localID.(models.FirebaseAuthResp))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fbJson not present"})
		return
	}
}

func DeleteUserDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		localID, localIDExist := c.Get("userID")

		if localIDExist {
			_, err := db.NamedExec(`DELETE FROM users where local_id = :local_id`,
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

// UserID      int       `json:"userId"`
// ServiceName string    `json:"serviceName"`
// Username    string    `json:"username"`
// Password    string    `json:"password"`
// UpdatedAt   time.Time `json:"updatedAt"`
// CreatedAt   time.Time `json:"createdAt"`
// DeletedAt   time.Time `json:"deletedAt"`

func AddCredentialDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// localID, localIDExist := c.Get("userID")
		credentials, a := c.Get("credentials")
		u, _ := c.Get("userID") // should rename to local id

		if a {
			q := fmt.Sprintf(`INSERT INTO credentials (user_id, service_name, username, password) 
			values ( (select user_id from users where local_id = "%s"), :service_name, :username,:password)`, u)
			_, err := db.NamedExec(q, credentials)

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

		u, _ := c.Get("userID") // should rename to local id
		serviceName := c.Param("serviceName")

		fmt.Println(u, serviceName)

		if serviceName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": ""})
			return
		}

		jason := models.Credentials{}
		err := db.Get(&jason, `select c.username, c.password, c.service_name 
			from passman.users u 
			join credentials c on u.user_id = c.user_id 
			where c.service_name = ?
			and u.local_id = ?`, serviceName, u)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, jason)
		return

	}
}
