package services

import (
	"pg-manager-backend/config"
	"pg-manager-backend/models"
	"time"
)

func RecordExpense(expense models.Expenditure) error {
	if expense.Date.IsZero() {
		expense.Date = time.Now()
	}
	return config.DB.Create(&expense).Error
}

func GetExpendituresByProperty(propertyID string) ([]models.Expenditure, error) {
	var expenses []models.Expenditure
	err := config.DB.Where("property_id = ?", propertyID).Order("date desc").Find(&expenses).Error
	return expenses, err
}

func GetMonthlyFinanceSummary(propertyID string) (float64, float64, error) {
	var totalExpense float64
	
	// Get first day of the current month
	currentMonth := time.Now().Format("2006-01") + "-01"

	// COALESCE ensures we return 0 if no records exist
	config.DB.Model(&models.Expenditure{}).
		Where("property_id = ? AND date >= ?", propertyID, currentMonth).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalExpense)

	return 0, totalExpense, nil
}