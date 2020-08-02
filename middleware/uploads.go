package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haxana/vida-backend/usecase/crawler"
)

// Uploads returns a handler for /api/uploads.
func Uploads() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tag, exists := c.GetQuery("tag"); !exists || len(tag) < 1 {
			c.JSON(http.StatusBadRequest, nil)
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"tag":     tag,
				"date":    time.Now().Format("2006-01-02"),
				"uploads": crawler.Uploads(tag),
			})
		}
	}
}

// Test is just a test middleware.
func Test() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tag, exists := c.GetQuery("tag"); !exists || len(tag) < 1 {
			c.JSON(http.StatusBadRequest, nil)
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"tag":  tag,
				"date": time.Now().Format("2006-01-02"),
				"uploads": map[string]interface{}{
					"user_A": 3,
					"user_B": 2,
				},
			})
		}
	}
}
