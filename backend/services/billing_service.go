package services

import (
	"fmt"
	"log"
	"pg-manager-backend/config"
	"pg-manager-backend/models"
	"pg-manager-backend/utils"
	"time"
)

func ProcessDailyBilling() {
	var tenants []models.TenantProfile
	today := time.Now().Truncate(24 * time.Hour)

	if err := config.DB.Where("status = ? AND next_billing_date <= ?", "active", today).Find(&tenants).Error; err != nil {
		log.Printf("Error fetching due tenants: %v", err)
		return
	}

	for _, tenant := range tenants {
		tx := config.DB.Begin()

		newBalance := tenant.Balance + tenant.MonthlyRent
		newNextBillingDate := tenant.NextBillingDate.AddDate(0, 0, 30)

		updates := map[string]interface{}{
			"balance":           newBalance,
			"last_billed_date":  &today, // Using pointer as per your struct
			"next_billing_date": newNextBillingDate,
		}

		if err := tx.Model(&tenant).Updates(updates).Error; err != nil {
			tx.Rollback()
			continue
		}
		tx.Commit()

		// Generate Link using MailID
		paymentLink, err := utils.GenerateRazorpayLink(tenant.UserID, tenant.MailID, newBalance, "Monthly Rent")
		if err != nil {
			log.Printf("⚠️ Billing link failed for %s: %v", tenant.Name, err)
			paymentLink = "[Link Unavailable]"
		}

		// TERMINAL LOGGING
		fmt.Println("\n--- [TERMINAL WHATSAPP SIMULATION: MONTHLY BILL] ---")
		fmt.Printf("To: %s (%s)\n", tenant.Name, tenant.PhoneNumber)
		fmt.Printf("Message: Namaste %s! 🏠\n"+
			"Your monthly rent of ₹%.2f is due. Total Balance: ₹%.2f.\n"+
			"Click here to pay: %s\n",
			tenant.Name, tenant.MonthlyRent, newBalance, paymentLink)

		// KEEPING THIS COMMENTED AS REQUESTED
		/*
		   message := fmt.Sprintf("Namaste %s! Your rent is due. Pay here: %s", tenant.Name, paymentLink)
		   utils.SendWhatsAppMessage(tenant.PhoneNumber, message)
		*/
	}
}
