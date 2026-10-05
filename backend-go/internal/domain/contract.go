package domain

import (
	"time"

	"github.com/google/uuid"
)

type ContractStatus string

const (
	ContractPendingEscrow ContractStatus = "PENDING_ESCROW"
	ContractActive        ContractStatus = "ACTIVE"
	ContractCompleted     ContractStatus = "COMPLETED"
	ContractDisputed      ContractStatus = "DISPUTED"
)

type Contract struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TenderID             uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"tender_id"`
	SellerCompanyID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"seller_company_id"` // Pabrik
	BuyerCompanyID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"buyer_company_id"`  // Lapak Pemenang
	ContractNumber       string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"contract_number"`
	AgreedPricePerKG     int64          `gorm:"type:bigint;not null" json:"agreed_price_per_kg"`
	EstimatedTotalAmount int64          `gorm:"type:bigint;not null" json:"estimated_total_amount"`
	EscrowLockedAmount   int64          `gorm:"type:bigint;default:0" json:"escrow_locked_amount"`
	TotalSettledAmount   int64          `gorm:"type:bigint;default:0" json:"total_settled_amount"`
	Status               ContractStatus `gorm:"type:varchar(50);default:'PENDING_ESCROW'" json:"status"`
	StartDate            time.Time      `json:"start_date"`
	EndDate              time.Time      `json:"end_date"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`

	Tender             *Tender              `gorm:"foreignKey:TenderID" json:"tender,omitempty"`
	SellerCompany      *Company             `gorm:"foreignKey:SellerCompanyID" json:"seller_company,omitempty"`
	BuyerCompany       *Company             `gorm:"foreignKey:BuyerCompanyID" json:"buyer_company,omitempty"`
	WeighbridgeTickets []WeighbridgeTicket  `gorm:"foreignKey:ContractID" json:"weighbridge_tickets,omitempty"`
}

type TicketStatus string

const (
	TicketPendingReview TicketStatus = "PENDING_REVIEW"
	TicketVerified      TicketStatus = "VERIFIED"
	TicketSettled       TicketStatus = "SETTLED"
	TicketDisputed      TicketStatus = "DISPUTED"
)

type WeighbridgeTicket struct {
	ID                     uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	ContractID             uuid.UUID    `gorm:"type:uuid;not null;index" json:"contract_id"`
	TicketNumber           string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"ticket_number"`
	TruckLicensePlate      string       `gorm:"type:varchar(50);not null;index" json:"truck_license_plate"`
	DriverName             string       `gorm:"type:varchar(255)" json:"driver_name"`
	GrossWeightKG          float64      `gorm:"type:numeric(14,4);not null" json:"gross_weight_kg"` // Truk + Muatan
	TareWeightKG           float64      `gorm:"type:numeric(14,4);not null" json:"tare_weight_kg"`  // Truk Kosong
	NetWeightKG            float64      `gorm:"type:numeric(14,4);not null" json:"net_weight_kg"`   // Gross - Tare
	RefractionPercentage   float64      `gorm:"type:numeric(6,2);default:0.0" json:"refraction_percentage"` // Potongan karat/kotoran (%)
	RefractionKG           float64      `gorm:"type:numeric(14,4);default:0.0" json:"refraction_kg"`
	FinalBillableWeightKG  float64      `gorm:"type:numeric(14,4);not null" json:"final_billable_weight_kg"`
	PricePerKG             int64        `gorm:"type:bigint;not null" json:"price_per_kg"`
	GrossPayableAmount     int64        `gorm:"type:bigint;not null" json:"gross_payable_amount"`
	PlatformFeeAmount      int64        `gorm:"type:bigint;default:0" json:"platform_fee_amount"`
	NetToSellerAmount      int64        `gorm:"type:bigint;not null" json:"net_to_seller_amount"`
	SlipPhotoURL           string       `gorm:"type:text" json:"slip_photo_url"`
	Status                 TicketStatus `gorm:"type:varchar(50);default:'PENDING_REVIEW'" json:"status"`
	WeighedAt              time.Time    `gorm:"not null" json:"weighed_at"`
	SettledAt              *time.Time   `json:"settled_at,omitempty"`
	CreatedAt              time.Time    `json:"created_at"`
	UpdatedAt              time.Time    `json:"updated_at"`

	Contract *Contract `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
}

// CalculateWeights computes the net weight, refraction deduction, and billable weight
func (w *WeighbridgeTicket) CalculateWeights(pricePerKG int64, feePercentage float64) {
	w.NetWeightKG = w.GrossWeightKG - w.TareWeightKG
	if w.NetWeightKG < 0 {
		w.NetWeightKG = 0
	}

	w.RefractionKG = (w.NetWeightKG * w.RefractionPercentage) / 100.0
	w.FinalBillableWeightKG = w.NetWeightKG - w.RefractionKG
	if w.FinalBillableWeightKG < 0 {
		w.FinalBillableWeightKG = 0
	}

	w.PricePerKG = pricePerKG
	w.GrossPayableAmount = int64(w.FinalBillableWeightKG * float64(pricePerKG))
	w.PlatformFeeAmount = int64(float64(w.GrossPayableAmount) * (feePercentage / 100.0))
	w.NetToSellerAmount = w.GrossPayableAmount - w.PlatformFeeAmount
}
