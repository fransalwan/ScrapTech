package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/usecase"
	"scrapflow-backend/pkg/response"
)

type LedgerHandler struct {
	db            *gorm.DB
	ledgerUsecase usecase.LedgerUsecase
}

func NewLedgerHandler(db *gorm.DB, lUsecase usecase.LedgerUsecase) *LedgerHandler {
	return &LedgerHandler{db: db, ledgerUsecase: lUsecase}
}

func (h *LedgerHandler) GetWalletBalance(c *gin.Context) {
	companyIDVal, _ := c.Get("company_id")
	companyID := companyIDVal.(uuid.UUID)

	balance, err := h.ledgerUsecase.GetCompanyWalletBalance(h.db, companyID)
	if err != nil {
		response.InternalError(c, "Failed to retrieve balance", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Wallet balance retrieved", gin.H{
		"company_id": companyID,
		"balance":    balance,
		"currency":   "IDR",
	})
}

func (h *LedgerHandler) GetTransactions(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	txs, err := h.ledgerUsecase.GetLedgerTransactions(h.db, limit, offset)
	if err != nil {
		response.InternalError(c, "Failed to retrieve ledger transactions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Ledger transactions retrieved", txs)
}

type DepositInput struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

func (h *LedgerHandler) DepositEscrow(c *gin.Context) {
	companyIDVal, _ := c.Get("company_id")
	companyID := companyIDVal.(uuid.UUID)

	idempotencyKeyVal, _ := c.Get("idempotency_key")
	idempotencyKey := idempotencyKeyVal.(string)

	var req DepositInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid deposit amount", err.Error())
		return
	}

	tx, err := h.ledgerUsecase.DepositEscrowVault(h.db, idempotencyKey, companyID, req.Amount)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Escrow vault deposit posted to double-entry ledger", tx)
}
