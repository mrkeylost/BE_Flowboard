package utils

import "github.com/gofiber/fiber/v3"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginationResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Status  int            `json:"status"`
	Data    interface{}    `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
	Meta    PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page      int    `json:"page" example:"1"`
	Limit     int    `json:"limit" example:"10"`
	Total     int    `json:"total" example:"100"`
	TotalPage int    `json:"total_pages" example:"10"`
	Search    string `json:"search" example:"name=john"`
	Sort      string `json:"sort" example:"-id"`
}

func Success(ctx fiber.Ctx, message string, data interface{}) error {
	return ctx.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Status:  fiber.StatusOK,
		Message: message,
		Data:    data,
	})
}

func PaginationSuccess(ctx fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return ctx.Status(fiber.StatusOK).JSON(PaginationResponse{
		Success: true,
		Status:  fiber.StatusOK,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Created(ctx fiber.Ctx, message string, data interface{}) error {
	return ctx.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Status:  fiber.StatusCreated,
		Message: message,
		Data:    data,
	})
}

func Unauthorized(ctx fiber.Ctx, message string, err string) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(Response{
		Success: false,
		Status:  fiber.StatusUnauthorized,
		Message: message,
		Error:   err,
	})
}

func BadRequest(ctx fiber.Ctx, message string, err string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Status:  fiber.StatusBadRequest,
		Message: message,
		Error:   err,
	})
}

func InternalServerError(ctx fiber.Ctx, message string, err string) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Status:  fiber.StatusInternalServerError,
		Message: message,
		Error:   err,
	})
}

func NotFound(ctx fiber.Ctx, message string, err string) error {
	return ctx.Status(fiber.StatusNotFound).JSON(Response{
		Success: false,
		Status:  fiber.StatusNotFound,
		Message: message,
		Error:   err,
	})
}

func PaginationNotFound(ctx fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return ctx.Status(fiber.StatusNotFound).JSON(PaginationResponse{
		Success: false,
		Status:  fiber.StatusNotFound,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
