package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/domain"
	"scrapflow-backend/internal/repository"
)

type FinancingUsecase interface {
	ApplyFacility(db *gorm.DB, companyID uuid.UUID, contractID uuid.UUID, requestedAmount int64, tenorDays int) (*domain.FinancingFacility, error)
	GetFacilitiesByBorrower(db *gorm.DB, borrowerID uuid.UUID) ([]domain.FinancingFacility, error)
	DisburseFacility(db *gorm.DB, idempotencyKey string, facilityID uuid.UUID) (*domain.FinancingFacility, error)
}

type financingUsecase struct {
	weighRepo   repository.WeighbridgeRepository
	tenderRepo  repository.TenderRepository
	ledgerRepo  repository.LedgerRepository
	companyRepo repository.CompanyRepository
}

func NewFinancingUsecase(
	wRepo repository.WeighbridgeRepository,
	tRepo repository.TenderRepository,
	lRepo repository.LedgerRepository,
	cRepo repository.CompanyRepository,
) FinancingUsecase {
	return &financingUsecase{
		weighRepo:   wRepo,
		tenderRepo:  tRepo,
		ledgerRepo:  lRepo,
		companyRepo: cRepo,
	}
}

func (u *financingUsecase) ApplyFacility(db *gorm.DB, companyID uuid.UUID, contractID uuid.UUID, requestedAmount int64, tenorDays int) (*domain.FinancingFacility, error) {
	contract, err := u.tenderRepo.GetContractByID(db, contractID)
	if err != nil {
		return nil, errors.New("contract not found")
	}

	if contract.BuyerCompanyID != companyID {
		return nil, errors.New("only the winning contract buyer can apply for financing on this contract")
	}

	company, err := u.companyRepo.GetCompanyByID(db, companyID)
	if err != nil {
		return nil, errors.New("company not found")
	}

	if requestedAmount > company.CreditLimit {
		return nil, fmt.Errorf("requested amount (%d) exceeds company approved credit limit (%d)", requestedAmount, company.CreditLimit)
	}

	now := time.Now()
	facility := domain.FinancingFacility{
		ID:                    uuid.New(),
		BorrowerCompanyID:     companyID,
		ContractID:            contractID,
		FacilityNumber:        fmt.Sprintf("SCF-%s-%d", now.Format("20060102"), now.UnixNano()%10000),
		RequestedAmount:       requestedAmount,
		ApprovedAmount:        requestedAmount,
		InterestRateAnnualPct: 12.0,
		TenorDays:             tenorDays,
		RepaymentDueDate:      now.Add(time.Duration(tenorDays) * 24 * time.Hour),
		Status:                domain.FinancingPendingReview,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := u.weighRepo.CreateFinancing(db, &facility); err != nil {
		return nil, err
	}
	return &facility, nil
}

func (u *financingUsecase) GetFacilitiesByBorrower(db *gorm.DB, borrowerID uuid.UUID) ([]domain.FinancingFacility, error) {
	return u.weighRepo.GetFinancingByBorrower(db, borrowerID)
}

func (u *financingUsecase) DisburseFacility(db *gorm.DB, idempotencyKey string, facilityID uuid.UUID) (*domain.FinancingFacility, error) {
	var facility *domain.FinancingFacility

	err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		facility, err = u.weighRepo.GetFinancingByID(tx, facilityID)
		if err != nil {
			return errors.New("financing facility not found")
		}

		if facility.Status != domain.FinancingPendingReview && facility.Status != domain.FinancingApproved {
			return errors.New("facility is not eligible for disbursement")
		}

		// Calculate 1% upfront origination fee
		originationFee := int64(float64(facility.ApprovedAmount) * 0.01)
		netEscrowInjected := facility.ApprovedAmount - originationFee

		// Double-Entry Ledger Posting:
		// [DEBIT]  1020 (SCF Financing Loan Receivables) : facility.ApprovedAmount
		// [CREDIT] 2020 (Escrow Milestone Payable)       : netEscrowInjected
		// [CREDIT] 4020 (SCF Financing Fee Revenue)      : originationFee
		loanReceivableAcc, err := u.ledgerRepo.GetAccountByCode(tx, "1020")
		if err != nil {
			return err
		}
		escrowAcc, err := u.ledgerRepo.GetAccountByCode(tx, "2020")
		if err != nil {
			return err
		}
		revenueAcc, err := u.ledgerRepo.GetAccountByCode(tx, "4020")
		if err != nil {
			return err
		}

		entries := []domain.LedgerEntry{
			{AccountID: loanReceivableAcc.ID, EntryType: domain.EntryDebit, Amount: facility.ApprovedAmount},
			{AccountID: escrowAcc.ID, EntryType: domain.EntryCredit, Amount: netEscrowInjected},
			{AccountID: revenueAcc.ID, EntryType: domain.EntryCredit, Amount: originationFee},
		}

		_, err = u.ledgerRepo.PostJournal(
			tx,
			idempotencyKey,
			domain.TxTypeSCFLoanDisburse,
			facility.FacilityNumber,
			fmt.Sprintf("Disbursement for facility %s into tender escrow", facility.FacilityNumber),
			entries,
		)
		if err != nil {
			return fmt.Errorf("disbursement ledger entry failed: %w", err)
		}

		facility.Status = domain.FinancingDisbursed
		return u.weighRepo.UpdateFinancing(tx, facility)
	})

	if err != nil {
		return nil, err
	}
	return facility, nil
}
