package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/usecase"
	"scrapflow-backend/pkg/response"
)

type FinancingHandler struct {
	db               *gorm.DB
	financingUsecase usecase.FinancingUsecase
}

func NewFinancingHandler(db *gorm.DB, fUsecase usecase.FinancingUsecase) *FinancingHandler {
	return &FinancingHandler{db: db, financingUsecase: fUsecase}
}

type ApplyFacilityInput struct {
	ContractID      uuid.UUID `json:"contract_id" binding:"required"`
	RequestedAmount int64     `json:"requested_amount" binding:"required,gt=0"`
	TenorDays       int       `json:"tenor_days" binding:"required,oneof=14 30 45 60"`
}

func (h *FinancingHandler) ApplyFacility(c *gin.Context) {
	companyIDVal, _ := c.Get("company_id")
	companyID := companyIDVal.(uuid.UUID)

	var req ApplyFacilityInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid financing application input", err.Error())
		return
	}

	facility, err := h.financingUsecase.ApplyFacility(h.db, companyID, req.ContractID, req.RequestedAmount, req.TenorDays)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Supply chain financing application submitted", facility)
}

func (h *FinancingHandler) GetMyFacilities(c *gin.Context) {
	companyIDVal, _ := c.Get("company_id")
	companyID := companyIDVal.(uuid.UUID)

	list, err := h.financingUsecase.GetFacilitiesByBorrower(h.db, companyID)
	if err != nil {
		response.InternalError(c, "Failed to retrieve facilities", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Financing facilities retrieved", list)
}

func (h *FinancingHandler) DisburseFacility(c *gin.Context) {
	idempotencyKeyVal, _ := c.Get("idempotency_key")
	idempotencyKey := idempotencyKeyVal.(string)

	facilityIDStr := c.Param("id")
	facilityID, err := uuid.Parse(facilityIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid facility ID format", nil)
		return
	}

	facility, err := h.financingUsecase.DisburseFacility(h.db, idempotencyKey, facilityID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Financing facility disbursed into tender escrow vault", facility)
}
