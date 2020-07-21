// Package main runs API server of VIDA.
package main

import "github.com/gin-gonic/gin"

func main() {
	engine := gin.Default()
	engine.Run(":8080")
}
