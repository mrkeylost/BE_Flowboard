package controller

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/mrkeylost/Flowboard_Backend/model"
	"github.com/mrkeylost/Flowboard_Backend/services"
	"github.com/mrkeylost/Flowboard_Backend/utils"
)

type AuthController struct {
	service services.AuthService
}

func NewUserController(service services.AuthService) *AuthController {
	return &AuthController{service}
}

func (controller *AuthController) Register(ctx fiber.Ctx) error {
	user := new(model.User)

	if err := ctx.Bind().Body(user); err != nil {
		return utils.BadRequest(ctx, "Parsing Data Failed", err.Error())
	}

	if err := controller.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Registration Failed", err.Error())
	}

	var userResponse model.UserResponse
	_ = copier.Copy(&userResponse, &user)

	return utils.Success(ctx, "Register Success", userResponse)
}

func (controller *AuthController) Login(ctx fiber.Ctx) error {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.Bind().Body(&requestBody); err != nil {
		return utils.BadRequest(ctx, "Request format invalid", err.Error())
	}

	user, err := controller.service.Login(requestBody.Email, requestBody.Password)
	if err != nil {
		return utils.Unauthorized(ctx, "Unauthorized user", err.Error())
	}

	token, _ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken, _ := utils.GenerateRefreshToken(user.InternalID)

	var userResponse model.UserResponse
	_ = copier.Copy(&userResponse, &user)

	return utils.Success(ctx, "Login Success", fiber.Map{
		"access_token":  token,
		"refresh_token": refreshToken,
		"user":          userResponse,
	})
}

func (controller *AuthController) UpdateUser(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	publicId, err := uuid.Parse(id)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid ID Format", err.Error())
	}

	var user model.User
	if err := ctx.Bind().Body(&user); err != nil {
		return utils.BadRequest(ctx, "Parsing Data Failed", err.Error())
	}

	user.PublicID = publicId

	if err := controller.service.UpdateUser(&user); err != nil {
		return utils.BadRequest(ctx, "Update data failed", err.Error())
	}

	updatedUser, err := controller.service.GetUserByPublicID(id)
	if err != nil {
		return utils.InternalServerError(ctx, "Data not found", err.Error())
	}

	var userResponse model.UserResponse
	err = copier.Copy(&userResponse, &updatedUser)
	if err != nil {
		return utils.InternalServerError(ctx, "Internal Server Error", err.Error())
	}

	return utils.Success(ctx, "Update user data success", userResponse)
}

func (controller *AuthController) DeleteUser(ctx fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	if err := controller.service.DeleteUser(uint(id)); err != nil {
		return utils.InternalServerError(ctx, "Delete user data failed", err.Error())
	}

	return utils.Success(ctx, "Delete user success", id)
}

func (controller *AuthController) GetUserDetail(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	user, err := controller.service.GetUserByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "User not found", err.Error())
	}

	var userResponse model.UserResponse
	err = copier.Copy(&userResponse, &user)
	if err != nil {
		return utils.InternalServerError(ctx, "Internal Server Error", err.Error())
	}

	return utils.Success(ctx, "Get user detail success", userResponse)
}

func (controller *AuthController) GetAllUser(ctx fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	search := ctx.Query("search", "")
	sort := ctx.Query("sort", "")

	users, total, err := controller.service.GetAllUser(search, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "User data not found", err.Error())
	}

	var userResponse []model.UserResponse
	_ = copier.Copy(&userResponse, &users)

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Search:    search,
		Sort:      sort,
	}

	if total == 0 {
		return utils.PaginationNotFound(ctx, "User data not found", userResponse, meta)
	}

	return utils.PaginationSuccess(ctx, "Get user data success", userResponse, meta)
}
