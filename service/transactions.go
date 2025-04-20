package services

import (
	"errors"
	"fmt"

	"github.com/DecodeWorms/BorderBlitz/models"
	"github.com/DecodeWorms/BorderBlitz/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionService handles business logic for transactions
type TransactionService struct {
	txRepo         repository.Transaction
	walletRepo     repository.Wallet
	stablecoinRepo repository.StableCoin
	db             *gorm.DB
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	db *gorm.DB,
	txRepo repository.Transaction,
	walletRepo repository.Wallet,
	stablecoinRepo repository.StableCoin,
) *TransactionService {
	return &TransactionService{
		db:             db,
		txRepo:         txRepo,
		walletRepo:     walletRepo,
		stablecoinRepo: stablecoinRepo,
	}
}

// Transfer transfers funds from one wallet to another
func (s *TransactionService) Transfer(
	senderWalletID, receiverWalletID uint,
	sourceCoinID, destCoinID uint,
	amount float64,
) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("transfer amount must be positive")
	}

	var transaction *models.Transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Create temporary repositories that use the transaction context
		walletRepoTx := &repository.WalletRepository{Db: tx}
		stablecoinRepoTx := &repository.StablecoinRepository{Db: tx}
		txRepoTx := &repository.TransactionRepository{Db: tx}

		// Get sender wallet and balance
		senderBalance, err := walletRepoTx.GetBalance(senderWalletID, sourceCoinID)
		if err != nil {
			return err
		}

		// Check if sender has enough funds
		if senderBalance.Amount < amount {
			return errors.New("insufficient funds")
		}

		// Get receiver wallet
		receiverWallet, err := walletRepoTx.FindByID(receiverWalletID)
		if err != nil {
			return fmt.Errorf("receiver wallet not found: %w", err)
		}

		// Get source and destination coins
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

		// Subtract from sender
		senderBalance.Amount -= amount
		if err := walletRepoTx.UpdateBalance(senderBalance); err != nil {
			return fmt.Errorf("failed to update sender balance: %w", err)
		}

		// Add to receiver
		receiverBalance, err := walletRepoTx.GetBalance(receiverWalletID, destCoinID)
		if err != nil {
			return fmt.Errorf("failed to get receiver balance: %w", err)
		}
		receiverBalance.Amount += destAmount
		if err := walletRepoTx.UpdateBalance(receiverBalance); err != nil {
			return fmt.Errorf("failed to update receiver balance: %w", err)
		}

		// Create transaction record
		txRef := uuid.New().String()
		transaction = &models.Transaction{
			Type:              models.TransactionTransfer,
			SenderWalletID:    &senderWalletID,
			ReceiverWalletID:  &receiverWalletID,
			ReceiverWallet:    receiverWallet,
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

// Deposit adds funds to a wallet (simulated)
func (s *TransactionService) Deposit(walletID uint, coinID uint, amount float64) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be positive")
	}

	var transaction *models.Transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		walletRepoTx := &repository.WalletRepository{Db: tx}
		stablecoinRepoTx := &repository.StablecoinRepository{Db: tx}
		txRepoTx := &repository.TransactionRepository{Db: tx}

		// Get wallet and balance
		balance, err := walletRepoTx.GetBalance(walletID, coinID)
		if err != nil {
			return err
		}

		// Get coin details
		coin, err := stablecoinRepoTx.FindByID(coinID)
		if err != nil {
			return fmt.Errorf("coin not found: %w", err)
		}

		// Add to balance
		balance.Amount += amount
		if err := walletRepoTx.UpdateBalance(balance); err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}

		// Create transaction record
		txRef := uuid.New().String()
		transaction = &models.Transaction{
			Type:              models.TransactionDeposit,
			ReceiverWalletID:  &walletID,
			SourceCoinID:      coinID,
			SourceCoin:        *coin,
			SourceAmount:      amount,
			DestinationCoinID: coinID,
			DestinationAmount: amount,
			ExchangeRate:      1.0,
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

// GetTransaction gets a transaction by ID
func (s *TransactionService) GetTransaction(id uint) (*models.Transaction, error) {
	return s.txRepo.FindByID(id)
}

// GetTransactionHistory gets all transactions for a wallet
func (s *TransactionService) GetTransactionHistory(walletID uint, limit, offset int) ([]models.Transaction, int64, error) {
	if limit <= 0 {
		limit = 10
	}

	transactions, err := s.txRepo.FindByWalletID(walletID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.txRepo.CountByWalletID(walletID)
	if err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}
