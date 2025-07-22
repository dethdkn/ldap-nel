package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.POST("/login", login)

	authenticated := server.Group("/", authenticate)
	authenticated.POST("/users", createUser)
}
