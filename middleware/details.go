package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maengsanha/vida-backend/usecase/search"
)

// Details returns a handler for /api/details.
func Details() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tag, exists := c.GetQuery("tag"); !exists || len(tag) < 1 {
			c.JSON(http.StatusBadRequest, nil)
		} else {
			c.JSON(http.StatusOK, search.Details(tag))
		}
	}
}
