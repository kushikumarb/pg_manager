package models

import (
	"time"

	"gorm.io/gorm"
)

// User handles both Owners and Tenants
type User struct {
	gorm.Model
	Name     string  `json:"name"`
	Email    *string `json:"email" gorm:"unique;default:null"`
	Password string  `json:"-"` // "-" ensures password never leaves the backend
	Role     string  `json:"role"`
	Phone    string  `json:"phone" gorm:"unique;not null"`
}

// Property represents a PG Building
type Property struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name"`
	Address   string         `json:"address"`
	OwnerID   uint           `json:"owner_id"`
	Rooms     []Room         `json:"rooms" gorm:"foreignKey:PropertyID"`
}

// Room represents an individual room
type Room struct {
	gorm.Model
	PropertyID uint    `json:"property_id"`
	RoomNumber string  `json:"room_no" gorm:"column:room_no"`
	Capacity   int     `json:"capacity"`
	Price      float64 `json:"price"`   // Rent per bed
	Deposit    float64 `json:"deposit"` // FIXED DEPOSIT for the room
	Occupied   int     `json:"occupied" gorm:"-"`
	IsFull     bool    `json:"is_full" gorm:"default:false"`
}

// TenantProfile links a user to a room and tracks details
type TenantProfile struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	PropertyID  uint   `json:"property_id"`
	User        User   `json:"user" gorm:"foreignKey:UserID"` // Allows .Preload("User")
	RoomID      uint   `json:"room_id"`
	Room        Room   `json:"room" gorm:"foreignKey:RoomID"` // Allows .Preload("Room")
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" gorm:"unique;not null"`
	Status      string `json:"status" gorm:"default:active"`

	// Personal Details
	FatherName        string `json:"father_name"`
	DOB               string `json:"dob"`
	Age               int    `json:"age"`
	PermanentAddress  string `json:"permanent_address"`
	EmergencyContact  string `json:"emergency_contact"`
	EmergencyRelation string `json:"emergency_relation"`
	MailID            string `json:"mail_id"`

	// Preferences & Logistics
	IsVegetarian  bool   `json:"is_vegetarian"`
	HasTwoWheeler bool   `json:"has_two_wheeler"`
	VehicleNo     string `json:"vehicle_no"`
	Education     string `json:"education"`
	Occupation    string `json:"occupation"`
	OfficeAddress string `json:"office_address"`
	IDProofType   string `json:"id_proof_type"`
	IDProofNo     string `json:"id_proof_no"`
	IDProofImage  string `json:"id_proof_image" gorm:"type:text"`

	// Financials
	MonthlyRent        float64    `json:"monthly_rent"`
	Deposit            float64    `json:"deposit"`
	MaintenanceCharges float64    `json:"maintenance_charges"`
	Balance            float64    `json:"balance" gorm:"default:0"`
	AdmissionDate      time.Time  `json:"admission_date"`
	LastBillingDate    *time.Time `json:"last_billed_date"`
	NextBillingDate    time.Time  `json:"next_billing_date" gorm:"index"` //Index for fast daily lookups

	//Verification for Digital Signature
	OTP        string `json:"otp"`
	IsVerified bool   `json:"is_verified" gorm:"default:false"`

	//Complain Fields
	LastComplaintDate *time.Time `json:"last_complaint_date"`
	LastBalanceCheck  *time.Time `json:"last_balance_check"`
}

// Complaint Model
type Complaint struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PropertyID  uint      `json:"property_id"`
	RoomNo      string    `json:"room_no"`
	TenantName  string    `json:"tenant_name"`
	PhoneNumber string    `json:"phone_number"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"default:'Pending'"`
	CreatedAt   time.Time `json:"created_at"`
}

// Expenditure Model
type Expenditure struct {
	ID          uint      `gorm:"primaryKey"`
	PropertyID  uint      `json:"property_id" binding:"required"`
	Amount      float64   `json:"amount" binding:"required"`
	Category    string    `json:"category" binding:"required"` // Can be "Electricity" or "New Lightbulbs"
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type Payment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PropertyID  uint      `json:"property_id"`
	TenantID    uint      `json:"tenant_id"`
	Amount      float64   `json:"amount"`
	PaymentType string    `json:"payment_type"` // e.g., "Rent", "Deposit", "Maintenance"
	Method      string    `json:"method"`       // e.g., "Cash", "UPI", "Bank Transfer"
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

type ArchivedTenant struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	OriginalUserID   uint      `json:"original_user_id"`
	Name             string    `json:"name"`
	PhoneNumber      string    `json:"phone_number"`
	FatherName       string    `json:"father_name"`
	PermanentAddress string    `json:"permanent_address"`
	IDProofNo        string    `json:"id_proof_no"`
	IDProofImage     string    `gorm:"type:text" json:"id_proof_image"`
	PropertyID       uint      `json:"property_id"`
	RoomID           uint      `json:"room_id"`
	AdmissionDate    time.Time `json:"admission_date"`
	CheckoutDate     time.Time `json:"checkout_date"`
}
