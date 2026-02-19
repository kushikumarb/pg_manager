package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	razorpay "github.com/razorpay/razorpay-go"
)

func GenerateRazorpayLink(userID uint, email string, amount float64, description string) (string, error) {
	// 1. Check for Credentials
	keyID := os.Getenv("RAZORPAY_KEY_ID")
	secret := os.Getenv("RAZORPAY_KEY_SECRET")

	if keyID == "" || secret == "" {
		return "", fmt.Errorf("razorpay credentials missing in environment")
	}

	client := razorpay.NewClient(keyID, secret)

	uniqueRef := fmt.Sprintf("BILL-%d-%d", userID, time.Now().Unix())

	// 2. Prepare Data
	data := map[string]interface{}{
		"amount":       int64(amount * 100),
		"currency":     "INR",
		"description":  description,
		"reference_id": uniqueRef,
		"customer": map[string]interface{}{
			"name":  fmt.Sprintf("Tenant ID: %d", userID),
			"email": email,
		},
		"notify": map[string]interface{}{
			"sms":   false,
			"email": true,
		},
		// 🔗 CRITICAL: This is what your Webhook Service uses to find the tenant!
		"notes": map[string]interface{}{
			"tenant_id": fmt.Sprint(userID),
		},
	}

	// 3. Create the Link
	body, err := client.PaymentLink.Create(data, nil)
	if err != nil {
		log.Printf("❌ Razorpay API Error: %v", err)
		return "", err
	}

	// 4. Extract Short URL
	if shortURL, ok := body["short_url"].(string); ok {
		return shortURL, nil
	}

	return "", fmt.Errorf("failed to parse short_url from Razorpay response")
}