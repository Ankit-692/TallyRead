package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"tallyRead.com/db"
	"tallyRead.com/routes"
    "github.com/gin-contrib/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Env file not loaded")
	}
	db.InitDB()
	server := gin.Default()
 	server.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:4200"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        AllowCredentials: true,
    }))
	routes.RegisterRoutes(server)
	server.Run(":8080")
}
