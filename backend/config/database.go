package config

import (
	"fmt"
	"log"
	"pg-manager-backend/models"

	"gorm.io/driver/postgres" // Changed from mysql
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// 1. PostgreSQL DSN Format is different from MySQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		App.DBHost, App.DBUser, App.DBPass, App.DBName, App.DBPort)

	// 2. Open connection using the Postgres driver
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("❌ Failed to connect to database: ", err)
	}

	// 3. Auto Migrate
	fmt.Println("🚀 Running migrations...")
	err = database.AutoMigrate(
		&models.User{},
		&models.Property{},
		&models.Room{},
		&models.TenantProfile{},
		&models.Complaint{},
		&models.Expenditure{},
		&models.Payment{},
		&models.ArchivedTenant{},
	)

	if err != nil {
		log.Fatal("❌ Migration Error:", err)
	}

	DB = database
	fmt.Println("✅ Database connection and migrations successful")
}
