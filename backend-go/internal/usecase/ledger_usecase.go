package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/domain"
	"scrapflow-backend/internal/repository"
)

type LedgerUsecase interface {
	GetCompanyWalletBalance(db *gorm.DB, companyID uuid.UUID) (int64, error)
	GetLedgerTransactions(db *gorm.DB, limit, offset int) ([]domain.Transaction, error)
	DepositEscrowVault(db *gorm.DB, idempotencyKey string, companyID uuid.UUID, amount int64) (*domain.Transaction, error)
}

type ledgerUsecase struct {
	ledgerRepo repository.LedgerRepository
}

func NewLedgerUsecase(lRepo repository.LedgerRepository) LedgerUsecase {
	return &ledgerUsecase{ledgerRepo: lRepo}
}

func (u *ledgerUsecase) GetCompanyWalletBalance(db *gorm.DB, companyID uuid.UUID) (int64, error) {
	acc, err := u.ledgerRepo.GetOrCreateCompanyAccount(
		db,
		companyID,
		domain.AccountLiability,
		"WALLET-"+companyID.String()[:8],
		"Withdrawable Wallet",
	)
	if err != nil {
		return 0, err
	}
	return u.ledgerRepo.GetAccountBalance(db, acc.ID)
}

func (u *ledgerUsecase) GetLedgerTransactions(db *gorm.DB, limit, offset int) ([]domain.Transaction, error) {
	if limit <= 0 {
		limit = 20
	}
	return u.ledgerRepo.GetLedgerTransactions(db, limit, offset)
}

func (u *ledgerUsecase) DepositEscrowVault(db *gorm.DB, idempotencyKey string, companyID uuid.UUID, amount int64) (*domain.Transaction, error) {
	var trans *domain.Transaction

	err := db.Transaction(func(tx *gorm.DB) error {
		// Debit 1010 (Bank platform custody)
		bankAcc, err := u.ledgerRepo.GetAccountByCode(tx, "1010")
		if err != nil {
			return err
		}

		// Credit 2020 (Escrow Milestone Payable)
		escrowAcc, err := u.ledgerRepo.GetAccountByCode(tx, "2020")
		if err != nil {
			return err
		}

		entries := []domain.LedgerEntry{
			{AccountID: bankAcc.ID, EntryType: domain.EntryDebit, Amount: amount},
			{AccountID: escrowAcc.ID, EntryType: domain.EntryCredit, Amount: amount},
		}

		trans, err = u.ledgerRepo.PostJournal(
			tx,
			idempotencyKey,
			domain.TxTypeEscrowLock,
			companyID.String(),
			fmt.Sprintf("Escrow deposit funded by company %s", companyID.String()[:8]),
			entries,
		)
		return err
	})

	return trans, err
}
