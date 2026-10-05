package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/domain"
	"scrapflow-backend/internal/repository"
)

type PaymentUsecase interface {
	CreateVirtualAccount(db *gorm.DB, companyID uuid.UUID, amount int64, channel domain.PaymentChannel, paymentType string) (*domain.PaymentOrder, error)
	ProcessWebhook(db *gorm.DB, payload *WebhookPayload, signatureHeader string, secretKey string) (*domain.PaymentOrder, error)
	ReconcilePendingOrders(db *gorm.DB) (*domain.ReconciliationReport, error)
}

type WebhookPayload struct {
	ReferenceID   string `json:"reference_id" binding:"required"`
	TransactionID string `json:"transaction_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	PaymentStatus string `json:"payment_status" binding:"required"` // SETTLEMENT, PAID, EXPIRE
	Timestamp     int64  `json:"timestamp" binding:"required"`
}

type paymentUsecase struct {
	paymentRepo repository.PaymentRepository
	ledgerRepo  repository.LedgerRepository
}

func NewPaymentUsecase(pRepo repository.PaymentRepository, lRepo repository.LedgerRepository) PaymentUsecase {
	return &paymentUsecase{paymentRepo: pRepo, ledgerRepo: lRepo}
}

func (u *paymentUsecase) CreateVirtualAccount(db *gorm.DB, companyID uuid.UUID, amount int64, channel domain.PaymentChannel, paymentType string) (*domain.PaymentOrder, error) {
	if amount <= 0 {
		return nil, errors.New("payment amount must be greater than zero")
	}

	now := time.Now()
	refID := fmt.Sprintf("INV-VA-%s-%d", now.Format("20060102"), now.UnixNano()%100000)

	// Simulated Bank VA routing number
	prefix := "8808" // Default BCA
	if channel == domain.ChannelMandiriVA {
		prefix = "8888"
	} else if channel == domain.ChannelBRIVA {
		prefix = "7777"
	}
	vaNumber := fmt.Sprintf("%s%d", prefix, now.Unix()%1000000000)

	order := domain.PaymentOrder{
		ID:          uuid.New(),
		CompanyID:   companyID,
		ReferenceID: refID,
		PaymentType: paymentType,
		Channel:     channel,
		Amount:      amount,
		VANumber:    vaNumber,
		Status:      domain.PaymentPending,
		ExpiresAt:   now.Add(24 * time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.paymentRepo.CreatePaymentOrder(db, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (u *paymentUsecase) ProcessWebhook(db *gorm.DB, payload *WebhookPayload, signatureHeader string, secretKey string) (*domain.PaymentOrder, error) {
	// 1. HMAC-SHA256 Signature Verification
	if secretKey != "" && signatureHeader != "" {
		message := fmt.Sprintf("%s:%d:%d", payload.ReferenceID, payload.Amount, payload.Timestamp)
		h := hmac.New(sha256.New, []byte(secretKey))
		h.Write([]byte(message))
		expectedSignature := hex.EncodeToString(h.Sum(nil))

		if signatureHeader != expectedSignature {
			return nil, errors.New("invalid webhook signature: possible tampering detected")
		}
	}

	var updatedOrder *domain.PaymentOrder

	err := db.Transaction(func(tx *gorm.DB) error {
		// 2. Pessimistic lock row to prevent duplicate webhook processing
		order, err := u.paymentRepo.GetPaymentOrderByReferenceIDForUpdate(tx, payload.ReferenceID)
		if err != nil {
			return fmt.Errorf("payment order '%s' not found", payload.ReferenceID)
		}

		// Idempotency: If already paid, return early with 200 OK without double crediting
		if order.Status == domain.PaymentPaid {
			updatedOrder = order
			return nil
		}

		if payload.PaymentStatus == "SETTLEMENT" || payload.PaymentStatus == "PAID" {
			now := time.Now()
			order.Status = domain.PaymentPaid
			order.PaidAt = &now
			order.ExternalID = payload.TransactionID

			// 3. Post to Double-Entry Ledger
			bankAcc, err := u.ledgerRepo.GetAccountByCode(tx, "1010")
			if err != nil {
				return err
			}

			// Credit account depends on payment type
			targetAccountCode := "2020" // Escrow Milestone
			if order.PaymentType == "BID_BOND" {
				targetAccountCode = "2010" // Bid-Bond Payable
			}
			targetAcc, err := u.ledgerRepo.GetAccountByCode(tx, targetAccountCode)
			if err != nil {
				return err
			}

			entries := []domain.LedgerEntry{
				{AccountID: bankAcc.ID, EntryType: domain.EntryDebit, Amount: order.Amount},
				{AccountID: targetAcc.ID, EntryType: domain.EntryCredit, Amount: order.Amount},
			}

			idempotencyKey := fmt.Sprintf("WH-PAY-%s", order.ReferenceID)
			_, err = u.ledgerRepo.PostJournal(
				tx,
				idempotencyKey,
				domain.TxTypeEscrowLock,
				order.ReferenceID,
				fmt.Sprintf("Payment Gateway VA Settlement for %s", order.ReferenceID),
				entries,
			)
			if err != nil {
				return fmt.Errorf("failed to record payment to ledger: %w", err)
			}
		} else if payload.PaymentStatus == "EXPIRE" {
			order.Status = domain.PaymentExpired
		} else {
			order.Status = domain.PaymentFailed
		}

		if err := u.paymentRepo.UpdatePaymentOrder(tx, order); err != nil {
			return err
		}

		updatedOrder = order
		return nil
	})

	if err != nil {
		return nil, err
	}
	return updatedOrder, nil
}

func (u *paymentUsecase) ReconcilePendingOrders(db *gorm.DB) (*domain.ReconciliationReport, error) {
	now := time.Now()
	report := &domain.ReconciliationReport{
		BatchID:    fmt.Sprintf("REC-%s", now.Format("20060102150405")),
		ExecutedAt: now,
	}

	// Find pending orders past expiration
	pendingOrders, err := u.paymentRepo.GetPendingPaymentOrders(db, now, 100)
	if err != nil {
		return nil, err
	}

	report.AuditedCount = len(pendingOrders)

	for i := range pendingOrders {
		order := pendingOrders[i]
		if now.After(order.ExpiresAt) {
			order.Status = domain.PaymentExpired
			order.ReconciledAt = &now
			_ = u.paymentRepo.UpdatePaymentOrder(db, &order)
			report.DiscrepancyCount++
		} else {
			report.MatchedCount++
		}
		report.TotalAmountIDR += order.Amount
	}

	return report, nil
}
