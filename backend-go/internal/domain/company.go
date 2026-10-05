package domain

import (
	"time"

	"github.com/google/uuid"
)

type EntityType string

const (
	EntityPabrik  EntityType = "PABRIK"
	EntityLapak   EntityType = "LAPAK"
	EntitySmelter EntityType = "SMELTER"
	EntityFunder  EntityType = "FUNDER"
)

type KYCStatus string

const (
	KYCPending  KYCStatus = "PENDING"
	KYCVerified KYCStatus = "VERIFIED"
	KYCRejected KYCStatus = "REJECTED"
)

type Company struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	LegalName   string     `gorm:"type:varchar(255);not null" json:"legal_name"`
	EntityType  EntityType `gorm:"type:varchar(50);not null" json:"entity_type"`
	NPWP        string     `gorm:"type:varchar(50);unique;not null" json:"npwp"`
	NIB         string     `gorm:"type:varchar(50)" json:"nib"`
	Address     string     `gorm:"type:text" json:"address"`
	KYCStatus   KYCStatus  `gorm:"type:varchar(50);default:'PENDING'" json:"kyc_status"`
	CreditLimit int64      `gorm:"type:bigint;default:0" json:"credit_limit"` // in Cents/IDR (no float)
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Users    []User    `gorm:"foreignKey:CompanyID" json:"users,omitempty"`
	Accounts []Account `gorm:"foreignKey:CompanyID" json:"accounts,omitempty"`
}

type UserRole string

const (
	RoleAdmin           UserRole = "ADMIN"
	RolePabrikManager   UserRole = "PABRIK_MANAGER"
	RoleLapakOwner      UserRole = "LAPAK_OWNER"
	RoleWeighOperator   UserRole = "WEIGH_OPERATOR"
	RoleFintechOfficer  UserRole = "FINTECH_OFFICER"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CompanyID    uuid.UUID `gorm:"type:uuid;not null;index" json:"company_id"`
	FullName     string    `gorm:"type:varchar(255);not null" json:"full_name"`
	Email        string    `gorm:"type:varchar(255);unique;not null;index" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         UserRole  `gorm:"type:varchar(50);not null" json:"role"`
	PhoneNumber  string    `gorm:"type:varchar(50)" json:"phone_number"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}
