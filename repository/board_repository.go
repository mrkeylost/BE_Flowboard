package repository

import (
	"time"

	"github.com/mrkeylost/Flowboard_Backend/config"
	"github.com/mrkeylost/Flowboard_Backend/model"
)

type BoardRepository interface {
	Create(board *model.Board) error
	Update(board *model.Board) error
	FindByPublicID(publicID string) (*model.Board, error)
	FindAllByUser(userPublicID, search, sort string, limit, offset int) ([]model.Board, int64, error)
	AddMember(boardID uint, userIDs []uint) error
	RemoveMember(boardID uint, userIDs []uint) error
}

type boardRepository struct{}

func NewBoardRepository() BoardRepository {
	return &boardRepository{}
}

func (repo *boardRepository) Create(board *model.Board) error {
	return config.DBConn.Create(board).Error
}

func (repo *boardRepository) Update(board *model.Board) error {
	return config.DBConn.Model(&model.Board{}).Where("public_id = ?", board.PublicID).Updates(map[string]interface{}{
		"title":       board.Title,
		"description": board.Description,
		"due_date":    board.DueDate,
	}).Error
}

func (repo *boardRepository) FindByPublicID(publicID string) (*model.Board, error) {
	var board model.Board

	err := config.DBConn.Where("public_id = ?", publicID).First(&board).Error

	return &board, err
}

func (repo *boardRepository) FindAllByUser(userPublicID, search, sort string, limit, offset int) ([]model.Board, int64, error) {
	var board []model.Board
	var total int64

	query := config.DBConn.Model(&model.Board{}).
		Where("owner_public_id = ? OR internal_id IN ("+
			"SELECT bm.board_internal_id FROM board_members bm "+
			"JOIN users u ON u.internal_id = bm.user_internal_id "+
			"WHERE u.public_id = ?)", userPublicID, userPublicID)

	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Limit(limit).Offset(offset).Find(&board).Error; err != nil {
		return nil, 0, err
	}

	return board, total, nil
}

func (repo *boardRepository) AddMember(boardID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}

	var members []model.BoardMember

	for _, value := range userIDs {
		members = append(members, model.BoardMember{
			BoardID:  int64(boardID),
			UserID:   int64(value),
			JoinedAt: time.Now(),
		})
	}

	return config.DBConn.Create(&members).Error
}

func (repo *boardRepository) RemoveMember(boardID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}

	return config.DBConn.
		Where("board_internal_id = ? AND user_internal_id IN (?)", boardID, userIDs).
		Delete(&model.BoardMember{}).
		Error

}
