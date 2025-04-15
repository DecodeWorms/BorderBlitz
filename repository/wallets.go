package repository

import (
	"errors"

	"github.com/DecodeWorms/BorderBlitz/models"
	"gorm.io/gorm"
)

// WalletRepository handles database operations for wallets
type WalletRepository struct {
	Db *gorm.DB
}

// NewWalletRepository creates a new wallet repository
func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{Db: db}
}

// Create creates a new wallet
func (r *WalletRepository) Create(wallet *models.Wallet) error {
	return r.Db.Create(wallet).Error
}

// FindByID finds a wallet by ID
func (r *WalletRepository) FindByID(id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.Db.Preload("Balances.Stablecoin").First(&wallet, id).Error
	return &wallet, err
}

// FindByUserID finds a wallet by user ID
func (r *WalletRepository) FindByUserID(userID string, userType string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.Db.Where("user_id = ? AND user_type = ?", userID, userType).
		Preload("Balances.Stablecoin").
		First(&wallet).Error

	return &wallet, err
}

// GetBalance gets the balance of a specific stablecoin in a wallet
func (r *WalletRepository) GetBalance(walletID uint, stablecoinID uint) (*models.Balance, error) {
	var balance models.Balance
	err := r.Db.Where("wallet_id = ? AND stablecoin_id = ?", walletID, stablecoinID).
		Preload("Stablecoin").
		First(&balance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a new balance record with zero amount
			balance = models.Balance{
				WalletID:     walletID,
				StablecoinID: stablecoinID,
				Amount:       0,
			}
			if err := r.Db.Create(&balance).Error; err != nil {
				return nil, err
			}
			return &balance, nil
		}
		return nil, err
	}
	return &balance, nil
}

// UpdateBalance updates the balance of a specific stablecoin in a wallet
func (r *WalletRepository) UpdateBalance(balance *models.Balance) error {
	return r.Db.Save(balance).Error
}

// ListAll lists all wallets
func (r *WalletRepository) ListAll() ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.Db.Preload("Balances.Stablecoin").Find(&wallets).Error
	return wallets, err
}
