package services

import (
	"errors"
	"fmt"
	"pg-manager-backend/config"
	"pg-manager-backend/models"
	"pg-manager-backend/utils" // Uncommented
	"time"
)

func HandleRazorpayPayment(payload map[string]interface{}) error {
	data, ok := payload["payload"].(map[string]interface{})
	if !ok {
		return errors.New("invalid payload structure")
	}

	payment := data["payment"].(map[string]interface{})
	entity := payment["entity"].(map[string]interface{})

	amountInRupees := entity["amount"].(float64) / 100
	notes := entity["notes"].(map[string]interface{})
	tenantIDStr := fmt.Sprintf("%v", notes["tenant_id"])

	tx := config.DB.Begin()

	var profile models.TenantProfile
	if err := tx.Where("user_id = ?", tenantIDStr).First(&profile).Error; err != nil {
		tx.Rollback()
		return errors.New("tenant not found")
	}

	newBalance := profile.Balance - amountInRupees
	if err := tx.Model(&profile).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update balance")
	}

	paymentRecord := models.Payment{
		TenantID:    profile.UserID,
		PropertyID:  profile.PropertyID,
		Amount:      amountInRupees,
		PaymentType: "Online-Payment",
		Method:      "Razorpay",
		Date:        time.Now(),
	}

	if err := tx.Create(&paymentRecord).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to record payment")
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("finalize transaction failed")
	}

	// --- NEW: Generate Receipt & Send WhatsApp ---
	go func() {
		fileName, err := utils.GenerateReceipt(paymentRecord, profile.Name)
		if err == nil {
			receiptURL := config.App.BaseURL + "/receipts/" + fileName
			msg := fmt.Sprintf("✅ Payment Successful!\n\nNamaste %s, your payment of ₹%.2f via Razorpay was successful.\nUpdated Balance: ₹%.2f.\nDownload Receipt: %s",
				profile.Name, amountInRupees, newBalance, receiptURL)

			utils.SendWhatsAppMessage(profile.PhoneNumber, msg)
		}
	}()

	return nil
}
