package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentPaid      PaymentStatus = "PAID"
	PaymentExpired   PaymentStatus = "EXPIRED"
	PaymentFailed    PaymentStatus = "FAILED"
	PaymentDisbursed PaymentStatus = "DISBURSED"
)

type PaymentChannel string

const (
	ChannelBCAVA     PaymentChannel = "BCA_VA"
	ChannelMandiriVA PaymentChannel = "MANDIRI_VA"
	ChannelBRIVA     PaymentChannel = "BRI_VA"
	ChannelQRIS      PaymentChannel = "QRIS"
	ChannelBankPayout PaymentChannel = "BANK_DISBURSEMENT"
)

type PaymentOrder struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CompanyID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"company_id"`
	ReferenceID    string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"reference_id"` // E.g. INV-VA-xxx, TND-xxx
	PaymentType    string         `gorm:"type:varchar(50);not null" json:"payment_type"`             // BID_BOND, ESCROW_DEPOSIT, TICKET_PAYOUT
	Channel        PaymentChannel `gorm:"type:varchar(50);not null" json:"channel"`
	Amount         int64          `gorm:"type:bigint;not null" json:"amount"`
	VANumber       string         `gorm:"type:varchar(100)" json:"va_number"`
	Status         PaymentStatus  `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	ExternalID     string         `gorm:"type:varchar(255)" json:"external_id"`
	SignatureHash  string         `gorm:"type:varchar(255)" json:"signature_hash"`
	PaidAt         *time.Time     `json:"paid_at,omitempty"`
	ReconciledAt   *time.Time     `json:"reconciled_at,omitempty"`
	ExpiresAt      time.Time      `gorm:"not null" json:"expires_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`

	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

type ReconciliationReport struct {
	BatchID          string    `json:"batch_id"`
	AuditedCount     int       `json:"audited_count"`
	MatchedCount     int       `json:"matched_count"`
	DiscrepancyCount int       `json:"discrepancy_count"`
	TotalAmountIDR   int64     `json:"total_amount_idr"`
	ExecutedAt       time.Time `json:"executed_at"`
}
