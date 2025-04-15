package handlers

import (
	"net/http"

	"github.com/DecodeWorms/BorderBlitz/models"
	services "github.com/DecodeWorms/BorderBlitz/service"
	"github.com/gin-gonic/gin"
)

type Users struct {
	user *services.Users
}

func NewUserHandler(user *services.Users) *Users {
	return &Users{
		user: user,
	}
}

func (u *Users) CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.CreateUserRequest
		if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		us, err := u.user.Create(&user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, us)

	}
}
