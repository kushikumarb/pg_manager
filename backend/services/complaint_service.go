package services

import (
	"errors"
	"pg-manager-backend/config"
	"pg-manager-backend/models"
	"time"
)

// RegisterComplaint saves the data from the QR form
func RegisterComplaint(data models.Complaint) error {
	var count int64

	// 1. Check if the phone number belongs to an active tenant in THIS property
	config.DB.Model(&models.TenantProfile{}).
		Where("phone_number = ? AND property_id = ?", data.PhoneNumber, data.PropertyID).
		Count(&count)

	if count == 0 {
		// Return a specific error if verification fails
		return errors.New("tenant_not_verified")
	}

	// 2. If verified, save the complaint
	data.CreatedAt = time.Now()
	data.Status = "Pending"
	return config.DB.Create(&data).Error
}

// GetAllComplaints returns complaints for the owner's dashboard
func GetAllComplaints(propertyID string) ([]models.Complaint, error) {
	var complaints []models.Complaint
	// Order by CreatedAt desc so the owner sees new issues first
	err := config.DB.Where("property_id = ?", propertyID).Order("created_at desc").Find(&complaints).Error
	return complaints, err
}

// ResolveComplaint updates the status in the DB
func ResolveComplaint(id string) error {
	return config.DB.Model(&models.Complaint{}).Where("id = ?", id).Update("status", "Resolved").Error
}
