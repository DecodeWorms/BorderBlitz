package repository

import (
	"github.com/DecodeWorms/BorderBlitz/models"
	"gorm.io/gorm"
)

// StablecoinRepository handles database operations for stablecoins
type StablecoinRepository struct {
	Db *gorm.DB
}

// NewStablecoinRepository creates a new stablecoin repository
func NewStablecoinRepository(db *gorm.DB) *StablecoinRepository {
	return &StablecoinRepository{Db: db}
}

// FindByID finds a stablecoin by ID
func (r *StablecoinRepository) FindByID(id uint) (*models.Stablecoin, error) {
	var coin models.Stablecoin
	err := r.Db.First(&coin, id).Error
	return &coin, err
}

// FindBySymbol finds a stablecoin by symbol
func (r *StablecoinRepository) FindBySymbol(symbol string) (*models.Stablecoin, error) {
	var coin models.Stablecoin
	err := r.Db.Where("symbol = ?", symbol).First(&coin).Error
	return &coin, err
}

// ListAll lists all stablecoins
func (r *StablecoinRepository) ListAll() ([]models.Stablecoin, error) {
	var coins []models.Stablecoin
	err := r.Db.Find(&coins).Error
	return coins, err
}

// UpdateRate updates the USD rate of a stablecoin
func (r *StablecoinRepository) UpdateRate(id uint, rate float64) error {
	return r.Db.Model(&models.Stablecoin{}).Where("id = ?", id).Update("usd_rate", rate).Error
}
