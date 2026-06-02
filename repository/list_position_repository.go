package repository

import (
	"github.com/google/uuid"
	"github.com/mrkeylost/Flowboard_Backend/config"
	"github.com/mrkeylost/Flowboard_Backend/model"
)

type ListPositionRepository interface {
	GetByBoard(boardPublicID string) (*model.ListPosition, error)
	CreateOrUpdate(boardPublicID string, listOrder []uuid.UUID) error
	GetListOrder(boardPublicID string) ([]uuid.UUID, error)
	UpdateListOrder(position *model.ListPosition) error
}

type listPositionRepository struct{}

func NewListPositionRepository() ListPositionRepository {
	return &listPositionRepository{}
}

func (repo *listPositionRepository) GetByBoard(boardPublicID string) (*model.ListPosition, error) {
	var position model.ListPosition

	err := config.DBConn.Joins("JOIN list_positions lp ON lp.board_internal_id = boards.internal_id").
		Where("boards.public_id = ?", boardPublicID).
		Error

	return &position, err
}

func (repo *listPositionRepository) CreateOrUpdate(boardPublicID string, listOrder []uuid.UUID) error {
	return config.DBConn.Exec(`
		INSERT INTO list_positions (board_internal_id, list_order)
		SELECT internal_id, ? FROM boards WHERE public_id = ?
		ON CONFLICT (board_internal_id)
		DO UPDATE SET list_order = EXCLUDE.list_order
	`, listOrder, boardPublicID).Error
}

func (repo *listPositionRepository) GetListOrder(boardPublicID string) ([]uuid.UUID, error) {
	findPosition, err := repo.GetByBoard(boardPublicID)
	if err != nil {
		return nil, err
	}

	return findPosition.ListOrder, err
}

func (repo *listPositionRepository) UpdateListOrder(position *model.ListPosition) error {
	return config.DBConn.Model(position).
		Where("internal_id = ?", position.InternalID).
		Update("list_order", position.ListOrder).
		Error
}
