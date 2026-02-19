package handlers

import (
	"fmt"
	"net/http"
	"pg-manager-backend/services"

	"github.com/gin-gonic/gin"
)

func RecordPayment(c *gin.Context) {
	// 1. Capture userID from URL parameter (:id)
	userIDStr := c.Param("id")

	// 2. Define struct to match the Frontend JSON exactly
	var input struct {
		Amount float64 `json:"amount" binding:"required"`
		Method string  `json:"method" binding:"required"`
	}

	// Bind JSON Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Convert string ID to uint
	var userID uint
	fmt.Sscanf(userIDStr, "%d", &userID)

	// 3. Call Service Layer
	newBalance, err := services.RecordManualPayment(userID, input.Amount, input.Method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Payment recorded successfully",
		"new_balance": newBalance,
	})
}

func GetPaymentHistory(c *gin.Context) {
	// Call the service logic
	results, err := services.GetAllPaymentHistory()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, results)
}
