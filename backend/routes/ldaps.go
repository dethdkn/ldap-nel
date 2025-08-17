package routes

import (
	"github.com/dethdkn/ldap-nel/backend/models"
	"github.com/gin-gonic/gin"
)

func createLdap(c *gin.Context) {
	var ldap models.Ldap
	if err := c.ShouldBindJSON(&ldap); err != nil {
		c.JSON(500, gin.H{"message": "Could not bind JSON"})
		return
	}

	if err := ldap.Save(); err != nil {
		c.JSON(500, gin.H{"message": "Failed to save LDAP configuration"})
		return
	}

	c.JSON(200, gin.H{"message": "LDAP configuration created successfully"})
}

func updateLdap(c *gin.Context) {
	var ldap models.Ldap
	if err := c.ShouldBindJSON(&ldap); err != nil {
		c.JSON(500, gin.H{"message": "Could not bind JSON"})
		return
	}

	if err := ldap.Update(); err != nil {
		c.JSON(500, gin.H{"message": "Failed to update LDAP configuration"})
		return
	}

	c.JSON(200, gin.H{"message": "LDAP configuration updated successfully"})
}

func deleteLdap(c *gin.Context) {
	var ldap models.Ldap
	if err := c.ShouldBindJSON(&ldap); err != nil {
		c.JSON(500, gin.H{"message": "Could not bind JSON"})
		return
	}

	if err := ldap.Delete(); err != nil {
		c.JSON(500, gin.H{"message": "Failed to delete LDAP configuration"})
		return
	}

	c.JSON(200, gin.H{"message": "LDAP configuration deleted successfully"})
}

func getLdap(c *gin.Context) {
	var req reqID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	ldap, err := models.GetLdapByID(int64(req.ID), false)
	if err != nil {
		c.JSON(404, gin.H{"message": "LDAP configuration not found"})
		return
	}

	c.JSON(200, ldap)
}

func getLdaps(c *gin.Context) {
	ldaps, err := models.GetAllLdaps()
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to retrieve LDAP configurations"})
		return
	}

	c.JSON(200, ldaps)
}

func getLdapsNames(c *gin.Context) {
	ldaps, err := models.GetAllLdapsNames()
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to retrieve LDAP configurations"})
		return
	}
	c.JSON(200, ldaps)
}
