package services

import (
	"errors"
	"fmt"
	"log"
	"pg-manager-backend/config"
	"pg-manager-backend/models"
	"pg-manager-backend/utils"
	"strings" // Required for splitting the reference_id
	"time"
)

func HandleRazorpayPayment(payload map[string]interface{}) error {
	// 1. Filter for the specific 'payment_link.paid' event
	event, _ := payload["event"].(string)
	if event != "payment_link.paid" {
		log.Printf("ℹ️ Webhook received for non-billable event: %s", event)
		return nil
	}

	// 2. Safely navigate the nested JSON payload
	data, ok := payload["payload"].(map[string]interface{})
	if !ok {
		return errors.New("invalid webhook payload: missing 'payload' field")
	}

	paymentLink, ok := data["payment_link"].(map[string]interface{})
	if !ok {
		return errors.New("invalid webhook payload: missing 'payment_link' field")
	}

	entity, ok := paymentLink["entity"].(map[string]interface{})
	if !ok {
		return errors.New("invalid webhook payload: missing 'entity' field")
	}

	// 3. Extract and Clean Tenant ID from reference_id (Format: BILL-ID-Timestamp)
	referenceID, _ := entity["reference_id"].(string)
	amountPaise, _ := entity["amount"].(float64)
	amountInRupees := amountPaise / 100

	// Split the reference_id "BILL-14-1771523347" to get "14"
	parts := strings.Split(referenceID, "-")
	if len(parts) < 2 {
		return fmt.Errorf("invalid reference_id format: %s", referenceID)
	}
	actualUserID := parts[1]

	log.Printf("💰 Razorpay Webhook: Received payment of ₹%.2f for Tenant UserID: %s", amountInRupees, actualUserID)

	// 4. Start Database Transaction
	tx := config.DB.Begin()

	var profile models.TenantProfile
	// Use the cleaned actualUserID to find the tenant
	if err := tx.Where("user_id = ?", actualUserID).First(&profile).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("tenant profile with UserID %s not found", actualUserID)
	}

	// Calculate and update the new balance
	newBalance := profile.Balance - amountInRupees
	if err := tx.Model(&profile).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update tenant balance")
	}

	// 5. Create a record in the Payment table
	paymentRecord := models.Payment{
		TenantID:    profile.UserID,
		PropertyID:  profile.PropertyID,
		Amount:      amountInRupees,
		PaymentType: "Rent-Payment",
		Method:      "Razorpay-Online",
		Date:        time.Now(),
	}

	if err := tx.Create(&paymentRecord).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to insert payment record into history")
	}

	// 6. Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.New("transaction commit failed")
	}

	// 7. Post-Payment Automation (Receipt & Terminal Log)
	go func() {
		log.Printf("✅ Payment successfully processed for %s.", profile.Name)

		fileName, err := utils.GenerateReceipt(paymentRecord, profile.Name)
		receiptURL := ""
		if err == nil {
			receiptURL = fmt.Sprintf("%s/receipts/%s", config.App.BaseURL, fileName)
		} else {
			log.Printf("⚠️ Receipt Generation Failed: %v", err)
		}

		// TERMINAL LOGGING (Simulating WhatsApp Message)
		fmt.Println("\n--- [TERMINAL WHATSAPP SIMULATION: PAYMENT SUCCESS] ---")
		fmt.Printf("To: %s (%s)\n", profile.Name, profile.PhoneNumber)
		fmt.Printf("Message: ✅ Payment Successful!\n"+
			"Namaste %s, your payment of ₹%.2f was received.\n"+
			"Current Balance: ₹%.2f.\n"+
			"Download Receipt: %s\n",
			profile.Name, amountInRupees, newBalance, receiptURL)

		// WhatsApp Part Commented Out as requested
		/*
			msg := fmt.Sprintf("✅ *Payment Successful!*\\n\\nNamaste %s, your payment of ₹%.2f via Razorpay was received.\\n*Current Balance:* ₹%.2f.\\n\\n📄 Download Receipt: %s",
				profile.Name, amountInRupees, newBalance, receiptURL)
			utils.SendWhatsAppMessage(profile.PhoneNumber, msg)
		*/
	}()

	return nil
}
