package repository

import (
	"github.com/mrkeylost/Flowboard_Backend/config"
	"github.com/mrkeylost/Flowboard_Backend/model"
)

type AuthRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type authRepository struct {
}

func NewUserRepository() AuthRepository {
	return &authRepository{}
}

func (repo *authRepository) Create(user *model.User) error {
	return config.DBConn.Create(user).Error
}

func (repo *authRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := config.DBConn.Where("email = ?", email).First(&user).Error

	return &user, err
}
