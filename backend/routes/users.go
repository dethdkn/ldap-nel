package routes

import (
	"github.com/dethdkn/ldap-nel/backend/models"
	"github.com/dethdkn/ldap-nel/backend/utils"
	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, gin.H{"message": "Could not bind JSON"})
		return
	}

	if err := user.Validate(); err != nil {
		c.JSON(401, gin.H{"message": "Invalid username or password"})
		return
	}

	token, err := utils.JWTGenerate(user.Username, user.Admin)
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not generate token"})
		return
	}

	c.JSON(200, gin.H{"message": "Login successful", "token": token})
}

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

	c.JSON(201, gin.H{"message": "User created successfully"})
}
