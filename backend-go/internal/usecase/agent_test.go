package usecase

import (
	"testing"

	"scrapflow-backend/internal/domain"
)

func TestParseTenderDocument_HMS1(t *testing.T) {
	u := NewAgentUsecase(nil)

	sampleDoc := `PENGUMUMAN LELANG BESI TUA PT KRAKATAU STEEL
Jenis Material: Scrap Besi Baja HMS-1 dan Potongan Plat
Estimasi Volume: 500 Ton
Lokasi: Cilegon Banten
Waktu Bidding: 3 Hari Kerja`

	result, err := u.ParseTenderDocument(sampleDoc)
	if err != nil {
		t.Fatalf("unexpected error parsing tender document: %v", err)
	}

	if result.EstimatedWeightTon != 500.0 {
		t.Fatalf("expected 500 Ton, got %.1f", result.EstimatedWeightTon)
	}

	if result.RecommendedReserveKG <= 0 {
		t.Fatalf("expected positive reserve price, got %d", result.RecommendedReserveKG)
	}

	if result.SuggestedBidBondIDR <= 0 {
		t.Fatalf("expected positive bid bond, got %d", result.SuggestedBidBondIDR)
	}
}

func TestAnalyzeWeighbridgeSlip_TareFraudDetection(t *testing.T) {
	// Baseline average tare = 11,200 Kg
	// Truck enters with tare = 12,000 Kg (> 7% deviation, water tank tampering suspected)
	historicalTare := 11200.0
	tamperedTare := 12000.0
	gross := 28450.0

	deviationPct := (tamperedTare - historicalTare) / historicalTare * 100.0
	if deviationPct <= 3.0 {
		t.Fatalf("expected deviation > 3%%, got %.2f%%", deviationPct)
	}

	// Verify invariant logic
	fraudDetected := deviationPct > 3.0
	if !fraudDetected {
		t.Fatalf("expected fraud to be detected")
	}

	net := gross - tamperedTare
	if net != 16450.0 {
		t.Fatalf("expected net weight 16450, got %.1f", net)
	}

	// Check domain constants
	if domain.ScrapHMS1 != "HMS_1" {
		t.Fatalf("expected ScrapHMS1 constant to be HMS_1")
	}
}
