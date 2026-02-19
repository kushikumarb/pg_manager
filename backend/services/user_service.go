package services

import (
	"pg-manager-backend/config"
	"pg-manager-backend/models"
)

// GetUserByID fetches a user but hides the password for security
func GetUserByID(id uint) (models.User, error) {
	var user models.User
	// We search by the ID inherited from gorm.Model
	err := config.DB.Omit("Password").First(&user, id).Error
	return user, err
}

func UpdateUser(id uint, name string, email string) error {
	// Since Email in your model is *string, we pass the address of the string
	return config.DB.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":  name,
		"email": &email, // Pass pointer for the *string field
	}).Error
}
