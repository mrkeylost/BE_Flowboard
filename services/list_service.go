package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mrkeylost/Flowboard_Backend/config"
	"github.com/mrkeylost/Flowboard_Backend/model"
	"github.com/mrkeylost/Flowboard_Backend/model/types"
	"github.com/mrkeylost/Flowboard_Backend/repository"
	"github.com/mrkeylost/Flowboard_Backend/utils"
	"gorm.io/gorm"
)

type ListWithOrder struct {
	Positions []uuid.UUID
	Lists     []model.List
}

type ListService interface {
	CreateList(list *model.List) error
	UpdateList(list *model.List) error
	DeleteList(id uint) error
	GetListByBoard(boardPublicID string) (*ListWithOrder, error)
	GetListByID(id uint) (*model.List, error)
	GetListByPublicID(publicID string) (*model.List, error)
	UpdatePositions(boardPublicID string, positions []uuid.UUID) error
}

type listService struct {
	listRepo    repository.ListRepository
	boardRepo   repository.BoardRepository
	listPosRepo repository.ListPositionRepository
}

func NewListService(listRepo repository.ListRepository, boardRepo repository.BoardRepository, listPosRepo repository.ListPositionRepository) ListService {
	return &listService{listRepo, boardRepo, listPosRepo}
}

func (service *listService) CreateList(list *model.List) error {
	findBoard, err := service.boardRepo.FindByPublicID(list.BoardPublicID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Board not found")
		}

		return fmt.Errorf("Failed to get board : %w", err)
	}

	list.BoardInternalID = findBoard.InternalID

	if list.PublicID == uuid.Nil {
		list.PublicID = uuid.New()
	}

	transaction := config.DBConn.Begin()

	recover := func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}
	defer recover()

	if err := transaction.Create(list).Error; err != nil {
		transaction.Rollback()
		return fmt.Errorf("failed to create list: %w", err)
	}

	var positions model.ListPosition

	res := transaction.Where("board_internal_id = ?", findBoard.InternalID).First(&positions)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		positions = model.ListPosition{
			PublicID:  uuid.New(),
			BoardID:   findBoard.InternalID,
			ListOrder: types.UUIDArray{list.PublicID},
		}

		if err := transaction.Create(&positions).Error; err != nil {
			transaction.Rollback()
			return fmt.Errorf("Failed to create list position: %w", err)
		}
	} else if res.Error != nil {
		transaction.Rollback()
		return fmt.Errorf("Failed to create list position: %w", res.Error)
	} else {
		positions.ListOrder = append(positions.ListOrder, list.PublicID)

		if err := transaction.Model(&positions).Update("list_order", positions.ListOrder); err.Error != nil {
			transaction.Rollback()
			return fmt.Errorf("Failed to update list position: %w", err.Error)
		}
	}

	if err := transaction.Commit().Error; err != nil {
		return fmt.Errorf("Transaction commit failed: %w", err)
	}

	return nil
}

func (service *listService) UpdateList(list *model.List) error {
	return service.listRepo.Update(list)
}

func (service *listService) DeleteList(id uint) error {
	return service.listRepo.Delete(id)
}

func (service *listService) GetListByBoard(boardPublicID string) (*ListWithOrder, error) {
	_, err := service.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return nil, errors.New("Board not found")
	}

	positions, err := service.listPosRepo.GetListOrder(boardPublicID)
	if err != nil {
		return nil, errors.New("List order not found: " + err.Error())
	}

	findList, err := service.listRepo.FindByBoard(boardPublicID)
	if err != nil {
		return nil, errors.New("List not found: " + err.Error())
	}

	orderedList := utils.SortListByPos(findList, positions)

	return &ListWithOrder{
		Positions: positions,
		Lists:     orderedList,
	}, nil
}

func (service *listService) GetListByID(id uint) (*model.List, error) {
	return service.listRepo.FindByID(id)
}

func (service *listService) GetListByPublicID(publicID string) (*model.List, error) {
	return service.listRepo.FindByPublicID(publicID)
}

func (service *listService) UpdatePositions(boardPublicID string, positions []uuid.UUID) error {
	findBoard, err := service.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("Board not found")
	}

	findPosition, err := service.listPosRepo.GetByBoard(findBoard.PublicID.String())
	if err != nil {
		return errors.New("List position not found")
	}

	findPosition.ListOrder = positions

	return service.listPosRepo.UpdateListOrder(findPosition)
}
