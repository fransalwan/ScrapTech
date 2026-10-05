package domain

type OCRAnalysisResult struct {
	SlipNumber           string  `json:"slip_number"`
	TruckLicensePlate    string  `json:"truck_license_plate"`
	GrossWeightKG        float64 `json:"gross_weight_kg"`
	TareWeightKG         float64 `json:"tare_weight_kg"`
	NetWeightKG          float64 `json:"net_weight_kg"`
	HistoricalTareAvgKG  float64 `json:"historical_tare_avg_kg"`
	TareDeviationPct     float64 `json:"tare_deviation_pct"`
	FraudAnomalyDetected bool    `json:"fraud_anomaly_detected"`
	FraudWarningMessage  string  `json:"fraud_warning_message,omitempty"`
	ConfidenceScore      float64 `json:"confidence_score"`
}

type TenderParseResult struct {
	ExtractedTitle        string        `json:"extracted_title"`
	DetectedScrapCategory ScrapCategory `json:"detected_scrap_category"`
	EstimatedWeightTon    float64       `json:"estimated_weight_ton"`
	EstimatedWeightKG     float64       `json:"estimated_weight_kg"`
	RecommendedReserveKG  int64         `json:"recommended_reserve_kg"`
	SuggestedBidBondIDR   int64         `json:"suggested_bid_bond_idr"`
	KeyTerms              []string      `json:"key_terms"`
}
