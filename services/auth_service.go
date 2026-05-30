package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mrkeylost/Flowboard_Backend/model"
	"github.com/mrkeylost/Flowboard_Backend/repository"
	"github.com/mrkeylost/Flowboard_Backend/utils"
)

type AuthService interface {
	Register(user *model.User) error
	Login(email, password string) (*model.User, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewUserService(repo repository.AuthRepository) AuthService {
	return &authService{repo}
}

func (service *authService) Register(user *model.User) error {
	findUser, _ := service.repo.FindByEmail(user.Email)
	if findUser.InternalID != 0 {
		return errors.New("Email already used")
	}

	hash, err := utils.Encrypt(user.Password)
	if err != nil {
		return err
	}

	user.Password = hash
	user.Role = "user"
	user.PublicID = uuid.New()

	return service.repo.Create(user)
}

func (service *authService) Login(email, password string) (*model.User, error) {
	findUser, err := service.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("Invalid email or password")
	}

	credentialMatch := utils.Compare(password, findUser.Password)
	if !credentialMatch {
		return nil, errors.New("Invalid email or password")
	}

	return findUser, nil
}
