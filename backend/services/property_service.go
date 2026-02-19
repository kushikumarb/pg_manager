package services

import (
	"errors"
	"pg-manager-backend/config"
	"pg-manager-backend/models"
)

// CreateProperty logic remains the same
func CreateProperty(name, address string, ownerID uint) (models.Property, error) {
	property := models.Property{
		Name:    name,
		Address: address,
		OwnerID: ownerID,
	}
	if err := config.DB.Create(&property).Error; err != nil {
		return models.Property{}, errors.New("could not create property in database")
	}
	return property, nil
}

// GetPropertyStats - UPDATED to include expenditure
func GetPropertyStats(propertyID string) (map[string]interface{}, error) {
	var roomCount, tenantCount, complaintCount int64

	// 1. Total Rooms
	config.DB.Model(&models.Room{}).Where("property_id = ?", propertyID).Count(&roomCount)

	// 2. Active Tenants
	config.DB.Table("tenant_profiles").
		Joins("JOIN rooms ON rooms.id = tenant_profiles.room_id").
		Where("rooms.property_id = ? AND tenant_profiles.status = ?", propertyID, "active").
		Count(&tenantCount)

	// 3. Pending Complaints
	config.DB.Table("complaints").
		Joins("JOIN tenant_profiles ON tenant_profiles.user_id = complaints.tenant_id").
		Joins("JOIN rooms ON rooms.id = tenant_profiles.room_id").
		Where("rooms.property_id = ? AND complaints.status = ?", propertyID, "pending").
		Count(&complaintCount)

	// 4. ADDED: Monthly Expenditure logic
	_, totalExpense, err := GetMonthlyFinanceSummary(propertyID)
	if err != nil {
		totalExpense = 0
	}

	return map[string]interface{}{
		"total_rooms":       roomCount,
		"active_tenants":    tenantCount,
		"pending_issues":    complaintCount,
		"total_expenditure": totalExpense, // This matches stats.total_expenditure in Vue
	}, nil
}

// GetDashboardData (Global Stats) remains same
func GetDashboardData(propertyID string) (map[string]interface{}, error) {
	var totalRooms, activeTenants, pendingIssues int64
	var totalRevenue, totalExpenditure float64

	// 1. Total Rooms
	config.DB.Model(&models.Room{}).Where("property_id = ?", propertyID).Count(&totalRooms)

	// 2. Active Tenants
	config.DB.Model(&models.TenantProfile{}).Where("property_id = ?", propertyID).Count(&activeTenants)

	// 3. UPDATED: Pending Complaints - QR System simplified! 🚀
	// Since we added property_id directly to the Complaint model, no JOIN is needed.
	config.DB.Model(&models.Complaint{}).
		Where("property_id = ? AND status = ?", propertyID, "Pending").
		Count(&pendingIssues)

	// 4. Total Revenue (Manual + Online)
	config.DB.Model(&models.Payment{}).
		Where("property_id = ?", propertyID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalRevenue)

	// 5. Total Expenditure
	config.DB.Model(&models.Expenditure{}).
		Where("property_id = ?", propertyID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpenditure)

	return map[string]interface{}{
		"total_rooms":       totalRooms,
		"active_tenants":    activeTenants,
		"pending_issues":    pendingIssues, // Matches pending_issues in Vue
		"total_revenue":     totalRevenue,
		"total_expenditure": totalExpenditure,
	}, nil
}

func GetAllProperties(ownerID uint) ([]models.Property, error) {
	var properties []models.Property
	err := config.DB.Where("owner_id = ?", ownerID).Find(&properties).Error
	return properties, err
}

func GetRoomsByPropertyID(propertyID string) ([]models.Room, error) {
	var rooms []models.Room
	err := config.DB.Where("property_id = ?", propertyID).Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	for i := range rooms {
		var count int64
		config.DB.Model(&models.TenantProfile{}).Where("room_id = ? AND status = ?", rooms[i].ID, "active").Count(&count)
		rooms[i].Occupied = int(count)
	}
	return rooms, nil
}
