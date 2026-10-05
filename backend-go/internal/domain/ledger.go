package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountAsset     AccountType = "ASSET"     // 1000s
	AccountLiability AccountType = "LIABILITY" // 2000s
	AccountEquity    AccountType = "EQUITY"    // 3000s
	AccountRevenue   AccountType = "REVENUE"   // 4000s
	AccountExpense   AccountType = "EXPENSE"   // 5000s
)

type EntryType string

const (
	EntryDebit  EntryType = "DEBIT"
	EntryCredit EntryType = "CREDIT"
)

type TransactionType string

const (
	TxTypeBidBondDeposit    TransactionType = "BID_BOND_DEPOSIT"
	TxTypeBidBondRefund     TransactionType = "BID_BOND_REFUND"
	TxTypeEscrowLock        TransactionType = "ESCROW_LOCK"
	TxTypeWeighbridgeSettle TransactionType = "WEIGHBRIDGE_SETTLE"
	TxTypeSCFLoanDisburse   TransactionType = "SCF_LOAN_DISBURSE"
	TxTypeSCFLoanRepay      TransactionType = "SCF_LOAN_REPAY"
	TxTypeUserWithdrawal    TransactionType = "USER_WITHDRAWAL"
)

// Account represents a ledger account (Node in the Chart of Accounts)
type Account struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	CompanyID   *uuid.UUID  `gorm:"type:uuid;index" json:"company_id,omitempty"` // null for platform system accounts
	AccountCode string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"account_code"`
	AccountName string      `gorm:"type:varchar(255);not null" json:"account_name"`
	AccountType AccountType `gorm:"type:varchar(50);not null" json:"account_type"`
	Currency    string      `gorm:"type:varchar(10);default:'IDR'" json:"currency"`
	IsActive    bool        `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`

	Entries []LedgerEntry `gorm:"foreignKey:AccountID" json:"entries,omitempty"`
}

// Transaction represents an atomic financial journal containing at least 2 entries
type Transaction struct {
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	IdempotencyKey string          `gorm:"type:varchar(255);uniqueIndex;not null" json:"idempotency_key"`
	TxType         TransactionType `gorm:"type:varchar(50);not null" json:"tx_type"`
	ReferenceID    string          `gorm:"type:varchar(255);index" json:"reference_id"` // TenderID, TicketID, etc.
	Description    string          `gorm:"type:text" json:"description"`
	PostedAt       time.Time       `gorm:"not null" json:"posted_at"`
	CreatedAt      time.Time       `json:"created_at"`

	Entries []LedgerEntry `gorm:"foreignKey:TransactionID" json:"entries"`
}

// LedgerEntry represents a single debit or credit line
type LedgerEntry struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TransactionID uuid.UUID `gorm:"type:uuid;not null;index" json:"transaction_id"`
	AccountID     uuid.UUID `gorm:"type:uuid;not null;index" json:"account_id"`
	EntryType     EntryType `gorm:"type:varchar(10);not null" json:"entry_type"` // DEBIT or CREDIT
	Amount        int64     `gorm:"type:bigint;not null" json:"amount"`           // Always strictly positive, in lowest currency unit
	CreatedAt     time.Time `json:"created_at"`

	Account *Account `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

var (
	ErrUnbalancedJournal = errors.New("unbalanced journal: total debits must equal total credits")
	ErrNonPositiveAmount = errors.New("ledger entry amount must be strictly greater than zero")
	ErrDuplicatePosting  = errors.New("transaction contains duplicate or invalid account postings")
)

// ValidateZeroSum enforces the fundamental double-entry invariant: Sum(Debits) - Sum(Credits) == 0
func ValidateZeroSum(entries []LedgerEntry) error {
	if len(entries) < 2 {
		return errors.New("a double-entry transaction must contain at least 2 entries")
	}

	var totalDebit int64
	var totalCredit int64

	for _, entry := range entries {
		if entry.Amount <= 0 {
			return ErrNonPositiveAmount
		}
		if entry.EntryType == EntryDebit {
			totalDebit += entry.Amount
		} else if entry.EntryType == EntryCredit {
			totalCredit += entry.Amount
		} else {
			return errors.New("invalid entry type, must be DEBIT or CREDIT")
		}
	}

	if totalDebit != totalCredit {
		return ErrUnbalancedJournal
	}

	return nil
}
