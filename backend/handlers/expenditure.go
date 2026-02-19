package handlers

import (
	"net/http"
	"pg-manager-backend/models"
	"pg-manager-backend/services"

	"github.com/gin-gonic/gin"
)

func AddExpense(c *gin.Context) {
	var input models.Expenditure
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.RecordExpense(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record expense"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Expense recorded successfully"})
}

// NEW: Handler to get expenses
func GetExpenditures(c *gin.Context) {
	propertyID := c.Query("property_id")
	expenses, err := services.GetExpendituresByProperty(propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}
	c.JSON(http.StatusOK, expenses)
}
