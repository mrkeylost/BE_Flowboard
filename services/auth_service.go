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
	UpdateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByPublicID(publicId string) (*model.User, error)
	GetAllUser(search, sort string, limit, offset int) ([]model.User, int64, error)
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

func (service *authService) UpdateUser(user *model.User) error {
	return service.repo.Update(user)
}

func (service *authService) GetUserByID(id uint) (*model.User, error) {
	return service.repo.FindByID(id)
}

func (service *authService) GetUserByPublicID(publicId string) (*model.User, error) {
	return service.repo.FindByPublicID(publicId)
}

func (service *authService) GetAllUser(search, sort string, limit, offset int) ([]model.User, int64, error) {
	return service.repo.FindAll(search, sort, limit, offset)
}
