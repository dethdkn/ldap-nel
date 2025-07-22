package routes

import (
	"github.com/dethdkn/ldap-nel/backend/utils"
	"github.com/gin-gonic/gin"
)

func authenticate(c *gin.Context) {
	token, err := c.Cookie("session")
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	username, admin, err := utils.JWTValidate(token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
		return
	}

	c.Set("username", username)
	c.Set("admin", admin)
	c.Next()
}
