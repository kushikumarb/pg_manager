package handlers

import (
	"fmt"
	"net/http"
	"pg-manager-backend/services"

	"github.com/gin-gonic/gin"
)

// AddProperty handles POST /api/properties
func AddProperty(c *gin.Context) {
	var input struct {
		Name    string `json:"name" binding:"required"`
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ownerIDFloat, _ := c.Get("user_id")
	ownerID := uint(ownerIDFloat.(float64))

	property, err := services.CreateProperty(input.Name, input.Address, ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Property added successfully",
		"property": property,
	})
}

// GetRoomsByProperty handles GET /api/properties/:id/rooms
func GetRoomsByProperty(c *gin.Context) {
	// Get ID from URL parameter (/properties/:id/rooms)
	propertyID := c.Param("id")

	rooms, err := services.GetRoomsByPropertyID(propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func GetOwnerDashboard(c *gin.Context) {
	propertyID := c.Query("property_id")

	// If viewing a specific property
	if propertyID != "" {
		stats, err := services.GetDashboardData(propertyID) // Now passing the string ID
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats)
		return
	}

	// Global Owner Overview (Optional: You can update GetDashboardData to handle uint as well)
	ownerIDFloat, _ := c.Get("user_id")
	ownerID := uint(ownerIDFloat.(float64))

	// Convert ownerID to string if you want to use the same service function
	ownerIDStr := fmt.Sprintf("%d", ownerID)
	dashboardData, err := services.GetDashboardData(ownerIDStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dashboardData)
}

func GetProperties(c *gin.Context) {
	ownerIDFloat, _ := c.Get("user_id")
	ownerID := uint(ownerIDFloat.(float64))

	properties, err := services.GetAllProperties(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch properties"})
		return
	}
	c.JSON(http.StatusOK, properties)
}
