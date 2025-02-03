package main

import (
	"pet-search-backend-go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8080")
}