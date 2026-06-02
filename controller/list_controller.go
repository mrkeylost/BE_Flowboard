package controller

import (
	"github.com/gofiber/fiber/v3"
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
