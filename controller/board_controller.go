package controller

import (
	"math"
	"strconv"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mrkeylost/Flowboard_Backend/model"
	"github.com/mrkeylost/Flowboard_Backend/services"
	"github.com/mrkeylost/Flowboard_Backend/utils"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(service services.BoardService) *BoardController {
	return &BoardController{service}
}

func (controller *BoardController) CreateBoard(ctx fiber.Ctx) error {
	board := new(model.Board)

	user := jwtware.FromContext(ctx)
	claims := user.Claims.(jwt.MapClaims)

	if err := ctx.Bind().Body(board); err != nil {
		return utils.BadRequest(ctx, "Parsing data failed", err.Error())
	}

	userID, err := uuid.Parse(claims["public_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "Token Invalid", err.Error())
	}
	board.OwnerPublicID = userID

	if err := controller.service.CreateBoard(board); err != nil {
		return utils.BadRequest(ctx, "Save data failed", err.Error())
	}

	return utils.Success(ctx, "Create board success", board)
}

func (controller *BoardController) UpdateBoard(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")
	board := new(model.Board)

	if err := ctx.Bind().Body(board); err != nil {
		return utils.BadRequest(ctx, "Parsing data failed", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Invalid ID", err.Error())
	}

	findBoard, err := controller.service.FindBoardByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board data not found", err.Error())
	}

	board.InternalID = findBoard.InternalID
	board.PublicID = findBoard.PublicID
	board.OwnerID = findBoard.OwnerID
	board.OwnerPublicID = findBoard.OwnerPublicID
	board.CreatedAt = findBoard.CreatedAt

	if err := controller.service.UpdateBoard(board); err != nil {
		return utils.BadRequest(ctx, "Update board failed", err.Error())
	}

	return utils.Success(ctx, "Update board success", board)
}

func (controller *BoardController) GetAllBoardByUserID(ctx fiber.Ctx) error {
	user := jwtware.FromContext(ctx)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["public_id"].(string)

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	search := ctx.Query("search", "")
	sort := ctx.Query("sort", "")

	boards, total, err := controller.service.FindBoardByUserID(userID, search, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Board data not found", err.Error())
	}

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Search:    search,
		Sort:      sort,
	}

	if total == 0 {
		return utils.PaginationNotFound(ctx, "Board data not found", boards, meta)
	}

	return utils.PaginationSuccess(ctx, "Get board data success", boards, meta)
}

func (controller *BoardController) AddBoardMembers(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.Bind().Body(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Parsing data failed", err.Error())
	}

	if err := controller.service.AddBoardMember(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Add board member failed", err.Error())
	}

	return utils.Success(ctx, "Add board member success", nil)
}

func (controller *BoardController) RemoveBoardMembers(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.Bind().Body(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Parsing data failed", err.Error())
	}

	if err := controller.service.RemoveBoardMember(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Remove board member failed", err.Error())
	}

	return utils.Success(ctx, "Remove board member success", nil)
}
