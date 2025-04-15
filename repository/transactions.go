package repository

import (
	"github.com/DecodeWorms/BorderBlitz/models"
	"gorm.io/gorm"
)

// TransactionRepository handles database operations for transactions
type TransactionRepository struct {
	Db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{Db: db}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(tx *models.Transaction) error {
	return r.Db.Create(tx).Error
}

// FindByID finds a transaction by ID
func (r *TransactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.Db.Preload("SourceCoin").
		Preload("DestinationCoin").
		Preload("SenderWallet").
		Preload("ReceiverWallet").
		First(&tx, id).Error
	return &tx, err
}

// FindByReference finds a transaction by reference
func (r *TransactionRepository) FindByReference(reference string) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.Db.Where("reference = ?", reference).First(&tx).Error
	return &tx, err
}

// FindByWalletID finds all transactions related to a wallet
func (r *TransactionRepository) FindByWalletID(walletID uint, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.Db.Where("sender_wallet_id = ? OR receiver_wallet_id = ?", walletID, walletID).
		Preload("SourceCoin").
		Preload("DestinationCoin").
		Preload("SenderWallet").
		Preload("ReceiverWallet").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

// CountByWalletID counts all transactions related to a wallet
func (r *TransactionRepository) CountByWalletID(walletID uint) (int64, error) {
	var count int64
	err := r.Db.Model(&models.Transaction{}).
		Where("sender_wallet_id = ? OR receiver_wallet_id = ?", walletID, walletID).
		Count(&count).Error
	return count, err
}
