package main

import (
	"log"
	"runtime"

	"github.com/haxana/vida-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var engine = gin.Default()

	engine.GET("/api/usercount", middleware.HandleUserCount)

	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("gin.Run: %v", err)
	}
}
