package worker

import (
	"log"
	"time"

	"gorm.io/gorm"

	"scrapflow-backend/internal/usecase"
)

type ReconciliationWorker struct {
	db             *gorm.DB
	paymentUsecase usecase.PaymentUsecase
	interval       time.Duration
	stopChan       chan struct{}
}

func NewReconciliationWorker(db *gorm.DB, pUsecase usecase.PaymentUsecase, interval time.Duration) *ReconciliationWorker {
	if interval <= 0 {
		interval = 15 * time.Minute
	}
	return &ReconciliationWorker{
		db:             db,
		paymentUsecase: pUsecase,
		interval:       interval,
		stopChan:       make(chan struct{}),
	}
}

// Start launches the background worker loop in a goroutine
func (w *ReconciliationWorker) Start() {
	log.Printf("[ReconciliationWorker] Started background audit worker (Interval: %v)...", w.interval)
	ticker := time.NewTicker(w.interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				w.RunOnce()
			case <-w.stopChan:
				ticker.Stop()
				log.Println("[ReconciliationWorker] Stopped background audit worker.")
				return
			}
		}
	}()
}

// Stop terminates the background worker loop
func (w *ReconciliationWorker) Stop() {
	close(w.stopChan)
}

// RunOnce performs an immediate reconciliation audit
func (w *ReconciliationWorker) RunOnce() {
	log.Println("[ReconciliationWorker] Running automated bank & payment gateway reconciliation...")
	report, err := w.paymentUsecase.ReconcilePendingOrders(w.db)
	if err != nil {
		log.Printf("[ReconciliationWorker] Error during reconciliation: %v", err)
		return
	}
	log.Printf("[ReconciliationWorker] Batch %s Complete: Audited=%d, Matched=%d, Expired/Discrepancies=%d, TotalIDR=Rp %d",
		report.BatchID, report.AuditedCount, report.MatchedCount, report.DiscrepancyCount, report.TotalAmountIDR)
}
