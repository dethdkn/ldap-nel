package main

import (
	"github.com/dethdkn/ldap-nel/api/db"
	"github.com/dethdkn/ldap-nel/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	server := gin.Default()

	server.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
