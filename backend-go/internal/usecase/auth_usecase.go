package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/domain"
	"scrapflow-backend/internal/repository"
	"scrapflow-backend/pkg/utils"
)

type AuthUsecase interface {
	RegisterCompanyAndUser(db *gorm.DB, req *RegisterRequest) (*RegisterResponse, error)
	Login(db *gorm.DB, email, password, jwtSecret string) (*LoginResponse, error)
}

type RegisterRequest struct {
	LegalName   string            `json:"legal_name" binding:"required"`
	EntityType  domain.EntityType `json:"entity_type" binding:"required"`
	NPWP        string            `json:"npwp" binding:"required"`
	NIB         string            `json:"nib"`
	Address     string            `json:"address"`
	FullName    string            `json:"full_name" binding:"required"`
	Email       string            `json:"email" binding:"required,email"`
	Password    string            `json:"password" binding:"required,min=6"`
	Role        domain.UserRole   `json:"role" binding:"required"`
	PhoneNumber string            `json:"phone_number"`
}

type RegisterResponse struct {
	CompanyID uuid.UUID `json:"company_id"`
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	LegalName string    `json:"legal_name"`
}

type LoginResponse struct {
	Token     string          `json:"token"`
	UserID    uuid.UUID       `json:"user_id"`
	Email     string          `json:"email"`
	FullName  string          `json:"full_name"`
	Role      domain.UserRole `json:"role"`
	CompanyID uuid.UUID       `json:"company_id"`
	LegalName string          `json:"legal_name"`
}

type authUsecase struct {
	companyRepo repository.CompanyRepository
	ledgerRepo  repository.LedgerRepository
}

func NewAuthUsecase(cRepo repository.CompanyRepository, lRepo repository.LedgerRepository) AuthUsecase {
	return &authUsecase{companyRepo: cRepo, ledgerRepo: lRepo}
}

func (u *authUsecase) RegisterCompanyAndUser(db *gorm.DB, req *RegisterRequest) (*RegisterResponse, error) {
	// Check if user email already exists
	if _, err := u.companyRepo.GetUserByEmail(db, req.Email); err == nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	companyID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	// Transaction to create Company, User, and Company Ledger Accounts
	err = db.Transaction(func(tx *gorm.DB) error {
		company := domain.Company{
			ID:          companyID,
			LegalName:   req.LegalName,
			EntityType:  req.EntityType,
			NPWP:        req.NPWP,
			NIB:         req.NIB,
			Address:     req.Address,
			KYCStatus:   domain.KYCVerified, // auto-verify for dev/sandbox
			CreditLimit: 1000000000,         // 1 Miliar default credit limit
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if err := u.companyRepo.CreateCompany(tx, &company); err != nil {
			return err
		}

		user := domain.User{
			ID:           userID,
			CompanyID:    companyID,
			FullName:     req.FullName,
			Email:        req.Email,
			PasswordHash: hashedPassword,
			Role:         req.Role,
			PhoneNumber:  req.PhoneNumber,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if err := u.companyRepo.CreateUser(tx, &user); err != nil {
			return err
		}

		// Create dedicated Ledger Accounts for this Company
		// Code: 2030-{CompanyID} (Withdrawable wallet balance)
		_, err := u.ledgerRepo.GetOrCreateCompanyAccount(tx, companyID, domain.AccountLiability, "WALLET-"+companyID.String()[:8], req.LegalName+" Withdrawable Wallet")
		return err
	})

	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		CompanyID: companyID,
		UserID:    userID,
		Email:     req.Email,
		LegalName: req.LegalName,
	}, nil
}

func (u *authUsecase) Login(db *gorm.DB, email, password, jwtSecret string) (*LoginResponse, error) {
	user, err := u.companyRepo.GetUserByEmail(db, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role), user.CompanyID, jwtSecret, 72*time.Hour)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:     token,
		UserID:    user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		CompanyID: user.CompanyID,
		LegalName: user.Company.LegalName,
	}, nil
}
