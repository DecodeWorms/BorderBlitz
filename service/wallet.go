package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DecodeWorms/BorderBlitz/models"
	"github.com/DecodeWorms/BorderBlitz/repository"
)

// WalletService handles business logic for wallets
type WalletService struct {
	walletRepo     repository.Wallet
	stableCoinRepo repository.StableCoin
	user           repository.User
}

// NewWalletService creates a new wallet service
func NewWalletService(walletRepo repository.Wallet, stablecoinRepo repository.StableCoin, user repository.User) *WalletService {
	return &WalletService{
		walletRepo:     walletRepo,
		stableCoinRepo: stablecoinRepo,
		user:           user,
	}
}

// CreateWallet creates a new wallet for a user
func (s *WalletService) CreateWallet(userID, email, userType string) (*models.Wallet, error) {
	// Validate user ID
	user, err := s.user.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("user's record is not found")
	}

	// Validate user type
	userType = strings.ToLower(userType)
	if userType != "email" && userType != "mobile" {
		return nil, errors.New("user type must be 'email' or 'mobile'")
	}

	// Check if wallet already exists
	existingWallet, err := s.walletRepo.FindByUserID(userID, userType)
	if err == nil {
		return existingWallet, fmt.Errorf("wallet already exists for this %s", userType)
	}

	// Create new wallet
	wallet := &models.Wallet{
		UserID:   user.ID,
		UserType: userType,
	}

	if err := s.walletRepo.Create(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetWallet gets a wallet by ID
func (s *WalletService) GetWallet(id uint) (*models.Wallet, error) {
	wallet, err := s.walletRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("error wallet not available")
	}
	return wallet, nil
}

// GetWalletByUserID gets a wallet by user ID and type
func (s *WalletService) GetWalletByUserID(userID string, userType string) (*models.Wallet, error) {
	wallet, err := s.walletRepo.FindByUserID(userID, userType)
	if err != nil {
		return nil, errors.New("error wallet not available")
	}
	return wallet, nil
}

// GetBalance gets the balance of a specific stablecoin in a wallet
func (s *WalletService) GetBalance(walletID uint, stablecoinID uint) (*models.Balance, error) {
	return s.walletRepo.GetBalance(walletID, stablecoinID)
}

// DepositFunds adds funds to a wallet
func (s *WalletService) DepositFunds(walletID uint, stablecoinID uint, amount float64) (*models.Balance, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be positive")
	}

	// Get current balance
	balance, err := s.walletRepo.GetBalance(walletID, stablecoinID)
	if err != nil {
		return nil, err
	}

	// Update balance
	balance.Amount += amount
	if err := s.walletRepo.UpdateBalance(balance); err != nil {
		return nil, err
	}

	return balance, nil
}

// ListAllWallets lists all wallets
func (s *WalletService) ListAllWallets() ([]models.Wallet, error) {
	return s.walletRepo.ListAll()
}
