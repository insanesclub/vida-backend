package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haxana/vida-backend/usecase/detectfakeuser"
)

func HandleUserCount(c *gin.Context) {
	c.JSON(http.StatusOK, detectfakeuser.Usernames(c.Query("tag")))
}
