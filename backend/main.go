package main

import (
	"github.com/dethdkn/ldap-nel/backend/db"
	"github.com/dethdkn/ldap-nel/backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	server := gin.Default()

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "LDAP-NEL API",
		})
	})

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
