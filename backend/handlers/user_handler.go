package handlers

import (
	"net/http"
	"pg-manager-backend/services"

	"github.com/gin-gonic/gin"
)

// GetOwnerProfile handles fetching current owner details
func GetOwnerProfile(c *gin.Context) {
	// 1. Change this to "user_id" to match your working middleware
	val, exists := c.Get("user_id")
	if !exists {
		// If this executes, it means the middleware didn't set the key correctly
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session context missing"})
		return
	}

	var userID uint
	// 2. Safe extraction from the interface{}
	switch v := val.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID format"})
		return
	}

	// 3. Call the service
	profile, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateOwnerProfile handles updating owner name/email
func UpdateOwnerProfile(c *gin.Context) {
    // 1. Use "user_id" to match your working Dashboard/Property routes
    val, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Session context missing"})
        return
    }

    // 2. SAFE TYPE EXTRACTION (The logic from GetOwnerProfile)
    var userID uint
    switch v := val.(type) {
    case float64:
        userID = uint(v)
    case uint:
        userID = v
    default:
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format in session"})
        return
    }

    // 3. Bind the incoming JSON from Vue
    var input struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
        return
    }

    // 4. Call the service to update the database
    err := services.UpdateUser(userID, input.Name, input.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}