package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidateZeroSum_Success(t *testing.T) {
	acc1 := uuid.New()
	acc2 := uuid.New()

	entries := []LedgerEntry{
		{AccountID: acc1, EntryType: EntryDebit, Amount: 50000000},
		{AccountID: acc2, EntryType: EntryCredit, Amount: 50000000},
	}

	err := ValidateZeroSum(entries)
	if err != nil {
		t.Fatalf("expected valid zero-sum journal, got error: %v", err)
	}
}

func TestValidateZeroSum_MultiSplitSuccess(t *testing.T) {
	accEscrow := uuid.New()
	accSeller := uuid.New()
	accFee := uuid.New()

	// Split settlement: 100M total = 99.5M to seller + 0.5M platform fee
	entries := []LedgerEntry{
		{AccountID: accEscrow, EntryType: EntryDebit, Amount: 100000000},
		{AccountID: accSeller, EntryType: EntryCredit, Amount: 99500000},
		{AccountID: accFee, EntryType: EntryCredit, Amount: 500000},
	}

	err := ValidateZeroSum(entries)
	if err != nil {
		t.Fatalf("expected valid zero-sum multi-split journal, got error: %v", err)
	}
}

func TestValidateZeroSum_Unbalanced(t *testing.T) {
	acc1 := uuid.New()
	acc2 := uuid.New()

	entries := []LedgerEntry{
		{AccountID: acc1, EntryType: EntryDebit, Amount: 50000000},
		{AccountID: acc2, EntryType: EntryCredit, Amount: 49000000}, // 1M discrepancy!
	}

	err := ValidateZeroSum(entries)
	if err == nil {
		t.Fatal("expected ErrUnbalancedJournal error, got nil")
	}

	if err != ErrUnbalancedJournal {
		t.Fatalf("expected %v, got %v", ErrUnbalancedJournal, err)
	}
}

func TestValidateZeroSum_NonPositiveAmount(t *testing.T) {
	acc1 := uuid.New()
	acc2 := uuid.New()

	entries := []LedgerEntry{
		{AccountID: acc1, EntryType: EntryDebit, Amount: -1000},
		{AccountID: acc2, EntryType: EntryCredit, Amount: -1000},
	}

	err := ValidateZeroSum(entries)
	if err != ErrNonPositiveAmount {
		t.Fatalf("expected ErrNonPositiveAmount, got %v", err)
	}
}

func TestValidateZeroSum_InsufficientEntries(t *testing.T) {
	acc1 := uuid.New()

	entries := []LedgerEntry{
		{AccountID: acc1, EntryType: EntryDebit, Amount: 100000},
	}

	err := ValidateZeroSum(entries)
	if err == nil {
		t.Fatal("expected error for single-entry journal, got nil")
	}
}
