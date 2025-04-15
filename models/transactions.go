package models

import (
	"time"

	"gorm.io/gorm"
)

// TransactionType defines the type of transaction
type TransactionType string

const (
	TransactionDeposit    TransactionType = "deposit"
	TransactionWithdrawal TransactionType = "withdrawal"
	TransactionTransfer   TransactionType = "transfer"
	TransactionSwap       TransactionType = "swap"
)

// Transaction represents a financial transaction in the system
type Transaction struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	Type              TransactionType `gorm:"not null" json:"type"`
	SenderWalletID    *uint           `gorm:"index" json:"sender_wallet_id,omitempty"`
	SenderWallet      *Wallet         `json:"sender_wallet,omitempty"`
	ReceiverWalletID  *uint           `gorm:"index" json:"receiver_wallet_id,omitempty"`
	ReceiverWallet    *Wallet         `json:"receiver_wallet,omitempty"`
	SourceCoinID      uint            `gorm:"index;not null" json:"source_coin_id"`
	SourceCoin        Stablecoin      `json:"source_coin"`
	SourceAmount      float64         `gorm:"not null" json:"source_amount"`
	DestinationCoinID uint            `gorm:"index;not null" json:"destination_coin_id"`
	DestinationCoin   Stablecoin      `json:"destination_coin"`
	DestinationAmount float64         `gorm:"not null" json:"destination_amount"`
	ExchangeRate      float64         `gorm:"not null" json:"exchange_rate"`
	Reference         string          `gorm:"not null;uniqueIndex" json:"reference"`
	Status            string          `gorm:"not null;default:'completed'" json:"status"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
}
