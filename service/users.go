package services

import (
	"github.com/DecodeWorms/BorderBlitz/models"
	"github.com/DecodeWorms/BorderBlitz/repository"
	"github.com/google/uuid"
)

type Users struct {
	user repository.User
}

func NewUsers(u repository.User) *Users {
	return &Users{
		user: u,
	}
}

func (u *Users) Create(data *models.CreateUserRequest) (*models.Users, error) {
	//Create user's record
	rec := &models.Users{
		ID:       uuid.New().String(),
		UserType: data.UserType,
	}

	if err := u.user.Create(rec); err != nil {
		return nil, err
	}
	return rec, nil
}
