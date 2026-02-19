package handlers

import (
	"fmt"
	"net/http"
	"pg-manager-backend/models"
	"pg-manager-backend/services"

	"github.com/gin-gonic/gin"
)

func OnboardTenant(c *gin.Context) {
	var input models.TenantProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	tenantID, err := services.OnboardTenant(input)
	if err != nil {
		// Professional tip: Use a status code that reflects the error (like Conflict if room is full)
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Tenant registration initiated. OTP sent to tenant.",
		"tenant_id": tenantID,
		"status":    "pending",
	})
}

// ConfirmAdmission is the new handler for Point 5 of your plan
func ConfirmAdmission(c *gin.Context) {
	var input struct {
		TenantID uint   `json:"tenant_id" binding:"required"`
		OTP      string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID and OTP are required"})
		return
	}

	err := services.VerifyTenantOTP(input.TenantID, input.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admission confirmed! Welcome message and payment link sent to tenant."})
}

func OffboardTenant(c *gin.Context) {
	tenantID := c.Param("id")

	err := services.OffboardTenant(tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant Removed. Chatbot access revoked"})
}

func CheckBalance(c *gin.Context) {
	// Get the UserID from the JWT token (set by middleware)
	userID, _ := c.Get("userID")

	balance, err := services.GetTenantBalance(userID.(uint))
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": balance,
		"message": fmt.Sprintf("Your current outstanding balance is ₹%.2f", balance),
	})
}

func GetTenants(c *gin.Context) {
	// Extract property_id from the query string (?property_id=1)
	propertyID := c.Query("property_id")

	if propertyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "property_id is required"})
		return
	}

	tenants, err := services.GetTenantsByProperty(propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tenants"})
		return
	}

	c.JSON(http.StatusOK, tenants)
}

func GetTenantProfile(c *gin.Context) {
	id := c.Param("id")
	profile, err := services.GetTenantByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func OffboardTenantHandler(c *gin.Context) {
	tenantID := c.Param("id")

	if err := services.OffboardTenant(tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant offboarded successfully"})
}

func GetArchivedTenantsHandler(c *gin.Context) {
	propertyID := c.Query("property_id")
	if propertyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Property ID is required"})
		return
	}

	archives, err := services.GetArchivedTenants(propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch backup data"})
		return
	}

	c.JSON(http.StatusOK, archives)
}
