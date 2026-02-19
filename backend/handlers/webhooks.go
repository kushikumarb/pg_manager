package handlers

import (
	"net/http"
	"pg-manager-backend/services"

	"github.com/gin-gonic/gin"
)

func RazorpayWebhook(c *gin.Context) {
	// In Test Mode, we read the JSON payload directly
	var payload map[string]interface{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Webhook Payload"})
		return
	}

	// Call the service to update balance and payment history
	err := services.HandleRazorpayPayment(payload)
	if err != nil {
		// We return 200 OK even on error to stop Razorpay from retrying
		// but we log the specific detail for debugging on your ASUS TUF
		c.JSON(http.StatusOK, gin.H{
			"status":  "processed_with_error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "payment_processed_successfully"})
}
