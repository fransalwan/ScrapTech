package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"scrapflow-backend/internal/domain"
)

type WeighbridgeRepository interface {
	CreateTicket(db *gorm.DB, ticket *domain.WeighbridgeTicket) error
	GetTicketByID(db *gorm.DB, id uuid.UUID) (*domain.WeighbridgeTicket, error)
	GetTicketByIDForUpdate(db *gorm.DB, id uuid.UUID) (*domain.WeighbridgeTicket, error)
	GetTicketsByContract(db *gorm.DB, contractID uuid.UUID) ([]domain.WeighbridgeTicket, error)
	UpdateTicket(db *gorm.DB, ticket *domain.WeighbridgeTicket) error

	CreateFinancing(db *gorm.DB, f *domain.FinancingFacility) error
	GetFinancingByID(db *gorm.DB, id uuid.UUID) (*domain.FinancingFacility, error)
	GetFinancingByBorrower(db *gorm.DB, borrowerID uuid.UUID) ([]domain.FinancingFacility, error)
	UpdateFinancing(db *gorm.DB, f *domain.FinancingFacility) error
}

type weighbridgeRepository struct{}

func NewWeighbridgeRepository() WeighbridgeRepository {
	return &weighbridgeRepository{}
}

func (r *weighbridgeRepository) CreateTicket(db *gorm.DB, ticket *domain.WeighbridgeTicket) error {
	return db.Create(ticket).Error
}

func (r *weighbridgeRepository) GetTicketByID(db *gorm.DB, id uuid.UUID) (*domain.WeighbridgeTicket, error) {
	var ticket domain.WeighbridgeTicket
	err := db.Preload("Contract.SellerCompany").
		Preload("Contract.BuyerCompany").
		First(&ticket, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *weighbridgeRepository) GetTicketByIDForUpdate(db *gorm.DB, id uuid.UUID) (*domain.WeighbridgeTicket, error) {
	var ticket domain.WeighbridgeTicket
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Contract").
		First(&ticket, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *weighbridgeRepository) GetTicketsByContract(db *gorm.DB, contractID uuid.UUID) ([]domain.WeighbridgeTicket, error) {
	var tickets []domain.WeighbridgeTicket
	err := db.Where("contract_id = ?", contractID).
		Order("weighed_at desc").
		Find(&tickets).Error
	return tickets, err
}

func (r *weighbridgeRepository) UpdateTicket(db *gorm.DB, ticket *domain.WeighbridgeTicket) error {
	return db.Save(ticket).Error
}

func (r *weighbridgeRepository) CreateFinancing(db *gorm.DB, f *domain.FinancingFacility) error {
	return db.Create(f).Error
}

func (r *weighbridgeRepository) GetFinancingByID(db *gorm.DB, id uuid.UUID) (*domain.FinancingFacility, error) {
	var f domain.FinancingFacility
	err := db.Preload("BorrowerCompany").
		Preload("Contract").
		First(&f, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *weighbridgeRepository) GetFinancingByBorrower(db *gorm.DB, borrowerID uuid.UUID) ([]domain.FinancingFacility, error) {
	var list []domain.FinancingFacility
	err := db.Preload("Contract").
		Where("borrower_company_id = ?", borrowerID).
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *weighbridgeRepository) UpdateFinancing(db *gorm.DB, f *domain.FinancingFacility) error {
	return db.Save(f).Error
}
