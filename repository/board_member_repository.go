package repository

import (
	"github.com/mrkeylost/Flowboard_Backend/config"
	"github.com/mrkeylost/Flowboard_Backend/model"
)

type BoardMemberRepository interface {
	GetMember(boardPublicID string) ([]model.User, error)
}

type boardMemberRepository struct{}

func NewBoardMemberRepository() BoardMemberRepository {
	return &boardMemberRepository{}
}

func (repo *boardMemberRepository) GetMember(boardPublicID string) ([]model.User, error) {
	var users []model.User

	err := config.DBConn.
		Joins("JOIN board_members bm ON bm.user_internal_id = users.internal_id").
		Joins("JOIN boards b ON b.internal_id = bm.board_internal_id").
		Where("b.public_id = ?", boardPublicID).
		Find(&users).
		Error

	return users, err
}
