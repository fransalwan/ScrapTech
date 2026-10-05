package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"scrapflow-backend/internal/domain"
)

type TenderRepository interface {
	CreateTender(db *gorm.DB, tender *domain.Tender) error
	GetTenderByID(db *gorm.DB, id uuid.UUID) (*domain.Tender, error)
	GetTenderByIDForUpdate(db *gorm.DB, id uuid.UUID) (*domain.Tender, error)
	GetTenders(db *gorm.DB, status string) ([]domain.Tender, error)
	UpdateTender(db *gorm.DB, tender *domain.Tender) error

	CreateBid(db *gorm.DB, bid *domain.Bid) error
	GetBidsByTenderID(db *gorm.DB, tenderID uuid.UUID) ([]domain.Bid, error)
	GetBidByID(db *gorm.DB, id uuid.UUID) (*domain.Bid, error)
	UpdateBid(db *gorm.DB, bid *domain.Bid) error

	CreateContract(db *gorm.DB, contract *domain.Contract) error
	GetContractByID(db *gorm.DB, id uuid.UUID) (*domain.Contract, error)
	GetContractByIDForUpdate(db *gorm.DB, id uuid.UUID) (*domain.Contract, error)
	GetContractsByCompany(db *gorm.DB, companyID uuid.UUID) ([]domain.Contract, error)
	UpdateContract(db *gorm.DB, contract *domain.Contract) error
}

type tenderRepository struct{}

func NewTenderRepository() TenderRepository {
	return &tenderRepository{}
}

func (r *tenderRepository) CreateTender(db *gorm.DB, tender *domain.Tender) error {
	return db.Create(tender).Error
}

func (r *tenderRepository) GetTenderByID(db *gorm.DB, id uuid.UUID) (*domain.Tender, error) {
	var tender domain.Tender
	err := db.Preload("PublisherCompany").
		Preload("Bids.BidderCompany").
		Preload("Contract").
		First(&tender, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tender, nil
}

// GetTenderByIDForUpdate uses PostgreSQL pessimistic row-locking
func (r *tenderRepository) GetTenderByIDForUpdate(db *gorm.DB, id uuid.UUID) (*domain.Tender, error) {
	var tender domain.Tender
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("PublisherCompany").
		Preload("Bids").
		First(&tender, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tender, nil
}

func (r *tenderRepository) GetTenders(db *gorm.DB, status string) ([]domain.Tender, error) {
	var list []domain.Tender
	q := db.Preload("PublisherCompany").Preload("Bids").Order("created_at desc")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *tenderRepository) UpdateTender(db *gorm.DB, tender *domain.Tender) error {
	return db.Save(tender).Error
}

func (r *tenderRepository) CreateBid(db *gorm.DB, bid *domain.Bid) error {
	return db.Create(bid).Error
}

func (r *tenderRepository) GetBidsByTenderID(db *gorm.DB, tenderID uuid.UUID) ([]domain.Bid, error) {
	var bids []domain.Bid
	err := db.Preload("BidderCompany").
		Where("tender_id = ?", tenderID).
		Order("price_per_kg desc").
		Find(&bids).Error
	return bids, err
}

func (r *tenderRepository) GetBidByID(db *gorm.DB, id uuid.UUID) (*domain.Bid, error) {
	var bid domain.Bid
	err := db.Preload("Tender").Preload("BidderCompany").First(&bid, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

func (r *tenderRepository) UpdateBid(db *gorm.DB, bid *domain.Bid) error {
	return db.Save(bid).Error
}

func (r *tenderRepository) CreateContract(db *gorm.DB, contract *domain.Contract) error {
	return db.Create(contract).Error
}

func (r *tenderRepository) GetContractByID(db *gorm.DB, id uuid.UUID) (*domain.Contract, error) {
	var contract domain.Contract
	err := db.Preload("Tender").
		Preload("SellerCompany").
		Preload("BuyerCompany").
		Preload("WeighbridgeTickets").
		First(&contract, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

// GetContractByIDForUpdate uses PostgreSQL pessimistic row-locking for atomic weighbridge settlement
func (r *tenderRepository) GetContractByIDForUpdate(db *gorm.DB, id uuid.UUID) (*domain.Contract, error) {
	var contract domain.Contract
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Tender").
		Preload("SellerCompany").
		Preload("BuyerCompany").
		First(&contract, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *tenderRepository) GetContractsByCompany(db *gorm.DB, companyID uuid.UUID) ([]domain.Contract, error) {
	var contracts []domain.Contract
	err := db.Preload("Tender").
		Preload("SellerCompany").
		Preload("BuyerCompany").
		Where("seller_company_id = ? OR buyer_company_id = ?", companyID, companyID).
		Order("created_at desc").
		Find(&contracts).Error
	return contracts, err
}

func (r *tenderRepository) UpdateContract(db *gorm.DB, contract *domain.Contract) error {
	return db.Save(contract).Error
}
