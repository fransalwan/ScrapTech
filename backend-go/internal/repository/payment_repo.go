package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"scrapflow-backend/internal/domain"
)

type PaymentRepository interface {
	CreatePaymentOrder(db *gorm.DB, p *domain.PaymentOrder) error
	GetPaymentOrderByReferenceID(db *gorm.DB, refID string) (*domain.PaymentOrder, error)
	GetPaymentOrderByReferenceIDForUpdate(db *gorm.DB, refID string) (*domain.PaymentOrder, error)
	UpdatePaymentOrder(db *gorm.DB, p *domain.PaymentOrder) error
	GetPendingPaymentOrders(db *gorm.DB, olderThan time.Time, limit int) ([]domain.PaymentOrder, error)
	GetPaymentOrdersByCompany(db *gorm.DB, companyID uuid.UUID) ([]domain.PaymentOrder, error)
}

type paymentRepository struct{}

func NewPaymentRepository() PaymentRepository {
	return &paymentRepository{}
}

func (r *paymentRepository) CreatePaymentOrder(db *gorm.DB, p *domain.PaymentOrder) error {
	return db.Create(p).Error
}

func (r *paymentRepository) GetPaymentOrderByReferenceID(db *gorm.DB, refID string) (*domain.PaymentOrder, error) {
	var p domain.PaymentOrder
	err := db.Preload("Company").First(&p, "reference_id = ?", refID).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) GetPaymentOrderByReferenceIDForUpdate(db *gorm.DB, refID string) (*domain.PaymentOrder, error) {
	var p domain.PaymentOrder
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Company").
		First(&p, "reference_id = ?", refID).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) UpdatePaymentOrder(db *gorm.DB, p *domain.PaymentOrder) error {
	return db.Save(p).Error
}

func (r *paymentRepository) GetPendingPaymentOrders(db *gorm.DB, olderThan time.Time, limit int) ([]domain.PaymentOrder, error) {
	var list []domain.PaymentOrder
	err := db.Where("status = ? AND created_at <= ?", domain.PaymentPending, olderThan).
		Limit(limit).
		Find(&list).Error
	return list, err
}

func (r *paymentRepository) GetPaymentOrdersByCompany(db *gorm.DB, companyID uuid.UUID) ([]domain.PaymentOrder, error) {
	var list []domain.PaymentOrder
	err := db.Where("company_id = ?", companyID).Order("created_at desc").Find(&list).Error
	return list, err
}
