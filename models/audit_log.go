package models

import (
	"time"
)

// AuditLog represents a compliance log entry for user actions
type AuditLog struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	UserIP    string    `json:"user_ip" gorm:"type:varchar(50)"`
	UserAgent string    `json:"user_agent" gorm:"type:text"`
	Country   string    `json:"country" gorm:"type:varchar(50)"`
	Browser   string    `json:"browser" gorm:"type:varchar(50)"`
	Action    string    `json:"action" gorm:"type:varchar(100)"`
	WalletID  *uint     `json:"wallet_id"`
	Wallet    *Wallet   `json:"-" gorm:"foreignkey:WalletID"`
	Details   string    `json:"details" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}
