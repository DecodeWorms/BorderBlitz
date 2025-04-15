package models

import (
	"time"

	"gorm.io/gorm"
)

// Stablecoin represents a digital currency pegged to a fiat currency
type Stablecoin struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Symbol    string         `gorm:"unique;not null" json:"symbol"`
	Name      string         `gorm:"not null" json:"name"`
	USDRate   float64        `gorm:"not null" json:"usd_rate"` // Exchange rate: 1 unit of this coin = X USD
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// GetUSDValue converts an amount of this stablecoin to USD
func (s *Stablecoin) GetUSDValue(amount float64) float64 {
	return amount * s.USDRate
}

// ConvertTo converts an amount of this stablecoin to another stablecoin
func (s *Stablecoin) ConvertTo(amount float64, target *Stablecoin) float64 {
	// First convert to USD, then to target currency
	usdValue := s.GetUSDValue(amount)
	return usdValue / target.USDRate
}
