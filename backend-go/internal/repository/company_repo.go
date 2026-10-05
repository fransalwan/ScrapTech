package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/domain"
)

type CompanyRepository interface {
	CreateCompany(db *gorm.DB, company *domain.Company) error
	GetCompanyByID(db *gorm.DB, id uuid.UUID) (*domain.Company, error)
	GetAllCompanies(db *gorm.DB) ([]domain.Company, error)
	CreateUser(db *gorm.DB, user *domain.User) error
	GetUserByEmail(db *gorm.DB, email string) (*domain.User, error)
	GetUserByID(db *gorm.DB, id uuid.UUID) (*domain.User, error)
}

type companyRepository struct{}

func NewCompanyRepository() CompanyRepository {
	return &companyRepository{}
}

func (r *companyRepository) CreateCompany(db *gorm.DB, company *domain.Company) error {
	return db.Create(company).Error
}

func (r *companyRepository) GetCompanyByID(db *gorm.DB, id uuid.UUID) (*domain.Company, error) {
	var company domain.Company
	err := db.Preload("Users").Preload("Accounts").First(&company, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) GetAllCompanies(db *gorm.DB) ([]domain.Company, error) {
	var list []domain.Company
	err := db.Find(&list).Error
	return list, err
}

func (r *companyRepository) CreateUser(db *gorm.DB, user *domain.User) error {
	return db.Create(user).Error
}

func (r *companyRepository) GetUserByEmail(db *gorm.DB, email string) (*domain.User, error) {
	var user domain.User
	err := db.Preload("Company").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *companyRepository) GetUserByID(db *gorm.DB, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := db.Preload("Company").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
