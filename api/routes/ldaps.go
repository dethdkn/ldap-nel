package routes

import (
	"github.com/dethdkn/ldap-nel/api/models"
	"github.com/gin-gonic/gin"
)

type reqLdap struct {
	ID int64  `json:"id" binding:"required"`
	DN string `json:"dn"`
}

type reqLdapAttributeValue struct {
	ID        int64  `json:"id" binding:"required"`
	DN        string `json:"dn" binding:"required"`
	Attribute string `json:"attribute" binding:"required"`
	Value     string `json:"value" binding:"required"`
}

func createLdap(c *gin.Context) {
	var ldap models.Ldap
	if err := c.ShouldBindJSON(&ldap); err != nil {
		c.JSON(500, gin.H{"message": "Could not bind JSON"})
		return
	}

	if err := ldap.Save(); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
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
		c.JSON(500, gin.H{"message": err.Error()})
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
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "LDAP configuration deleted successfully"})
}

func getLdap(c *gin.Context) {
	var req reqID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Could not bind JSON"})
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

func getChilds(c *gin.Context) {
	var req reqLdap
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Could not bind JSON"})
		return
	}

	dn, children, err := models.GetLdapChilds(req.ID, req.DN)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"dn": dn, "childs": children})
}

func getAttributes(c *gin.Context) {
	var req reqLdap
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Could not bind JSON"})
		return
	}

	attributes, err := models.GetLdapAttributes(req.ID, req.DN)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"attributes": attributes})
}

func getPossibleAttributes(c *gin.Context) {
	var req reqLdap
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Could not bind JSON"})
		return
	}

	attributes, err := models.GetLdapPossibleAttributes(req.ID, req.DN)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"possibleAttributes": attributes})
}

func addAttributeValue(c *gin.Context) {
	var req reqLdapAttributeValue

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Could not bind JSON"})
		return
	}

	if err := models.AddLdapAttributeValue(req.ID, req.DN, req.Attribute, req.Value); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Attribute value added successfully"})
}

func updateAttributeValue(c *gin.Context) {
	var req struct {
		ID        int64  `json:"id" binding:"required"`
		DN        string `json:"dn" binding:"required"`
		Attribute string `json:"attribute" binding:"required"`
		Value     string `json:"value" binding:"required"`
		NewValue  string `json:"newValue" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Could not bind JSON"})
		return
	}

	if err := models.UpdateLdapAttributeValue(req.ID, req.DN, req.Attribute, req.Value, req.NewValue); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Attribute value updated successfully"})
}

func deleteAttributeValue(c *gin.Context) {
	var req reqLdapAttributeValue

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Could not bind JSON"})
		return
	}

	if err := models.DeleteLdapAttributeValue(req.ID, req.DN, req.Attribute, req.Value); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Attribute value deleted successfully"})
}
