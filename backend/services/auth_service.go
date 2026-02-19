package services

import (
	"errors"
	"os"
	"pg-manager-backend/config"
	"pg-manager-backend/models"
	
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// --- NEW: Exported Hashing Functions ---

// HashPassword is now accessible by config.seedOwner
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash is a clean way to verify passwords
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// --- End of New Functions ---

func RegisterUser(name, email, password, role string) error {
	// Use the new helper function
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     name,
		Email:    &email,
		Password: hashedPassword,
		Role:     role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return errors.New("user already exists or database error")
	}

	return nil
}

func AuthenticateUser(email, password string) (string, string, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", errors.New("invalid email or password")
	}

	// Use the new helper function
	if !CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "my_secret_key"
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return tokenString, user.Role, nil
}
