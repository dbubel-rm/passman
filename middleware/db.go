package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AddUserDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		localID, exists := c.Get("localId")
		if exists {
			_, err := db.NamedExec(`INSERT INTO users (local_id) VALUES (:local_id)`,
				map[string]interface{}{
					"local_id": localID,
				})
			if err != nil {
				log.Print(err.Error())
			}
			return
		}

	}
}
