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
	AddMember(boardID uint, userIDs []uint) error
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

func (repo *boardRepository) AddMember(boardID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}

	currentTime := time.Now()
	var members []model.BoardMember

	for _, value := range userIDs {
		members = append(members, model.BoardMember{
			BoardID:  int64(boardID),
			UserID:   int64(value),
			JoinedAt: currentTime,
		})
	}

	return config.DBConn.Create(&members).Error
}
