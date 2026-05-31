package repository

import (
	"strings"

	"github.com/mrkeylost/Flowboard_Backend/config"
	"github.com/mrkeylost/Flowboard_Backend/model"
)

type AuthRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindByPublicID(publiCId string) (*model.User, error)
	FindAll(search, sort string, limit, offset int) ([]model.User, int64, error)
}

type authRepository struct {
}

func NewUserRepository() AuthRepository {
	return &authRepository{}
}

func (repo *authRepository) Create(user *model.User) error {
	return config.DBConn.Create(user).Error
}

func (repo *authRepository) Update(user *model.User) error {
	return config.DBConn.Model(&model.User{}).Where("public_id = ?", user.PublicID).Updates(map[string]interface{}{
		"name": user.Name,
	}).Error
}

func (repo *authRepository) Delete(id uint) error {
	return config.DBConn.Delete(&model.User{}, id).Error
}

func (repo *authRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := config.DBConn.Where("email = ?", email).First(&user).Error

	return &user, err
}

func (repo *authRepository) FindByID(id uint) (*model.User, error) {
	var user model.User

	err := config.DBConn.First(&user, id).Error

	return &user, err
}

func (repo *authRepository) FindByPublicID(publiCId string) (*model.User, error) {
	var user model.User

	err := config.DBConn.Where("public_id = ?", publiCId).First(&user).Error

	return &user, err
}

func (repo *authRepository) FindAll(search, sort string, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := config.DBConn.Model(&model.User{})

	if search != "" {
		queryPattern := "%" + search + "%"

		// insensitive search
		db = db.Where("name Ilike ? OR email Ilike ?", queryPattern, queryPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sort != "" {
		switch sort {
		case "-id":
			sort = "-internal_id"
		case "id":
			sort = "internal_id"
		}

		if after, ok := strings.CutPrefix(sort, "-"); ok {
			sort = after + " DESC"
		} else {
			sort += " ASC"
		}

		db = db.Order(sort)
	}

	err := db.Limit(limit).Offset(offset).Find(&users).Error

	return users, total, err
}
