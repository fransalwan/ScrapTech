package domain

import (
	"time"

	"github.com/google/uuid"
)

type FinancingStatus string

const (
	FinancingPendingReview FinancingStatus = "PENDING_REVIEW"
	FinancingApproved      FinancingStatus = "APPROVED"
	FinancingDisbursed     FinancingStatus = "DISBURSED"
	FinancingRepaid        FinancingStatus = "REPAID"
	FinancingDefaulted     FinancingStatus = "DEFAULTED"
)

type FinancingFacility struct {
	ID                    uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	BorrowerCompanyID     uuid.UUID       `gorm:"type:uuid;not null;index" json:"borrower_company_id"` // Lapak
	ContractID            uuid.UUID       `gorm:"type:uuid;not null;index" json:"contract_id"`         // Dasar SPK Tender
	FunderCompanyID       *uuid.UUID      `gorm:"type:uuid;index" json:"funder_company_id,omitempty"`  // P2P/Bank Funder
	FacilityNumber        string          `gorm:"type:varchar(100);uniqueIndex;not null" json:"facility_number"`
	RequestedAmount       int64           `gorm:"type:bigint;not null" json:"requested_amount"`
	ApprovedAmount        int64           `gorm:"type:bigint;default:0" json:"approved_amount"`
	InterestRateAnnualPct float64         `gorm:"type:numeric(5,2);default:12.0" json:"interest_rate_annual_pct"`
	TenorDays             int             `gorm:"not null" json:"tenor_days"` // e.g. 14, 30, 45 days
	RepaymentDueDate      time.Time       `gorm:"not null" json:"repayment_due_date"`
	Status                FinancingStatus `gorm:"type:varchar(50);default:'PENDING_REVIEW'" json:"status"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`

	BorrowerCompany *Company  `gorm:"foreignKey:BorrowerCompanyID" json:"borrower_company,omitempty"`
	Contract        *Contract `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
}
