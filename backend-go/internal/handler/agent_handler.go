package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"scrapflow-backend/internal/usecase"
	"scrapflow-backend/pkg/response"
)

type AgentHandler struct {
	db           *gorm.DB
	agentUsecase usecase.AgentUsecase
}

func NewAgentHandler(db *gorm.DB, aUsecase usecase.AgentUsecase) *AgentHandler {
	return &AgentHandler{db: db, agentUsecase: aUsecase}
}

type OCRSlipInput struct {
	ImageURL          string  `json:"image_url"`
	TruckLicensePlate string  `json:"truck_license_plate" binding:"required"`
	ManualGrossKG     float64 `json:"manual_gross_kg"`
	ManualTareKG      float64 `json:"manual_tare_kg"`
}

func (h *AgentHandler) AnalyzeSlip(c *gin.Context) {
	var req OCRSlipInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid OCR input", err.Error())
		return
	}

	result, err := h.agentUsecase.AnalyzeWeighbridgeSlip(h.db, req.ImageURL, req.TruckLicensePlate, req.ManualGrossKG, req.ManualTareKG)
	if err != nil {
		response.InternalError(c, "OCR slip analysis failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Weighbridge slip analyzed with AI fraud-detection engine", result)
}

type ParseTenderInput struct {
	DocumentContent string `json:"document_content" binding:"required"`
}

func (h *AgentHandler) ParseTender(c *gin.Context) {
	var req ParseTenderInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid document content", err.Error())
		return
	}

	result, err := h.agentUsecase.ParseTenderDocument(req.DocumentContent)
	if err != nil {
		response.InternalError(c, "Tender document parsing failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Tender specifications extracted by AI agent", result)
}
