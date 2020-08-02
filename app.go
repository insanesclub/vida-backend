package main

import (
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/haxana/vida-backend/middleware"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// gin.SetMode(gin.ReleaseMode)

	var engine = gin.Default()
	var api = engine.Group("/api")
	{
		api.GET("/uploads", middleware.Uploads())
		api.GET("/test", middleware.Test())
	}
	engine.Run(":8080")
}
