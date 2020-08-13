package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Test() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tag, exists := c.GetQuery("tag"); !exists || len(tag) < 1 {
			c.JSON(http.StatusBadRequest, nil)
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"tag":  tag,
				"date": time.Now().UTC().Format(time.RFC3339[:10]),
				"uploads": map[string]interface{}{
					"user_A": 3,
					"user_B": 2,
				},
			})
		}
	}
}
