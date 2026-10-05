package usecase

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"scrapflow-backend/internal/domain"
	"scrapflow-backend/internal/repository"
)

type AgentUsecase interface {
	AnalyzeWeighbridgeSlip(db *gorm.DB, imageURL, truckPlate string, manualGross, manualTare float64) (*domain.OCRAnalysisResult, error)
	ParseTenderDocument(rawContent string) (*domain.TenderParseResult, error)
}

type agentUsecase struct {
	weighRepo repository.WeighbridgeRepository
}

func NewAgentUsecase(wRepo repository.WeighbridgeRepository) AgentUsecase {
	return &agentUsecase{weighRepo: wRepo}
}

func (u *agentUsecase) AnalyzeWeighbridgeSlip(db *gorm.DB, imageURL, truckPlate string, manualGross, manualTare float64) (*domain.OCRAnalysisResult, error) {
	// Baseline historical tare for trucks (in real system, aggregated from past WeighbridgeTicket records)
	historicalTare := 11200.0 // default baseline 11.2 Ton
	var pastTickets []domain.WeighbridgeTicket
	db.Where("truck_license_plate = ?", truckPlate).Limit(5).Find(&pastTickets)
	if len(pastTickets) > 0 {
		var totalTare float64
		for _, t := range pastTickets {
			totalTare += t.TareWeightKG
		}
		historicalTare = totalTare / float64(len(pastTickets))
	}

	// Simulated high-fidelity OCR reading
	extractedGross := manualGross
	if extractedGross <= 0 {
		extractedGross = 28450.0
	}
	extractedTare := manualTare
	if extractedTare <= 0 {
		extractedTare = 11200.0
	}

	extractedNet := math.Max(0, extractedGross-extractedTare)
	deviationPct := math.Abs(extractedTare-historicalTare) / historicalTare * 100.0

	fraudDetected := false
	fraudMsg := ""
	if deviationPct > 3.0 {
		fraudDetected = true
		fraudMsg = fmt.Sprintf("CRITICAL ANOMALY: Truk %s memiliki deviasi berat tara %.2f%% dari rata-rata historis (%.1f Kg vs %.1f Kg). Indikasi manipulasi air tangki atau modifikasi tara sebelum timbang.",
			truckPlate, deviationPct, extractedTare, historicalTare)
	}

	return &domain.OCRAnalysisResult{
		SlipNumber:           fmt.Sprintf("OCR-WB-%s", truckPlate),
		TruckLicensePlate:    truckPlate,
		GrossWeightKG:        extractedGross,
		TareWeightKG:         extractedTare,
		NetWeightKG:          extractedNet,
		HistoricalTareAvgKG:  historicalTare,
		TareDeviationPct:     math.Round(deviationPct*100) / 100,
		FraudAnomalyDetected: fraudDetected,
		FraudWarningMessage:  fraudMsg,
		ConfidenceScore:      0.978,
	}, nil
}

func (u *agentUsecase) ParseTenderDocument(rawContent string) (*domain.TenderParseResult, error) {
	upper := strings.ToUpper(rawContent)

	category := domain.ScrapHMS1
	if strings.Contains(upper, "PLATE") || strings.Contains(upper, "PLAT") {
		category = domain.ScrapPlate
	} else if strings.Contains(upper, "CAST IRON") || strings.Contains(upper, "COR") {
		category = domain.ScrapCastIron
	} else if strings.Contains(upper, "COPPER") || strings.Contains(upper, "TEMBAGA") {
		category = domain.ScrapCopper
	} else if strings.Contains(upper, "HMS 2") || strings.Contains(upper, "HMS-2") {
		category = domain.ScrapHMS2
	}

	// Regex search for tonnage (e.g. "450 TON", "300 TON", "150.000 KG")
	tonnage := 250.0 // default
	re := regexp.MustCompile(`(\d+[\.,]?\d*)\s*(TON|TONS|KG|KILOGRAM)`)
	matches := re.FindStringSubmatch(upper)
	if len(matches) >= 3 {
		numStr := strings.ReplaceAll(matches[1], ",", ".")
		val, err := strconv.ParseFloat(numStr, 64)
		if err == nil {
			if strings.Contains(matches[2], "KG") {
				tonnage = val / 1000.0
			} else {
				tonnage = val
			}
		}
	}

	weightKG := tonnage * 1000.0
	reservePrice := int64(6200)
	if category == domain.ScrapPlate {
		reservePrice = 6700
	} else if category == domain.ScrapCopper {
		reservePrice = 95000
	} else if category == domain.ScrapCastIron {
		reservePrice = 5800
	}

	totalEstValue := int64(weightKG * float64(reservePrice))
	bidBond := int64(float64(totalEstValue) * 0.05) // 5% bid bond

	return &domain.TenderParseResult{
		ExtractedTitle:        fmt.Sprintf("Lelang Terbuka Scrap %s (Est. %.1f Ton)", category, tonnage),
		DetectedScrapCategory: category,
		EstimatedWeightTon:    tonnage,
		EstimatedWeightKG:     weightKG,
		RecommendedReserveKG:  reservePrice,
		SuggestedBidBondIDR:   bidBond,
		KeyTerms: []string{
			"Syarat K3 Lapangan Wajib",
			"Pembayaran per Ritase Timbangan",
			"Jaminan Bid-Bond Ditahan di Escrow",
		},
	}, nil
}
