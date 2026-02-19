package services

import (
	"fmt"
	"pg-manager-backend/config"
	"pg-manager-backend/models"
	"pg-manager-backend/utils"
	"time"
	// "pg-manager-backend/utils"
)

// RecordManualPayment handles Cash, UPI, or Bank transfers recorded by the owner
func RecordManualPayment(userID uint, amount float64, method string) (float64, error) {
	tx := config.DB.Begin()

	var profile models.TenantProfile
	if err := tx.First(&profile, "user_id = ?", userID).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Balance logic
	newBalance := profile.Balance - amount
	if err := tx.Model(&profile).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	paymentRecord := models.Payment{
		TenantID:    profile.UserID,
		PropertyID:  profile.PropertyID,
		Amount:      amount,
		PaymentType: "Manual-Payment",
		Method:      method,
		Date:        time.Now(),
	}

	if err := tx.Create(&paymentRecord).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	// --- NEW: Generate Receipt & Send WhatsApp ---
	go func() {
		fileName, err := utils.GenerateReceipt(paymentRecord, profile.Name)
		if err == nil {
			receiptURL := config.App.BaseURL + "/receipts/" + fileName
			msg := fmt.Sprintf("✅ Payment Received!\n\nNamaste %s, we received ₹%.2f via %s.\nNew Balance: ₹%.2f\nDownload Receipt: %s",
				profile.Name, amount, method, newBalance, receiptURL)

			utils.SendWhatsAppMessage(profile.PhoneNumber, msg)
		}
	}()

	return newBalance, nil
}

type PaymentResponse struct {
	models.Payment
	TenantName string `json:"tenant_name"`
}

func GetAllPaymentHistory() ([]PaymentResponse, error) {
	var results []PaymentResponse

	// Get all payments
	config.DB.Table("payments").Order("date desc").Scan(&results)

	for i := range results {
		var name string
		// Try active table first
		config.DB.Table("tenant_profiles").
			Where("user_id = ?", results[i].TenantID).
			Select("name").Scan(&name)

		// If not found, check the backup table
		if name == "" {
			config.DB.Table("archived_tenants").
				Where("original_user_id = ?", results[i].TenantID).
				Select("name").Scan(&name)
		}

		if name != "" {
			results[i].TenantName = name
		} else {
			results[i].TenantName = "Unknown (ID: " + fmt.Sprint(results[i].TenantID) + ")"
		}
	}
	return results, nil
}
