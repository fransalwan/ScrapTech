package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/usecase"
	"scrapflow-backend/pkg/response"
)

type WeighbridgeHandler struct {
	db           *gorm.DB
	weighUsecase usecase.WeighbridgeUsecase
}

func NewWeighbridgeHandler(db *gorm.DB, wUsecase usecase.WeighbridgeUsecase) *WeighbridgeHandler {
	return &WeighbridgeHandler{db: db, weighUsecase: wUsecase}
}

func (h *WeighbridgeHandler) SubmitTicket(c *gin.Context) {
	contractIDStr := c.Param("contract_id")
	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid contract ID format", nil)
		return
	}

	var req usecase.SubmitTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid ticket payload", err.Error())
		return
	}

	ticket, err := h.weighUsecase.SubmitTicket(h.db, contractID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Weighbridge ticket registered successfully", ticket)
}

func (h *WeighbridgeHandler) GetTicketsByContract(c *gin.Context) {
	contractIDStr := c.Param("contract_id")
	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid contract ID format", nil)
		return
	}

	tickets, err := h.weighUsecase.GetTicketsByContract(h.db, contractID)
	if err != nil {
		response.InternalError(c, "Failed to fetch tickets", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Tickets fetched successfully", tickets)
}

func (h *WeighbridgeHandler) SettleTicket(c *gin.Context) {
	ticketIDStr := c.Param("ticket_id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid ticket ID format", nil)
		return
	}

	idempotencyKeyVal, _ := c.Get("idempotency_key")
	idempotencyKey := idempotencyKeyVal.(string)

	settled, err := h.weighUsecase.VerifyAndSettleTicket(h.db, idempotencyKey, ticketID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Weighbridge ticket atomically settled via Double-Entry Ledger", settled)
}
