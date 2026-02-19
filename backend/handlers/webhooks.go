package handlers

import(
	"net/http"
	"pg-manager-backend/services"
	"github.com/gin-gonic/gin"
)

func RazorpayWebhook(c *gin.Context) {
	// In a real product, you'd verify the Razorpay signature here for security
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err !=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Payload"})
		return
	}

	// Logic to Extract payment ID and Tenant ID from the Webhook
	err := services.HandleRazorpayPayment(payload)
	if err != nil{
		c.JSON(http.StatusOK, gin.H{"status": "processed_with_error", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}