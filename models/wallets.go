package models

import (
	"time"

	"gorm.io/gorm"
)

// Wallet represents a user's digital wallet containing multiple stablecoins
type Wallet struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"index;not null" json:"user_id"` // Can be email or phone
	UserType  string         `gorm:"not null" json:"user_type"`     // "email" or "mobile"
	Balances  []Balance      `gorm:"foreignKey:WalletID" json:"balances"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Balance represents the amount of a specific stablecoin in a wallet
type Balance struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	WalletID     uint           `gorm:"index;not null" json:"wallet_id"`
	StablecoinID uint           `gorm:"index;not null" json:"stablecoin_id"`
	Stablecoin   Stablecoin     `json:"stablecoin"`
	Amount       float64        `gorm:"not null;default:0" json:"amount"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// GetTotalUSDValue calculates the total value of all balances in USD
func (w *Wallet) GetTotalUSDValue() float64 {
	var total float64
	for _, balance := range w.Balances {
		total += balance.Stablecoin.GetUSDValue(balance.Amount)
	}
	return total
}
