package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log" // Used for terminal verification
	"pg-manager-backend/config"
	"pg-manager-backend/models"
	"pg-manager-backend/utils"

	// "pg-manager-backend/utils" // Commented out: No active credentials
	"time"

	"gorm.io/gorm"
)

// Random OTP Generator remains the same
func generateOTP() string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, 6)
	io.ReadAtLeast(rand.Reader, b, 6)
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// OnboardTenant handles initial registration and OTP dispatch
func OnboardTenant(input models.TenantProfile) (uint, error) {
	tx := config.DB.Begin()

	// 1. ROOM CAPACITY CHECK
	var room models.Room
	if err := tx.First(&room, input.RoomID).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("room not found")
	}

	var activeTenants int64
	tx.Model(&models.TenantProfile{}).Where("room_id = ? AND status = ?", input.RoomID, "active").Count(&activeTenants)

	if activeTenants >= int64(room.Capacity) {
		tx.Rollback()
		return 0, fmt.Errorf("room %s is full", room.RoomNumber)
	}

	// 2. CREATE USER ENTITY
	newUser := models.User{
		Name:  input.Name,
		Phone: input.PhoneNumber,
		Role:  "tenant",
	}
	// This will fail if the phone number already exists in the 'users' table
	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 3. SETUP FULL PROFILE
	otpCode := generateOTP()
	input.UserID = newUser.ID
	input.OTP = otpCode
	input.IsVerified = false
	input.Status = "pending"
	input.AdmissionDate = time.Now()
	input.NextBillingDate = input.AdmissionDate.AddDate(0, 1, 0)

	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	// 4. TERMINAL LOGGING
	log.Printf("------------------------------------------------")
	log.Printf("ADMISSION KYC COMPLETE - OTP GENERATED")
	log.Printf("Tenant: %s | Phone: %s", input.Name, input.PhoneNumber)
	log.Printf("OTP Code: %s", otpCode)
	log.Printf("------------------------------------------------")

	return newUser.ID, nil
}

// VerifyTenantOTP confirms the OTP and activates the tenant
// ... existing imports ...

func VerifyTenantOTP(tenantID uint, inputOTP string) error {
	var profile models.TenantProfile
	if err := config.DB.Preload("Room").Where("user_id = ?", tenantID).First(&profile).Error; err != nil {
		return errors.New("tenant profile not found")
	}

	if profile.OTP != inputOTP {
		return errors.New("invalid OTP: verification failed")
	}

	// Calculate Initial Total: Monthly Rent + Security Deposit
	initialDue := profile.MonthlyRent + profile.Room.Deposit

	updates := map[string]interface{}{
		"is_verified": true,
		"status":      "active",
		"balance":     initialDue, // Set initial balance
	}

	if err := config.DB.Model(&profile).Updates(updates).Error; err != nil {
		return err
	}

	// 🔗 GENERATE INITIAL LINK (Rent + Deposit)
	paymentLink := utils.GenerateRazorpayLink(profile.UserID, initialDue, "Initial Rent + Deposit")

	// TERMINAL LOGGING (Simulating WhatsApp)
	log.Printf("\n--- WHATSAPP SIMULATION (OTP VERIFIED) ---")
	log.Printf("To: %s (%s)", profile.Name, profile.PhoneNumber)
	log.Printf("Message: ✅ Admission Confirmed! Your initial total (Rent+Deposit) is ₹%.2f. Pay here: %s", initialDue, paymentLink)
	log.Printf("------------------------------------------\n")

	// Commented out: WhatsApp notification
	/*
		welcomeMsg := fmt.Sprintf("✅ *Admission Confirmed!*\n🔗 Pay here: %s", paymentLink)
		utils.SendWhatsAppMessage(profile.PhoneNumber, welcomeMsg)
	*/

	return nil
}

// GetTenantBalance includes terminal logging for balance checks
func GetTenantBalance(userID uint) (float64, error) {
	var profile models.TenantProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return 0, errors.New("tenant profile not found")
	}

	if profile.LastBalanceCheck != nil {
		lastCheck := profile.LastBalanceCheck.Format("2006-01-02")
		today := time.Now().Format("2006-01-02")
		if lastCheck == today {
			return 0, errors.New("limit reached: you can only check your balance once per day")
		}
	}

	now := time.Now()
	config.DB.Model(&profile).Update("last_balance_check", &now)

	// TERMINAL LOGGING: Simulate balance inquiry
	log.Printf("💰 Balance Inquiry: Tenant %s | Balance: ₹%.2f", profile.Name, profile.Balance)

	return profile.Balance, nil
}

// GetTenantsByProperty remains the same
func GetTenantsByProperty(propertyID string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := config.DB.Table("tenant_profiles").
		// UPDATED: Changed users.phone to tenant_profiles.phone_number
		Select("tenant_profiles.user_id, users.name, rooms.room_no AS room_no, tenant_profiles.phone_number, tenant_profiles.status").
		Joins("JOIN users ON users.id = tenant_profiles.user_id").
		Joins("JOIN rooms ON rooms.id = tenant_profiles.room_id").
		Where("rooms.property_id = ?", propertyID).
		Scan(&results).Error

	return results, err
}

func GetTenantByID(id string) (models.TenantProfile, error) {
	var profile models.TenantProfile
	// We use Preload to get the Room details and user details in one shot
	err := config.DB.Preload("Room").First(&profile, "user_id =?", id).Error
	return profile, err
}

func OffboardTenant(tenantID string) error {
	var profile models.TenantProfile
	if err := config.DB.Where("user_id = ?", tenantID).First(&profile).Error; err != nil {
		return errors.New("tenant not found")
	}

	if profile.Balance > 0 {
		return fmt.Errorf("cannot offboard: pending balance ₹%.2f", profile.Balance)
	}

	// Use a transaction to ensure either everything happens or nothing happens
	return config.DB.Transaction(func(tx *gorm.DB) error {

		// 1. MOVE DATA TO BACKUP TABLE
		archive := models.ArchivedTenant{
			OriginalUserID:   profile.UserID,
			Name:             profile.Name,
			PhoneNumber:      profile.PhoneNumber,
			FatherName:       profile.FatherName,
			PermanentAddress: profile.PermanentAddress,
			IDProofNo:        profile.IDProofNo,
			IDProofImage:     profile.IDProofImage,
			PropertyID:       profile.PropertyID,
			RoomID:           profile.RoomID,
			AdmissionDate:    profile.AdmissionDate,
			CheckoutDate:     time.Now(),
		}

		if err := tx.Create(&archive).Error; err != nil {
			return fmt.Errorf("failed to backup data: %v", err)
		}

		// 2. DELETE ACTIVE RECORDS
		// Logic: In Postgres, delete the 'Profile' first, then the 'User'.
		// This respects Foreign Key constraints without needing to disable them.

		// Delete Profile (Child record)
		if err := tx.Unscoped().Where("user_id = ?", tenantID).Delete(&models.TenantProfile{}).Error; err != nil {
			return err
		}

		// Delete User (Parent record)
		if err := tx.Unscoped().Where("id = ?", tenantID).Delete(&models.User{}).Error; err != nil {
			return err
		}

		return nil // Transaction will commit here
	})
}

// GetArchivedTenants retrieves all records from the backup table
func GetArchivedTenants(propertyID string) ([]models.ArchivedTenant, error) {
	var archives []models.ArchivedTenant

	// We filter by PropertyID so the owner only sees history for their PG
	err := config.DB.Where("property_id = ?", propertyID).
		Order("checkout_date desc").
		Find(&archives).Error

	return archives, err
}
