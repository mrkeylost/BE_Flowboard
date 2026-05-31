package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mrkeylost/Flowboard_Backend/model"
	"github.com/mrkeylost/Flowboard_Backend/repository"
)

type BoardService interface {
	CreateBoard(board *model.Board) error
	UpdateBoard(board *model.Board) error
	FindBoardByPublicID(publicID string) (*model.Board, error)
	AddBoardMember(boardPublicID string, userPublicIDs []string) error
}

type boardService struct {
	boardRepo       repository.BoardRepository
	userRepo        repository.AuthRepository
	boardMemberRepo repository.BoardMemberRepository
}

func NewBoardService(boardRepo repository.BoardRepository, userRepo repository.AuthRepository, boardMemberRepo repository.BoardMemberRepository) BoardService {
	return &boardService{boardRepo, userRepo, boardMemberRepo}
}

func (service *boardService) CreateBoard(board *model.Board) error {
	owner, err := service.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("Owner not found")
	}

	board.PublicID = uuid.New()
	board.OwnerID = owner.InternalID

	return service.boardRepo.Create(board)
}

func (service *boardService) UpdateBoard(board *model.Board) error {
	return service.boardRepo.Update(board)
}

func (service *boardService) FindBoardByPublicID(publicID string) (*model.Board, error) {
	return service.boardRepo.FindByPublicID(publicID)
}

func (service *boardService) AddBoardMember(boardPublicID string, userPublicIDs []string) error {
	findBoard, err := service.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("Board not found")
	}

	var userInternalIDs []uint
	for _, value := range userPublicIDs {
		user, err := service.userRepo.FindByPublicID(value)
		if err != nil {
			return errors.New("User not found: " + value)
		}

		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	memberList, err := service.boardMemberRepo.GetMember(string(findBoard.PublicID.String()))
	if err != nil {
		return err
	}

	memberMap := make(map[uint]bool)
	for _, member := range memberList {
		memberMap[uint(member.InternalID)] = true
	}

	var newMembersID []uint
	for _, userId := range userInternalIDs {
		if !memberMap[userId] {
			newMembersID = append(newMembersID, userId)
		}
	}

	if len(newMembersID) == 0 {
		return nil
	}

	return service.boardRepo.AddMember(uint(findBoard.InternalID), newMembersID)
}
