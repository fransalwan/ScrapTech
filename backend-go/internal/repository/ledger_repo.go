package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/domain"
)

type LedgerRepository interface {
	PostJournal(tx *gorm.DB, idempotencyKey string, txType domain.TransactionType, refID, desc string, entries []domain.LedgerEntry) (*domain.Transaction, error)
	GetAccountBalance(db *gorm.DB, accountID uuid.UUID) (int64, error)
	GetAccountByCode(db *gorm.DB, code string) (*domain.Account, error)
	GetOrCreateCompanyAccount(db *gorm.DB, companyID uuid.UUID, accountType domain.AccountType, code, name string) (*domain.Account, error)
	GetLedgerTransactions(db *gorm.DB, limit, offset int) ([]domain.Transaction, error)
}

type ledgerRepository struct{}

func NewLedgerRepository() LedgerRepository {
	return &ledgerRepository{}
}

func (r *ledgerRepository) PostJournal(tx *gorm.DB, idempotencyKey string, txType domain.TransactionType, refID, desc string, entries []domain.LedgerEntry) (*domain.Transaction, error) {
	if tx == nil {
		return nil, errors.New("ledger posting requires an active database transaction")
	}

	// 1. Verify Zero-Sum Double-Entry Invariant (Sum(Debits) == Sum(Credits))
	if err := domain.ValidateZeroSum(entries); err != nil {
		return nil, fmt.Errorf("double-entry violation: %w", err)
	}

	// 2. Check for duplicate idempotency key in database
	var existing domain.Transaction
	if err := tx.Where("idempotency_key = ?", idempotencyKey).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("transaction with idempotency key '%s' has already been posted", idempotencyKey)
	}

	// 3. Create Transaction Header
	now := time.Now()
	trans := domain.Transaction{
		ID:             uuid.New(),
		IdempotencyKey: idempotencyKey,
		TxType:         txType,
		ReferenceID:    refID,
		Description:    desc,
		PostedAt:       now,
		CreatedAt:      now,
	}

	if err := tx.Create(&trans).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction header: %w", err)
	}

	// 4. Attach & Insert all ledger entries atomically
	for i := range entries {
		entries[i].ID = uuid.New()
		entries[i].TransactionID = trans.ID
		entries[i].CreatedAt = now
		if err := tx.Create(&entries[i]).Error; err != nil {
			return nil, fmt.Errorf("failed to insert ledger entry #%d: %w", i, err)
		}
	}

	trans.Entries = entries
	return &trans, nil
}

func (r *ledgerRepository) GetAccountBalance(db *gorm.DB, accountID uuid.UUID) (int64, error) {
	var account domain.Account
	if err := db.First(&account, "id = ?", accountID).Error; err != nil {
		return 0, err
	}

	// Calculate Debit Sum
	var debitSum struct{ Total int64 }
	db.Model(&domain.LedgerEntry{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("account_id = ? AND entry_type = ?", accountID, domain.EntryDebit).
		Scan(&debitSum)

	// Calculate Credit Sum
	var creditSum struct{ Total int64 }
	db.Model(&domain.LedgerEntry{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("account_id = ? AND entry_type = ?", accountID, domain.EntryCredit).
		Scan(&creditSum)

	// In Accounting:
	// Asset / Expense Normal Balance: Debit - Credit
	// Liability / Equity / Revenue Normal Balance: Credit - Debit
	switch account.AccountType {
	case domain.AccountAsset, domain.AccountExpense:
		return debitSum.Total - creditSum.Total, nil
	default:
		return creditSum.Total - debitSum.Total, nil
	}
}

func (r *ledgerRepository) GetAccountByCode(db *gorm.DB, code string) (*domain.Account, error) {
	var acc domain.Account
	err := db.Where("account_code = ?", code).First(&acc).Error
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *ledgerRepository) GetOrCreateCompanyAccount(db *gorm.DB, companyID uuid.UUID, accountType domain.AccountType, code, name string) (*domain.Account, error) {
	var acc domain.Account
	err := db.Where("company_id = ? AND account_code = ?", companyID, code).First(&acc).Error
	if err == nil {
		return &acc, nil
	}

	acc = domain.Account{
		ID:          uuid.New(),
		CompanyID:   &companyID,
		AccountCode: code,
		AccountName: name,
		AccountType: accountType,
		Currency:    "IDR",
		IsActive:    true,
	}

	if err := db.Create(&acc).Error; err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *ledgerRepository) GetLedgerTransactions(db *gorm.DB, limit, offset int) ([]domain.Transaction, error) {
	var txs []domain.Transaction
	err := db.Preload("Entries.Account").
		Order("posted_at desc").
		Limit(limit).
		Offset(offset).
		Find(&txs).Error
	return txs, err
}
