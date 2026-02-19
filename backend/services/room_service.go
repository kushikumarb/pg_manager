package services

import (
	"errors"
	"pg-manager-backend/config"
	"pg-manager-backend/models"
)

// CreateRoom logic with ownership validation
func CreateRoom(propertyID, ownerID uint, roomNumber string, capacity int, price float64, deposit float64) (models.Room, error) {
	// 1. Security check: Verify property ownership
	var property models.Property
	if err := config.DB.Where("id = ? AND owner_id = ?", propertyID, ownerID).First(&property).Error; err != nil {
		return models.Room{}, errors.New("unauthorized: you do not own this property")
	}

	// 2. Prepare Room Object
	room := models.Room{
		PropertyID: propertyID,
		RoomNumber: roomNumber,
		Capacity:   capacity,
		Price:      price,
		Deposit:    deposit, // ADD THIS LINE
		IsFull:     false,
	}

	// 3. Save to Database
	if err := config.DB.Create(&room).Error; err != nil {
		return models.Room{}, errors.New("failed to create room in database")
	}

	return room, nil
}


func DeleteRoom(roomID string) error {
	// 1. Check for active tenants in the room
	var activeCount int64
	err := config.DB.Model(&models.TenantProfile{}).
		Where("room_id = ? AND status = ?", roomID, "active").
		Count(&activeCount).Error

	if err != nil {
		return err
	}

	// 2. Block deletion if occupancy is greater than zero
	if activeCount > 0 {
		return errors.New("cannot delete room: active tenants are currently assigned to it")
	}

	// 3. Perform the delete (Soft delete if gorm.Model is used)
	return config.DB.Delete(&models.Room{}, roomID).Error
}
