package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"scrapflow-backend/internal/domain"
	"scrapflow-backend/internal/repository"
)

type WeighbridgeUsecase interface {
	SubmitTicket(db *gorm.DB, contractID uuid.UUID, req *SubmitTicketRequest) (*domain.WeighbridgeTicket, error)
	GetTicketsByContract(db *gorm.DB, contractID uuid.UUID) ([]domain.WeighbridgeTicket, error)
	VerifyAndSettleTicket(db *gorm.DB, idempotencyKey string, ticketID uuid.UUID) (*domain.WeighbridgeTicket, error)
}

type SubmitTicketRequest struct {
	TruckLicensePlate    string  `json:"truck_license_plate" binding:"required"`
	DriverName           string  `json:"driver_name"`
	GrossWeightKG        float64 `json:"gross_weight_kg" binding:"required,gt=0"`
	TareWeightKG         float64 `json:"tare_weight_kg" binding:"required,gt=0"`
	RefractionPercentage float64 `json:"refraction_percentage" binding:"gte=0,lte=100"`
	SlipPhotoURL         string  `json:"slip_photo_url"`
}

type weighbridgeUsecase struct {
	weighRepo   repository.WeighbridgeRepository
	tenderRepo  repository.TenderRepository
	ledgerRepo  repository.LedgerRepository
	companyRepo repository.CompanyRepository
}

func NewWeighbridgeUsecase(
	wRepo repository.WeighbridgeRepository,
	tRepo repository.TenderRepository,
	lRepo repository.LedgerRepository,
	cRepo repository.CompanyRepository,
) WeighbridgeUsecase {
	return &weighbridgeUsecase{
		weighRepo:   wRepo,
		tenderRepo:  tRepo,
		ledgerRepo:  lRepo,
		companyRepo: cRepo,
	}
}

func (u *weighbridgeUsecase) SubmitTicket(db *gorm.DB, contractID uuid.UUID, req *SubmitTicketRequest) (*domain.WeighbridgeTicket, error) {
	contract, err := u.tenderRepo.GetContractByID(db, contractID)
	if err != nil {
		return nil, errors.New("contract not found")
	}

	if contract.Status != domain.ContractActive {
		return nil, errors.New("cannot submit weighbridge tickets for inactive contract")
	}

	if req.GrossWeightKG <= req.TareWeightKG {
		return nil, errors.New("gross weight must be strictly greater than tare weight")
	}

	now := time.Now()
	ticket := domain.WeighbridgeTicket{
		ID:                   uuid.New(),
		ContractID:           contractID,
		TicketNumber:         fmt.Sprintf("WB-%s-%d", now.Format("20060102"), now.UnixNano()%10000),
		TruckLicensePlate:    req.TruckLicensePlate,
		DriverName:           req.DriverName,
		GrossWeightKG:        req.GrossWeightKG,
		TareWeightKG:         req.TareWeightKG,
		RefractionPercentage: req.RefractionPercentage,
		SlipPhotoURL:         req.SlipPhotoURL,
		Status:               domain.TicketVerified, // Auto-verified for streamlined flow
		WeighedAt:            now,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	// Calculate billable weight and amounts (0.5% platform fee)
	ticket.CalculateWeights(contract.AgreedPricePerKG, 0.5)

	if err := u.weighRepo.CreateTicket(db, &ticket); err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (u *weighbridgeUsecase) GetTicketsByContract(db *gorm.DB, contractID uuid.UUID) ([]domain.WeighbridgeTicket, error) {
	return u.weighRepo.GetTicketsByContract(db, contractID)
}

func (u *weighbridgeUsecase) VerifyAndSettleTicket(db *gorm.DB, idempotencyKey string, ticketID uuid.UUID) (*domain.WeighbridgeTicket, error) {
	var settledTicket *domain.WeighbridgeTicket

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Pessimistic lock ticket
		ticket, err := u.weighRepo.GetTicketByIDForUpdate(tx, ticketID)
		if err != nil {
			return errors.New("weighbridge ticket not found")
		}

		if ticket.Status == domain.TicketSettled {
			return errors.New("ticket has already been settled")
		}

		// 2. Pessimistic lock contract
		contract, err := u.tenderRepo.GetContractByIDForUpdate(tx, ticket.ContractID)
		if err != nil {
			return errors.New("contract not found")
		}

		// 3. Resolve accounts for Double-Entry Settlement
		// Milestone Escrow Account (2020)
		escrowAcc, err := u.ledgerRepo.GetAccountByCode(tx, "2020")
		if err != nil {
			return err
		}

		// Platform Fee Account (4010)
		feeAcc, err := u.ledgerRepo.GetAccountByCode(tx, "4010")
		if err != nil {
			return err
		}

		// Seller (Pabrik) Withdrawable Wallet Account
		sellerWalletAcc, err := u.ledgerRepo.GetOrCreateCompanyAccount(
			tx,
			contract.SellerCompanyID,
			domain.AccountLiability,
			"WALLET-"+contract.SellerCompanyID.String()[:8],
			"Seller Withdrawable Balance",
		)
		if err != nil {
			return err
		}

		// 4. Double-Entry Posting:
		// [DEBIT]  2020 (Escrow Milestone Payable)       : GrossPayableAmount
		// [CREDIT] Wallet (Seller Withdrawable Wallet)   : NetToSellerAmount
		// [CREDIT] 4010 (Platform Escrow Fee Revenue)    : PlatformFeeAmount
		entries := []domain.LedgerEntry{
			{
				AccountID: escrowAcc.ID,
				EntryType: domain.EntryDebit,
				Amount:    ticket.GrossPayableAmount,
			},
			{
				AccountID: sellerWalletAcc.ID,
				EntryType: domain.EntryCredit,
				Amount:    ticket.NetToSellerAmount,
			},
			{
				AccountID: feeAcc.ID,
				EntryType: domain.EntryCredit,
				Amount:    ticket.PlatformFeeAmount,
			},
		}

		_, err = u.ledgerRepo.PostJournal(
			tx,
			idempotencyKey,
			domain.TxTypeWeighbridgeSettle,
			ticket.TicketNumber,
			fmt.Sprintf("Weighbridge settlement %s (%.2f Kg)", ticket.TicketNumber, ticket.FinalBillableWeightKG),
			entries,
		)
		if err != nil {
			return fmt.Errorf("ledger settlement failed: %w", err)
		}

		// 5. Update Ticket and Contract state
		now := time.Now()
		ticket.Status = domain.TicketSettled
		ticket.SettledAt = &now
		if err := u.weighRepo.UpdateTicket(tx, ticket); err != nil {
			return err
		}

		contract.TotalSettledAmount += ticket.GrossPayableAmount
		if err := u.tenderRepo.UpdateContract(tx, contract); err != nil {
			return err
		}

		settledTicket = ticket
		return nil
	})

	if err != nil {
		return nil, err
	}
	return settledTicket, nil
}
