package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.POST("/login", login)

	authenticated := server.Group("/", authenticate)
	authenticated.PUT("/password", updatePassword)

	authenticatedAdmin := authenticated.Group("/", authenticateAdmin)
	authenticatedAdmin.GET("/user", getUser)
	authenticatedAdmin.GET("/users", getUsers)
	authenticatedAdmin.POST("/user", createUser)
	authenticatedAdmin.PUT("/user", updateUser)
	authenticatedAdmin.DELETE("/user", deleteUser)
}
