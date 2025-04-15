package services

import (
	"errors"
	"fmt"

	"github.com/DecodeWorms/BorderBlitz/models"
	"github.com/DecodeWorms/BorderBlitz/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SwapService handles business logic for swapping stablecoins
type SwapService struct {
	db             *gorm.DB
	walletRepo     *repository.WalletRepository
	stablecoinRepo *repository.StablecoinRepository
	txRepo         *repository.TransactionRepository
}

// NewSwapService creates a new swap service
func NewSwapService(
	db *gorm.DB,
	walletRepo *repository.WalletRepository,
	stablecoinRepo *repository.StablecoinRepository,
	txRepo *repository.TransactionRepository,
) *SwapService {
	return &SwapService{
		db:             db,
		walletRepo:     walletRepo,
		stablecoinRepo: stablecoinRepo,
		txRepo:         txRepo,
	}
}

// Swap swaps one stablecoin for another within the same wallet
func (s *SwapService) Swap(
	walletID uint,
	sourceCoinID, destCoinID uint,
	amount float64,
) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("swap amount must be positive")
	}

	if sourceCoinID == destCoinID {
		return nil, errors.New("source and destination coins must be different")
	}

	var transaction *models.Transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		walletRepoTx := &repository.WalletRepository{Db: tx}
		stablecoinRepoTx := &repository.StablecoinRepository{Db: tx}
		txRepoTx := &repository.TransactionRepository{Db: tx}

		// Get source balance
		sourceBalance, err := walletRepoTx.GetBalance(walletID, sourceCoinID)
		if err != nil {
			return err
		}

		// Check if wallet has enough funds
		if sourceBalance.Amount < amount {
			return errors.New("insufficient funds")
		}

		// Get destination balance
		destBalance, err := walletRepoTx.GetBalance(walletID, destCoinID)
		if err != nil {
			return err
		}

		// Get coin details
		sourceCoin, err := stablecoinRepoTx.FindByID(sourceCoinID)
		if err != nil {
			return fmt.Errorf("source coin not found: %w", err)
		}

		destCoin, err := stablecoinRepoTx.FindByID(destCoinID)
		if err != nil {
			return fmt.Errorf("destination coin not found: %w", err)
		}

		// Calculate exchange rate and destination amount
		exchangeRate := sourceCoin.USDRate / destCoin.USDRate
		destAmount := amount * exchangeRate

		// Update balances
		sourceBalance.Amount -= amount
		if err := walletRepoTx.UpdateBalance(sourceBalance); err != nil {
			return fmt.Errorf("failed to update source balance: %w", err)
		}

		destBalance.Amount += destAmount
		if err := walletRepoTx.UpdateBalance(destBalance); err != nil {
			return fmt.Errorf("failed to update destination balance: %w", err)
		}

		// Create transaction record
		txRef := uuid.New().String()
		transaction = &models.Transaction{
			Type:              models.TransactionSwap,
			SenderWalletID:    &walletID,
			ReceiverWalletID:  &walletID,
			SourceCoinID:      sourceCoinID,
			SourceAmount:      amount,
			DestinationCoinID: destCoinID,
			DestinationAmount: destAmount,
			ExchangeRate:      exchangeRate,
			Reference:         txRef,
			Status:            "completed",
		}

		if err := txRepoTx.Create(transaction); err != nil {
			return fmt.Errorf("failed to create transaction record: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetExchangeRate gets the exchange rate between two stablecoins
func (s *SwapService) GetExchangeRate(sourceCoinID, destCoinID uint) (float64, error) {
	sourceCoin, err := s.stablecoinRepo.FindByID(sourceCoinID)
	if err != nil {
		return 0, fmt.Errorf("source coin not found: %w", err)
	}

	destCoin, err := s.stablecoinRepo.FindByID(destCoinID)
	if err != nil {
		return 0, fmt.Errorf("destination coin not found: %w", err)
	}

	return sourceCoin.USDRate / destCoin.USDRate, nil
}
