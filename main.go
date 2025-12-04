package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"tallyRead.com/db"
	"tallyRead.com/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Env file not loaded")
	}
	db.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8080")
}
