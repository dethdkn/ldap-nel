package main

import (
	"github.com/dethdkn/ldap-nel/backend/db"
	"github.com/dethdkn/ldap-nel/backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":3001")
}
