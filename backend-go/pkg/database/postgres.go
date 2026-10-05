package database

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"scrapflow-backend/internal/domain"
)

func InitDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	log.Println("PostgreSQL connection established successfully.")

	// Auto-Migrate schema
	err = db.AutoMigrate(
		&domain.Company{},
		&domain.User{},
		&domain.Account{},
		&domain.Transaction{},
		&domain.LedgerEntry{},
		&domain.Tender{},
		&domain.Bid{},
		&domain.Contract{},
		&domain.WeighbridgeTicket{},
		&domain.FinancingFacility{},
		&domain.PaymentOrder{},
	)
	if err != nil {
		return nil, err
	}

	log.Println("Database AutoMigrate executed successfully.")

	seedSystemAccounts(db)

	return db, nil
}

func seedSystemAccounts(db *gorm.DB) {
	systemAccounts := []domain.Account{
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000001010"),
			AccountCode: "1010",
			AccountName: "Cash in Bank - Platform Custody",
			AccountType: domain.AccountAsset,
			Currency:    "IDR",
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000001020"),
			AccountCode: "1020",
			AccountName: "SCF Financing Loan Receivables",
			AccountType: domain.AccountAsset,
			Currency:    "IDR",
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000002010"),
			AccountCode: "2010",
			AccountName: "Escrow Bid-Bond Payable (Lapak Deposits)",
			AccountType: domain.AccountLiability,
			Currency:    "IDR",
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000002020"),
			AccountCode: "2020",
			AccountName: "Escrow Milestone Payable (Pabrik Contracts)",
			AccountType: domain.AccountLiability,
			Currency:    "IDR",
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000004010"),
			AccountCode: "4010",
			AccountName: "Escrow & Settlement Fee Revenue",
			AccountType: domain.AccountRevenue,
			Currency:    "IDR",
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000004020"),
			AccountCode: "4020",
			AccountName: "SCF Financing Fee & Interest Revenue",
			AccountType: domain.AccountRevenue,
			Currency:    "IDR",
			IsActive:    true,
		},
	}

	for _, acc := range systemAccounts {
		var existing domain.Account
		if err := db.Where("account_code = ?", acc.AccountCode).First(&existing).Error; err != nil {
			db.Create(&acc)
		}
	}
}
