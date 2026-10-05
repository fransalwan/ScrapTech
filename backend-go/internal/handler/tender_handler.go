package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/usecase"
	"scrapflow-backend/pkg/response"
)

type TenderHandler struct {
	db            *gorm.DB
	tenderUsecase usecase.TenderUsecase
}

func NewTenderHandler(db *gorm.DB, tUsecase usecase.TenderUsecase) *TenderHandler {
	return &TenderHandler{db: db, tenderUsecase: tUsecase}
}

func (h *TenderHandler) CreateTender(c *gin.Context) {
	companyIDVal, _ := c.Get("company_id")
	companyID := companyIDVal.(uuid.UUID)

	var req usecase.PublishTenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid tender payload", err.Error())
		return
	}

	tender, err := h.tenderUsecase.PublishTender(h.db, companyID, &req)
	if err != nil {
		response.InternalError(c, "Failed to create tender", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Tender published successfully", tender)
}

func (h *TenderHandler) GetTenders(c *gin.Context) {
	status := c.Query("status")
	list, err := h.tenderUsecase.GetTenders(h.db, status)
	if err != nil {
		response.InternalError(c, "Failed to fetch tenders", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Tenders fetched successfully", list)
}

func (h *TenderHandler) GetTenderDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid tender ID format", nil)
		return
	}

	tender, err := h.tenderUsecase.GetTenderDetail(h.db, id)
	if err != nil {
		response.NotFound(c, "Tender not found")
		return
	}

	response.Success(c, http.StatusOK, "Tender detail fetched", tender)
}

type SubmitBidInput struct {
	PricePerKG int64 `json:"price_per_kg" binding:"required,gt=0"`
}

func (h *TenderHandler) SubmitBid(c *gin.Context) {
	companyIDVal, _ := c.Get("company_id")
	companyID := companyIDVal.(uuid.UUID)

	idempotencyKeyVal, _ := c.Get("idempotency_key")
	idempotencyKey := idempotencyKeyVal.(string)

	tenderIDStr := c.Param("id")
	tenderID, err := uuid.Parse(tenderIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid tender ID format", nil)
		return
	}

	var req SubmitBidInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid bid input", err.Error())
		return
	}

	bid, err := h.tenderUsecase.SubmitBid(h.db, idempotencyKey, companyID, tenderID, req.PricePerKG)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Bid submitted and bid-bond safely locked in escrow", bid)
}

type AwardTenderInput struct {
	WinningBidID uuid.UUID `json:"winning_bid_id" binding:"required"`
}

func (h *TenderHandler) AwardTender(c *gin.Context) {
	idempotencyKeyVal, _ := c.Get("idempotency_key")
	idempotencyKey := idempotencyKeyVal.(string)

	tenderIDStr := c.Param("id")
	tenderID, err := uuid.Parse(tenderIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid tender ID format", nil)
		return
	}

	var req AwardTenderInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid award input", err.Error())
		return
	}

	contract, err := h.tenderUsecase.AwardTender(h.db, idempotencyKey, tenderID, req.WinningBidID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Tender awarded and official SPK contract created", contract)
}
