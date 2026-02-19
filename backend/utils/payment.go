package utils

import "fmt"

// GenerateRazorpayLink creates a payment link.
// Once you have credentials, this will call the Razorpay API.
func GenerateRazorpayLink(userID uint, amount float64, description string) string {
	// DUMMY LINK: In production, this becomes a real API call to Razorpay
	// Example: return razorpayClient.PaymentLink.Create(params)
	return fmt.Sprintf("https://rzp.io/i/pay_u%d_amt%.0f", userID, amount)
}