package handlers

import (
	"net/http"
	"pg-manager-backend/models"
	"pg-manager-backend/services"

	"github.com/gin-gonic/gin"
)

// PublicRaiseComplaint is the endpoint for the QR code form
func PublicRaiseComplaint(c *gin.Context) {
	var input models.Complaint
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err := services.RegisterComplaint(input)
	if err != nil {
		if err.Error() == "tenant_not_verified" {
			// This triggers the "Verification Failed" message in your Vue app
			c.JSON(403, gin.H{"error": "Phone number not registered with this property"})
			return
		}
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"message": "Complaint submitted successfully"})
}

// GetComplaints fetches all complaints for a specific property
func GetComplaints(c *gin.Context) {
	propertyID := c.Query("property_id")
	if propertyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Property ID is required"})
		return
	}

	complaints, err := services.GetAllComplaints(propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch complaints"})
		return
	}

	c.JSON(http.StatusOK, complaints)
}

// MarkComplaintResolved updates the status so the owner can clear their dashboard
func MarkComplaintResolved(c *gin.Context) {
	id := c.Param("id")
	if err := services.ResolveComplaint(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Complaint marked as resolved"})
}
