package handlers

import (
	"net/http"
	"pg-manager-backend/services"

	"github.com/gin-gonic/gin"
)

func AddRoom(c *gin.Context) {
	var input struct {
		PropertyID uint    `json:"property_id" binding:"required"`
		RoomNumber string  `json:"room_no" binding:"required"`
		Capacity   int     `json:"capacity" binding:"required"`
		Price      float64 `json:"price" binding:"required"`
		Deposit    float64 `json:"deposit" binding:"required"` // ADD THIS LINE
	}

	// 1. Validate JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Get Owner ID from the JWT context
	ownerIDFloat, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User context not found"})
		return
	}
	ownerID := uint(ownerIDFloat.(float64))

	// 3. Call Service Layer
	room, err := services.CreateRoom(input.PropertyID, ownerID, input.RoomNumber, input.Capacity, input.Price, input.Deposit)
	if err != nil {
		// Distinguish between unauthorized and server errors
		if err.Error() == "unauthorized: you do not own this property" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room added successfully!",
		"room":    room,
	})
}

func RemoveRoom(c *gin.Context) {
	// Extract the room ID from the URL parameter
	roomID := c.Param("id")

	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	// Call the service layer
	err := services.DeleteRoom(roomID)
	if err != nil {
		// Handle the specific safety check error
		if err.Error() == "cannot delete room: active tenants are currently assigned to it" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}
