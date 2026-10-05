package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/domain"
	"scrapflow-backend/internal/usecase"
	"scrapflow-backend/pkg/response"
)

type PaymentHandler struct {
	db             *gorm.DB
	paymentUsecase usecase.PaymentUsecase
	webhookSecret  string
}

func NewPaymentHandler(db *gorm.DB, pUsecase usecase.PaymentUsecase, webhookSecret string) *PaymentHandler {
	return &PaymentHandler{
		db:             db,
		paymentUsecase: pUsecase,
		webhookSecret:  webhookSecret,
	}
}

type CreateVAInput struct {
	Amount      int64                 `json:"amount" binding:"required,gt=0"`
	Channel     domain.PaymentChannel `json:"channel" binding:"required"`
	PaymentType string                `json:"payment_type" binding:"required"` // BID_BOND or ESCROW_DEPOSIT
}

func (h *PaymentHandler) CreateVirtualAccount(c *gin.Context) {
	companyIDVal, _ := c.Get("company_id")
	companyID := companyIDVal.(uuid.UUID)

	var req CreateVAInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid virtual account request", err.Error())
		return
	}

	order, err := h.paymentUsecase.CreateVirtualAccount(h.db, companyID, req.Amount, req.Channel, req.PaymentType)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Virtual Account created successfully", order)
}

func (h *PaymentHandler) ProcessWebhook(c *gin.Context) {
	signature := c.GetHeader("X-Signature")

	var payload usecase.WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, "Invalid webhook payload structure", err.Error())
		return
	}

	order, err := h.paymentUsecase.ProcessWebhook(h.db, &payload, signature, h.webhookSecret)
	if err != nil {
		response.BadRequest(c, "Webhook processing failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Webhook acknowledged and settled in double-entry ledger", order)
}

func (h *PaymentHandler) TriggerReconciliation(c *gin.Context) {
	report, err := h.paymentUsecase.ReconcilePendingOrders(h.db)
	if err != nil {
		response.InternalError(c, "Reconciliation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Reconciliation audit executed", report)
}
