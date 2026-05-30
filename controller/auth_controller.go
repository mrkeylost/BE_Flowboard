package controller

import (
	"github.com/gofiber/fiber/v3"
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
