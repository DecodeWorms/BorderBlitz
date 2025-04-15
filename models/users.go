package models

type Users struct {
	ID       string `gorm:"primaryKey" json:"id"`
	UserType string `gorm:"user_type" json:"user_type"`
}
