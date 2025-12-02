package main

import (
	"github.com/gin-gonic/gin"
	"tallyRead.com/routes"
)

func main() {
	server := gin.Default()
	routes.RegisterRoutes(server)
}
