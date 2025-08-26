package main

import (
	"github.com/dethdkn/ldap-nel/api/db"
	"github.com/dethdkn/ldap-nel/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
