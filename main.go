package main

import (
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/maengsanha/vida-backend/middleware"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()
	api := engine.Group("/api")
	{
		api.GET("/uploads", middleware.Uploads())
	}
	engine.Run(":8080")
}
