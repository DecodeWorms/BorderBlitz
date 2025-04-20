package repository

import (
	"github.com/DecodeWorms/BorderBlitz/models"
	"gorm.io/gorm"
)

type Users struct {
	Db *gorm.DB
}

func NewUsers(db *gorm.DB) *Users {
	return &Users{
		Db: db,
	}
}

// Create creates new user's record
func (u *Users) Create(data *models.Users) error {
	return u.Db.Create(data).Error
}

func (u *Users) FindByUserID(userID string) (*models.Users, error) {
	var user *models.Users
	err := u.Db.Where("id = ?", userID).First(&user).Error
	return user, err
}

func (u *Users) FindByPhoneNumber(pNumber string) (*models.Users, error) {
	var user *models.Users
	err := u.Db.Where("phone_number=?", pNumber).First(&user).Error
	return user, err
}

func (u *Users) FindByEmail(email string) (*models.Users, error) {
	var user *models.Users
	err := u.Db.Where("email=?", email).First(&user).Error
	return user, err
}

var _ User = &Users{}
