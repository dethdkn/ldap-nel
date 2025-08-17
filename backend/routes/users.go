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

	c.JSON(200, gin.H{"message": "Login successful", "token": token, "isAdmin": user.Admin})
}

func createUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, gin.H{"message": "Could not bind JSON"})
		return
	}

	user.Username = utils.SanitizeUsername(user.Username)

	exists, err := models.AlreadyExists(user.Username)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to check if user exists"})
		return
	}

	if exists {
		c.JSON(400, gin.H{"message": "User already exists"})
		return
	}

	if !utils.ValidatePassword(user.Password) {
		c.JSON(400, gin.H{"message": "Password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one digit, and one special character"})
		return
	}

	if err := user.Save(); err != nil {
		c.JSON(500, gin.H{"message": "Failed to save user"})
		return
	}

	c.JSON(200, gin.H{"message": "User created successfully"})
}

func updateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, gin.H{"message": "Could not bind JSON"})
		return
	}

	originalUser, err := models.GetUserByID(user.ID)
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	user.Username = utils.SanitizeUsername(user.Username)

	if originalUser.Username != user.Username {
		exists, err := models.AlreadyExists(user.Username)
		if err != nil {
			c.JSON(500, gin.H{"message": "Failed to check if user exists"})
			return
		}

		if exists {
			c.JSON(400, gin.H{"message": "User already exists"})
			return
		}
	}

	if user.Password != "" {
		if !utils.ValidatePassword(user.Password) {
			c.JSON(400, gin.H{"message": "Password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one digit, and one special character"})
			return
		}
	}

	if err := user.Update(); err != nil {
		c.JSON(500, gin.H{"message": "Failed to update user"})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})
}

func updatePassword(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, gin.H{"message": "Could not bind JSON"})
		return
	}

	loggedInUsername, exists := c.Get("username")
	if !exists || loggedInUsername != user.Username {
		c.JSON(403, gin.H{"message": "You can only update your own password"})
		return
	}

	if err := user.UpdatePassword(); err != nil {
		c.JSON(500, gin.H{"message": "Failed to update password"})
		return
	}

	c.JSON(200, gin.H{"message": "Password updated successfully"})
}

func deleteUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, gin.H{"message": "Could not bind JSON"})
		return
	}

	if err := user.Delete(); err != nil {
		c.JSON(500, gin.H{"message": "Failed to delete user"})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

func getUser(c *gin.Context) {
	var req reqID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Could not bind JSON"})
		return
	}

	user, err := models.GetUserByID(int64(req.ID))
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	c.JSON(200, user)
}

func getUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to retrieve users"})
		return
	}

	c.JSON(200, users)
}

func isUsersEmpty(c *gin.Context) {
	empty, err := models.IsUsersEmpty()
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to check if users are empty"})
		return
	}

	c.JSON(200, gin.H{"empty": empty})
}

func createFirstUser(c *gin.Context) {
	empty, err := models.IsUsersEmpty()
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to check if users are empty"})
		return
	}

	if !empty {
		c.JSON(400, gin.H{"message": "Users already exist"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, gin.H{"message": "Could not bind JSON"})
		return
	}

	user.Admin = true // First user is always admin

	if err := user.Save(); err != nil {
		c.JSON(500, gin.H{"message": "Failed to save user"})
		return
	}

	c.JSON(200, gin.H{"message": "User created successfully"})
}
