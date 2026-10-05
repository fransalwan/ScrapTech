package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"scrapflow-backend/internal/config"
	"scrapflow-backend/internal/handler"
	"scrapflow-backend/internal/middleware"
	"scrapflow-backend/internal/repository"
	"scrapflow-backend/internal/usecase"
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

	// 2. Initialize Usecases
	authUsecase := usecase.NewAuthUsecase(companyRepo, ledgerRepo)
	tenderUsecase := usecase.NewTenderUsecase(tenderRepo, ledgerRepo)
	weighUsecase := usecase.NewWeighbridgeUsecase(weighRepo, tenderRepo, ledgerRepo, companyRepo)
	ledgerUsecase := usecase.NewLedgerUsecase(ledgerRepo)
	financingUsecase := usecase.NewFinancingUsecase(weighRepo, tenderRepo, ledgerRepo, companyRepo)

	// 3. Initialize HTTP Handlers
	authHandler := handler.NewAuthHandler(db, authUsecase, cfg.JWTSecret)
	tenderHandler := handler.NewTenderHandler(db, tenderUsecase)
	weighHandler := handler.NewWeighbridgeHandler(db, weighUsecase)
	ledgerHandler := handler.NewLedgerHandler(db, ledgerUsecase)
	financingHandler := handler.NewFinancingHandler(db, financingUsecase)

	// 4. Setup Engine
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "ScrapFlow Fintech Core API is healthy", gin.H{
			"version": "1.0.0-fintech-core",
			"status":  "ONLINE",
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
				// Bidding requires both Auth and Idempotency-Key
				tendersProtected.POST("/:id/bid", middleware.IdempotencyMiddleware(), tenderHandler.SubmitBid)
				// Awarding requires both Auth and Idempotency-Key
				tendersProtected.POST("/:id/award", middleware.IdempotencyMiddleware(), tenderHandler.AwardTender)
			}
		}

		// Weighbridge & Field Settlement Routes
		protectedAPI := apiV1.Group("")
		protectedAPI.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Weighbridge Tickets
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
		}
	}

	log.Fatal(r.Run(":" + cfg.Port))
}
