package main

import (
	"example.com/REST-api/db"
	"example.com/REST-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize storage before registering handlers because model functions
	// use the package-level database handle.
	db.InitDB()

	server := gin.Default()

	routes.RegisterRoutes(server)

	// Run starts the HTTP server and blocks while it accepts requests.
	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
