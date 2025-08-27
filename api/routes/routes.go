package routes

import "github.com/gin-gonic/gin"

type reqID struct {
	ID int `json:"id" binding:"required"`
}

func RegisterRoutes(server *gin.Engine) {
	server.GET("/users-empty", isUsersEmpty)
	server.POST("/first-user", createFirstUser)
	server.POST("/login", login)

	authenticated := server.Group("/", authenticate)
	authenticated.GET("/check-session", checkSession)
	authenticated.GET("/saved-ldaps", getLdapsNames)
	authenticated.PUT("/password", updatePassword)

	authenticated.POST("/ldap-childs", getChilds)
	authenticated.POST("/ldap-attributes", getAttributes)

	authenticatedAdmin := authenticated.Group("/", authenticateAdmin)
	authenticatedAdmin.GET("/user", getUser)
	authenticatedAdmin.GET("/users", getUsers)
	authenticatedAdmin.POST("/user", createUser)
	authenticatedAdmin.PUT("/user", updateUser)
	authenticatedAdmin.DELETE("/user", deleteUser)

	authenticatedAdmin.GET("/ldap", getLdap)
	authenticatedAdmin.GET("/ldaps", getLdaps)
	authenticatedAdmin.POST("/ldap", createLdap)
	authenticatedAdmin.PUT("/ldap", updateLdap)
	authenticatedAdmin.DELETE("/ldap", deleteLdap)

	authenticatedAdmin.PUT("/attribute-value", updateAttributeValue)
	authenticatedAdmin.DELETE("/attribute-value", deleteAttributeValue)
}
