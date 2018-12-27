package middleware

import (
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
		localID, localIDExist := c.Get(models.LocalID)
		fbResp, fbRespExist := c.Get(models.FirebaseResp)

		if localIDExist && fbRespExist {
			_, err := db.NamedExec(`DELETE FROM users where local_id = :local_id`,
				map[string]interface{}{
					"local_id": localID,
				})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, fbResp)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Context parameters not present"})
		return
	}
}
