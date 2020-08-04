package middleware

import (
	"net/http"
	"time"

	"github.com/haxana/vida-backend/usecase/crawl"

	"github.com/gin-gonic/gin"
)

// Uploads returns a handler for /api/uploads.
func Uploads() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tag, exists := c.GetQuery("tag"); !exists || len(tag) < 1 {
			c.JSON(http.StatusBadRequest, nil)
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"tag":     tag,
				"date":    time.Now().Format(time.RFC3339[:10]),
				"uploads": crawl.Uploads(tag),
			})
		}
	}
}
