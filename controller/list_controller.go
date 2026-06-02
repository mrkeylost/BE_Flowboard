package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mrkeylost/Flowboard_Backend/model"
	"github.com/mrkeylost/Flowboard_Backend/services"
	"github.com/mrkeylost/Flowboard_Backend/utils"
)

type ListController struct {
	service services.ListService
}

func NewListController(service services.ListService) *ListController {
	return &ListController{service}
}

func (controller *ListController) CreateList(ctx fiber.Ctx) error {
	list := new(model.List)

	if err := ctx.Bind().Body(list); err != nil {
		return utils.BadRequest(ctx, "Parsing data failed", err.Error())
	}

	if err := controller.service.CreateList(list); err != nil {
		return utils.BadRequest(ctx, "Create list failed", err.Error())
	}

	return utils.Success(ctx, "Create list success", list)
}

func (controller *ListController) UpdateList(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(model.List)

	if err := ctx.Bind().Body(list); err != nil {
		return utils.BadRequest(ctx, "Parsing data failed", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Invalid ID", err.Error())
	}

	findList, err := controller.service.GetListByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}

	list.InternalID = findList.InternalID
	list.PublicID = findList.PublicID

	if err := controller.service.UpdateList(list); err != nil {
		return utils.BadRequest(ctx, "Update list failed", err.Error())
	}

	updatedList, err := controller.service.GetListByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}

	return utils.Success(ctx, "Update list success", updatedList)
}

func (controller *ListController) DeleteList(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Invalid ID", err.Error())
	}

	findList, err := controller.service.GetListByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}

	if err := controller.service.DeleteList(uint(findList.InternalID)); err != nil {
		return utils.InternalServerError(ctx, "Delete list failed", err.Error())
	}

	return utils.Success(ctx, "Delete list success", publicID)
}

func (controller *ListController) GetListOnBoard(ctx fiber.Ctx) error {
	boardPublicID := ctx.Params("id")

	if _, err := uuid.Parse(boardPublicID); err != nil {
		return utils.BadRequest(ctx, "Invalid ID", err.Error())
	}

	orderedList, err := controller.service.GetListByBoard(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}

	return utils.Success(ctx, "Get list success", orderedList)
}
