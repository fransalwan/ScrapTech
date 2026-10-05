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

type TenderUsecase interface {
	PublishTender(db *gorm.DB, companyID uuid.UUID, req *PublishTenderRequest) (*domain.Tender, error)
	GetTenders(db *gorm.DB, status string) ([]domain.Tender, error)
	GetTenderDetail(db *gorm.DB, id uuid.UUID) (*domain.Tender, error)
	SubmitBid(db *gorm.DB, idempotencyKey string, companyID uuid.UUID, tenderID uuid.UUID, pricePerKG int64) (*domain.Bid, error)
	AwardTender(db *gorm.DB, idempotencyKey string, tenderID uuid.UUID, winningBidID uuid.UUID) (*domain.Contract, error)
}

type PublishTenderRequest struct {
	Title             string               `json:"title" binding:"required"`
	ScrapType         domain.ScrapCategory `json:"scrap_type" binding:"required"`
	EstimatedWeightKG float64              `json:"estimated_weight_kg" binding:"required,gt=0"`
	ReservePricePerKG int64                `json:"reserve_price_per_kg" binding:"required,gt=0"`
	BidBondAmount     int64                `json:"bid_bond_amount" binding:"required,gt=0"`
	Location          string               `json:"location" binding:"required"`
	BiddingDurationHr int                  `json:"bidding_duration_hr" binding:"required,gt=0"`
}

type tenderUsecase struct {
	tenderRepo repository.TenderRepository
	ledgerRepo repository.LedgerRepository
}

func NewTenderUsecase(tRepo repository.TenderRepository, lRepo repository.LedgerRepository) TenderUsecase {
	return &tenderUsecase{tenderRepo: tRepo, ledgerRepo: lRepo}
}

func (u *tenderUsecase) PublishTender(db *gorm.DB, companyID uuid.UUID, req *PublishTenderRequest) (*domain.Tender, error) {
	now := time.Now()
	tender := domain.Tender{
		ID:                 uuid.New(),
		PublisherCompanyID: companyID,
		TenderNumber:       fmt.Sprintf("TND-%s-%d", now.Format("20060102"), now.UnixNano()%10000),
		Title:              req.Title,
		ScrapType:          req.ScrapType,
		EstimatedWeightKG:  req.EstimatedWeightKG,
		ReservePricePerKG:  req.ReservePricePerKG,
		BidBondAmount:      req.BidBondAmount,
		Location:           req.Location,
		BiddingDeadline:    now.Add(time.Duration(req.BiddingDurationHr) * time.Hour),
		Status:             domain.TenderBidding,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := u.tenderRepo.CreateTender(db, &tender); err != nil {
		return nil, err
	}
	return &tender, nil
}

func (u *tenderUsecase) GetTenders(db *gorm.DB, status string) ([]domain.Tender, error) {
	return u.tenderRepo.GetTenders(db, status)
}

func (u *tenderUsecase) GetTenderDetail(db *gorm.DB, id uuid.UUID) (*domain.Tender, error) {
	return u.tenderRepo.GetTenderByID(db, id)
}

func (u *tenderUsecase) SubmitBid(db *gorm.DB, idempotencyKey string, companyID uuid.UUID, tenderID uuid.UUID, pricePerKG int64) (*domain.Bid, error) {
	var createdBid *domain.Bid

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Pessimistic lock tender to prevent bidding on closed auctions
		tender, err := u.tenderRepo.GetTenderByIDForUpdate(tx, tenderID)
		if err != nil {
			return errors.New("tender not found")
		}

		if tender.Status != domain.TenderBidding {
			return errors.New("tender is not accepting bids")
		}

		if time.Now().After(tender.BiddingDeadline) {
			return errors.New("bidding deadline has passed")
		}

		if pricePerKG < tender.ReservePricePerKG {
			return fmt.Errorf("bid price per kg (%d) is below reserve price (%d)", pricePerKG, tender.ReservePricePerKG)
		}

		// Calculate total offer amount
		totalOffer := int64(tender.EstimatedWeightKG * float64(pricePerKG))

		// 2. Lock Bid-Bond via Double-Entry Ledger
		// Debit 1010 (Platform Custody Bank) & Credit 2010 (Bid-Bond Payable Lapak)
		bankAcc, err := u.ledgerRepo.GetAccountByCode(tx, "1010")
		if err != nil {
			return err
		}
		bidBondAcc, err := u.ledgerRepo.GetAccountByCode(tx, "2010")
		if err != nil {
			return err
		}

		entries := []domain.LedgerEntry{
			{AccountID: bankAcc.ID, EntryType: domain.EntryDebit, Amount: tender.BidBondAmount},
			{AccountID: bidBondAcc.ID, EntryType: domain.EntryCredit, Amount: tender.BidBondAmount},
		}

		_, err = u.ledgerRepo.PostJournal(
			tx,
			idempotencyKey,
			domain.TxTypeBidBondDeposit,
			tender.TenderNumber,
			fmt.Sprintf("Bid bond locked for Tender %s", tender.TenderNumber),
			entries,
		)
		if err != nil {
			return fmt.Errorf("failed to lock bid bond in ledger: %w", err)
		}

		// 3. Create Bid record
		now := time.Now()
		bid := domain.Bid{
			ID:               uuid.New(),
			TenderID:         tenderID,
			BidderCompanyID:  companyID,
			PricePerKG:       pricePerKG,
			TotalOfferAmount: totalOffer,
			BidBondLocked:    true,
			Status:           domain.BidSubmitted,
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		if err := u.tenderRepo.CreateBid(tx, &bid); err != nil {
			return err
		}

		createdBid = &bid
		return nil
	})

	if err != nil {
		return nil, err
	}
	return createdBid, nil
}

func (u *tenderUsecase) AwardTender(db *gorm.DB, idempotencyKey string, tenderID uuid.UUID, winningBidID uuid.UUID) (*domain.Contract, error) {
	var contract *domain.Contract

	err := db.Transaction(func(tx *gorm.DB) error {
		tender, err := u.tenderRepo.GetTenderByIDForUpdate(tx, tenderID)
		if err != nil {
			return errors.New("tender not found")
		}

		if tender.Status != domain.TenderBidding {
			return errors.New("tender cannot be awarded in current status")
		}

		winningBid, err := u.tenderRepo.GetBidByID(tx, winningBidID)
		if err != nil {
			return errors.New("winning bid not found")
		}

		now := time.Now()

		// 1. Mark Tender as Awarded
		tender.Status = domain.TenderAwarded
		if err := u.tenderRepo.UpdateTender(tx, tender); err != nil {
			return err
		}

		// 2. Mark winning & losing bids
		allBids, _ := u.tenderRepo.GetBidsByTenderID(tx, tenderID)
		for _, b := range allBids {
			if b.ID == winningBidID {
				b.Status = domain.BidWon
			} else {
				b.Status = domain.BidLost
				// Auto-Refund Bid Bond for losers via Double-Entry Ledger
				bankAcc, _ := u.ledgerRepo.GetAccountByCode(tx, "1010")
				bidBondAcc, _ := u.ledgerRepo.GetAccountByCode(tx, "2010")
				entries := []domain.LedgerEntry{
					{AccountID: bidBondAcc.ID, EntryType: domain.EntryDebit, Amount: tender.BidBondAmount},
					{AccountID: bankAcc.ID, EntryType: domain.EntryCredit, Amount: tender.BidBondAmount},
				}
				refundKey := fmt.Sprintf("REFUND-%s-%s", tender.TenderNumber, b.ID.String()[:8])
				_, _ = u.ledgerRepo.PostJournal(tx, refundKey, domain.TxTypeBidBondRefund, tender.TenderNumber, "Bid bond refund for losing bidder", entries)
			}
			_ = u.tenderRepo.UpdateBid(tx, &b)
		}

		// 3. Generate Official SPK Contract
		newContract := domain.Contract{
			ID:                   uuid.New(),
			TenderID:             tenderID,
			SellerCompanyID:      tender.PublisherCompanyID,
			BuyerCompanyID:       winningBid.BidderCompanyID,
			ContractNumber:       fmt.Sprintf("SPK-%s-%d", now.Format("20060102"), now.UnixNano()%10000),
			AgreedPricePerKG:     winningBid.PricePerKG,
			EstimatedTotalAmount: winningBid.TotalOfferAmount,
			EscrowLockedAmount:   tender.BidBondAmount, // Seed initial escrow from bid bond
			Status:               domain.ContractActive,
			StartDate:            now,
			EndDate:              now.Add(30 * 24 * time.Hour),
			CreatedAt:            now,
			UpdatedAt:            now,
		}

		if err := u.tenderRepo.CreateContract(tx, &newContract); err != nil {
			return err
		}

		contract = &newContract
		return nil
	})

	if err != nil {
		return nil, err
	}
	return contract, nil
}
