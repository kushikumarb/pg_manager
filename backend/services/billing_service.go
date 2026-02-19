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

	// 1. Find active tenants due for billing
	if err := config.DB.Where("status = ? AND next_billing_date <= ?", "active", today).Find(&tenants).Error; err != nil {
		log.Printf("Error fetching due tenants: %v", err)
		return
	}

	for _, tenant := range tenants {
		tx := config.DB.Begin()

		// 2. Update Balance and Next Billing Date (30 days)
		newBalance := tenant.Balance + tenant.MonthlyRent
		newNextBillingDate := tenant.NextBillingDate.AddDate(0, 0, 30)

		updates := map[string]interface{}{
			"balance":           newBalance,
			"last_billing_date": today,
			"next_billing_date": newNextBillingDate,
		}

		if err := tx.Model(&tenant).Updates(updates).Error; err != nil {
			tx.Rollback()
			continue
		}
		tx.Commit()

		// 3. Generate Simulated Link
		paymentLink := utils.GenerateRazorpayLink(tenant.UserID, newBalance,"Monthly-Rent")

		// 4. TERMINAL LOGGING instead of WhatsApp
		fmt.Println("\n--- [TERMINAL WHATSAPP SIMULATION: MONTHLY BILL] ---")
		fmt.Printf("To: %s (%s)\n", tenant.Name, tenant.PhoneNumber)
		fmt.Printf("Message: Namaste %s! 🏠\n"+
			"Your monthly rent of ₹%.2f is due. Total Balance: ₹%.2f.\n"+
			"Click here to pay: %s\n",
			tenant.Name, tenant.MonthlyRent, newBalance, paymentLink)

		/* Commented out for now
		utils.SendWhatsAppMessage(tenant.PhoneNumber, message)
		*/
	}
}
