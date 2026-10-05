package domain

import (
	"time"

	"github.com/google/uuid"
)

type ScrapCategory string

const (
	ScrapHMS1     ScrapCategory = "HMS_1"     // Heavy Melting Steel Grade 1
	ScrapHMS2     ScrapCategory = "HMS_2"     // Heavy Melting Steel Grade 2
	ScrapPlate    ScrapCategory = "PLATE"     // Potongan plat kapal/konstruksi
	ScrapCastIron ScrapCategory = "CAST_IRON" // Besi cor / blok mesin
	ScrapCopper   ScrapCategory = "COPPER"    // Tembaga / kabel kupas
)

type TenderStatus string

const (
	TenderDraft     TenderStatus = "DRAFT"
	TenderPublished TenderStatus = "PUBLISHED"
	TenderBidding   TenderStatus = "BIDDING"
	TenderAwarded   TenderStatus = "AWARDED"
	TenderExecuting TenderStatus = "EXECUTING"
	TenderCompleted TenderStatus = "COMPLETED"
	TenderCancelled TenderStatus = "CANCELLED"
)

type Tender struct {
	ID                 uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	PublisherCompanyID uuid.UUID     `gorm:"type:uuid;not null;index" json:"publisher_company_id"`
	TenderNumber       string        `gorm:"type:varchar(100);uniqueIndex;not null" json:"tender_number"`
	Title              string        `gorm:"type:varchar(255);not null" json:"title"`
	ScrapType          ScrapCategory `gorm:"type:varchar(50);not null" json:"scrap_type"`
	EstimatedWeightKG  float64       `gorm:"type:numeric(14,4);not null" json:"estimated_weight_kg"`
	ReservePricePerKG  int64         `gorm:"type:bigint;not null" json:"reserve_price_per_kg"`
	BidBondAmount      int64         `gorm:"type:bigint;not null" json:"bid_bond_amount"` // Jaminan lelang
	Location           string        `gorm:"type:text;not null" json:"location"`
	BiddingDeadline    time.Time     `gorm:"not null" json:"bidding_deadline"`
	Status             TenderStatus  `gorm:"type:varchar(50);default:'PUBLISHED'" json:"status"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`

	PublisherCompany *Company  `gorm:"foreignKey:PublisherCompanyID" json:"publisher_company,omitempty"`
	Bids             []Bid     `gorm:"foreignKey:TenderID" json:"bids,omitempty"`
	Contract         *Contract `gorm:"foreignKey:TenderID" json:"contract,omitempty"`
}

type BidStatus string

const (
	BidSubmitted BidStatus = "SUBMITTED"
	BidWon       BidStatus = "WON"
	BidLost      BidStatus = "LOST"
	BidCancelled BidStatus = "CANCELLED"
)

type Bid struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenderID         uuid.UUID `gorm:"type:uuid;not null;index" json:"tender_id"`
	BidderCompanyID  uuid.UUID `gorm:"type:uuid;not null;index" json:"bidder_company_id"`
	PricePerKG       int64     `gorm:"type:bigint;not null" json:"price_per_kg"`
	TotalOfferAmount int64     `gorm:"type:bigint;not null" json:"total_offer_amount"`
	BidBondLocked    bool      `gorm:"default:false" json:"bid_bond_locked"`
	Status           BidStatus `gorm:"type:varchar(50);default:'SUBMITTED'" json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Tender        *Tender  `gorm:"foreignKey:TenderID" json:"tender,omitempty"`
	BidderCompany *Company `gorm:"foreignKey:BidderCompanyID" json:"bidder_company,omitempty"`
}
