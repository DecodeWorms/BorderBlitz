package repository

import "github.com/DecodeWorms/BorderBlitz/models"

//go:generate mockgen -source=datastore.go -destination=../mocks/datastore_mock.go -package=mocks
type User interface {
	Create(data *models.Users) error
	FindByUserID(userID string) (*models.Users, error)
	FindByPhoneNumber(pNumber string) (*models.Users, error)
	FindByEmail(email string) (*models.Users, error)
}

type Wallet interface {
	Create(wallet *models.Wallet) error
	FindByID(id uint) (*models.Wallet, error)
	FindByUserID(userID string, userType string) (*models.Wallet, error)
	GetBalance(walletID uint, stablecoinID uint) (*models.Balance, error)
	UpdateBalance(balance *models.Balance) error
	ListAll() ([]models.Wallet, error)
}

type Transaction interface {
	Create(tx *models.Transaction) error
	FindByID(id uint) (*models.Transaction, error)
	FindByReference(reference string) (*models.Transaction, error)
	FindByWalletID(walletID uint, limit, offset int) ([]models.Transaction, error)
	CountByWalletID(walletID uint) (int64, error)
}

type StableCoin interface {
	FindByID(id uint) (*models.Stablecoin, error)
	FindBySymbol(symbol string) (*models.Stablecoin, error)
	ListAll() ([]models.Stablecoin, error)
	UpdateRate(id uint, rate float64) error
}
