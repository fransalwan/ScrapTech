package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"scrapflow-backend/internal/config"
	"scrapflow-backend/internal/handler"
	"scrapflow-backend/internal/middleware"
	"scrapflow-backend/internal/repository"
	"scrapflow-backend/internal/usecase"
	"scrapflow-backend/internal/worker"
	"scrapflow-backend/pkg/database"
	"scrapflow-backend/pkg/response"
)

func main() {
	cfg := config.LoadConfig()

	log.Printf("Starting ScrapFlow B2B Commodity Fintech Engine on port %s...", cfg.Port)

	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Fatal: Database initialization failed: %v", err)
	}

	// 1. Initialize Repositories
	companyRepo := repository.NewCompanyRepository()
	ledgerRepo := repository.NewLedgerRepository()
	tenderRepo := repository.NewTenderRepository()
	weighRepo := repository.NewWeighbridgeRepository()
	paymentRepo := repository.NewPaymentRepository()

	// 2. Initialize Usecases
	authUsecase := usecase.NewAuthUsecase(companyRepo, ledgerRepo)
	tenderUsecase := usecase.NewTenderUsecase(tenderRepo, ledgerRepo)
	weighUsecase := usecase.NewWeighbridgeUsecase(weighRepo, tenderRepo, ledgerRepo, companyRepo)
	ledgerUsecase := usecase.NewLedgerUsecase(ledgerRepo)
	financingUsecase := usecase.NewFinancingUsecase(weighRepo, tenderRepo, ledgerRepo, companyRepo)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, ledgerRepo)
	agentUsecase := usecase.NewAgentUsecase(weighRepo)

	// 3. Initialize & Start Background Reconciliation Worker (Audits every 15 mins)
	reconcileWorker := worker.NewReconciliationWorker(db, paymentUsecase, 15*time.Minute)
	reconcileWorker.Start()
	defer reconcileWorker.Stop()

	// 4. Initialize HTTP Handlers
	authHandler := handler.NewAuthHandler(db, authUsecase, cfg.JWTSecret)
	tenderHandler := handler.NewTenderHandler(db, tenderUsecase)
	weighHandler := handler.NewWeighbridgeHandler(db, weighUsecase)
	ledgerHandler := handler.NewLedgerHandler(db, ledgerUsecase)
	financingHandler := handler.NewFinancingHandler(db, financingUsecase)
	paymentHandler := handler.NewPaymentHandler(db, paymentUsecase, cfg.JWTSecret)
	agentHandler := handler.NewAgentHandler(db, agentUsecase)

	// 5. Setup Engine
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "ScrapFlow Fintech Core API is healthy", gin.H{
			"version":        "1.0.0-fintech-core",
			"status":         "ONLINE",
			"ledger_engine":  "ACTIVE",
			"audit_worker":   "RUNNING",
			"ai_agent_layer": "ENABLED",
		})
	})

	apiV1 := r.Group("/api/v1")
	{
		// Auth Routes (Public)
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Public & Protected Tender Routes
		tendersGroup := apiV1.Group("/tenders")
		{
			tendersGroup.GET("", tenderHandler.GetTenders)
			tendersGroup.GET("/:id", tenderHandler.GetTenderDetail)

			// Protected routes
			tendersProtected := tendersGroup.Group("")
			tendersProtected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				tendersProtected.POST("", tenderHandler.CreateTender)
				tendersProtected.POST("/:id/bid", middleware.IdempotencyMiddleware(), tenderHandler.SubmitBid)
				tendersProtected.POST("/:id/award", middleware.IdempotencyMiddleware(), tenderHandler.AwardTender)
			}
		}

		// Public Webhook Ingestion (HMAC Signature Verified)
		apiV1.POST("/payments/webhook", paymentHandler.ProcessWebhook)

		// Protected Operations
		protectedAPI := apiV1.Group("")
		protectedAPI.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Weighbridge Tickets & Field Settlement
			protectedAPI.POST("/contracts/:contract_id/tickets", weighHandler.SubmitTicket)
			protectedAPI.GET("/contracts/:contract_id/tickets", weighHandler.GetTicketsByContract)
			protectedAPI.POST("/tickets/:ticket_id/settle", middleware.IdempotencyMiddleware(), weighHandler.SettleTicket)

			// Ledger & Wallet Balance
			protectedAPI.GET("/ledger/balance", ledgerHandler.GetWalletBalance)
			protectedAPI.GET("/ledger/transactions", ledgerHandler.GetTransactions)
			protectedAPI.POST("/ledger/deposit", middleware.IdempotencyMiddleware(), ledgerHandler.DepositEscrow)

			// Supply Chain Financing
			protectedAPI.POST("/financing/apply", financingHandler.ApplyFacility)
			protectedAPI.GET("/financing/my-facilities", financingHandler.GetMyFacilities)
			protectedAPI.POST("/financing/:id/disburse", middleware.IdempotencyMiddleware(), financingHandler.DisburseFacility)

			// Payment Gateway & Reconciliation Trigger
			protectedAPI.POST("/payments/va/create", paymentHandler.CreateVirtualAccount)
			protectedAPI.POST("/payments/reconcile", paymentHandler.TriggerReconciliation)

			// AI Agents (Weighbridge OCR Anomaly & Tender Spec Parsing)
			protectedAPI.POST("/agents/ocr-slip", agentHandler.AnalyzeSlip)
			protectedAPI.POST("/agents/parse-tender", agentHandler.ParseTender)
		}
	}

	log.Fatal(r.Run(":" + cfg.Port))
}
