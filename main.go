package main

import (
	"example.com/chat/db"
	"example.com/chat/routes"
	"github.com/gin-gonic/gin"
)



func main() {

	db.InitDB()

	router := gin.Default()

	routes.RegisterRoutes(router)

	router.Run(":8080")
}
