package repository

import (
	"github.com/google/uuid"
	"github.com/mrkeylost/Flowboard_Backend/config"
	"github.com/mrkeylost/Flowboard_Backend/model"
)

type ListRepository interface {
	Create(list *model.List) error
	Update(list *model.List) error
	Delete(id uint) error
	FindByBoard(boardPublicID string) ([]model.List, error)
	FindByPublicID(publicID string) (*model.List, error)
	FindByID(id uint) (*model.List, error)
	UpdatePosition(boardPublicID string, position []string) error
	GetCardPosition(listPublicID string) ([]uuid.UUID, error)
}

type listRepository struct{}

func NewListRepository() ListRepository {
	return &listRepository{}
}

func (repo *listRepository) Create(list *model.List) error {
	return config.DBConn.Create(list).Error
}

func (repo *listRepository) Update(list *model.List) error {
	return config.DBConn.Model(&model.List{}).Where("public_id = ?", list.PublicID).Updates(map[string]interface{}{
		"title": list.Title,
	}).Error
}

func (repo *listRepository) Delete(id uint) error {
	return config.DBConn.Delete(&model.List{}, id).Error
}

func (repo *listRepository) UpdatePosition(boardPublicID string, position []string) error {
	return config.DBConn.Model(&model.ListPosition{}).
		Where("board_internal_id = (SELECT internal_id FROM boards WHERE public_id = ?)", boardPublicID).
		Update("list_order", position).
		Error
}

func (repo *listRepository) GetCardPosition(listPublicID string) ([]uuid.UUID, error) {
	var position model.CardPosition
	err := config.DBConn.
		Table("lists l").
		Joins("JOIN card_position cp ON cp.list_internal_id = l.internal_id").
		Where("l.public_id = ?", listPublicID).
		Error

	return position.CardOrder, err
}

func (repo *listRepository) FindByBoard(boardPublicID string) ([]model.List, error) {
	var list []model.List
	err := config.DBConn.
		Where("board_public_id = ?", boardPublicID).
		Order("internal_id ASC").
		Find(&list).
		Error

	return list, err
}

func (repo *listRepository) FindByPublicID(publicID string) (*model.List, error) {
	var list model.List
	err := config.DBConn.
		Where("public_id = ?", publicID).
		Order("internal_id ASC").
		Find(&list).
		Error

	return &list, err
}

func (repo *listRepository) FindByID(id uint) (*model.List, error) {
	var list model.List
	err := config.DBConn.
		First(&list, id).
		Error

	return &list, err
}
