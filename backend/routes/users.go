package routes

import (
	"github.com/dethdkn/ldap-nel/backend/models"
	"github.com/gin-gonic/gin"
)

func createUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, gin.H{"message": "Could not bind JSON"})
		return
	}

	if err := user.Save(); err != nil {
		c.JSON(500, gin.H{"message": "Failed to save user"})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully", "id": user.ID})
}
