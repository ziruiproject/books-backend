package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"starter/internal/adapters/api/http"
	"starter/internal/core/pagination"
	"starter/internal/core/user"
	ivalidator "starter/internal/core/validator"
)

type UserHandler struct {
	UserUsecase user.Usecase
}

func NewUserHandler(userUsecase user.Usecase) *UserHandler {
	return &UserHandler{
		UserUsecase: userUsecase,
	}
}

func (handler *UserHandler) Create(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	request := new(user.CreateRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}

	response, err := handler.UserUsecase.Save(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create user")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create user"),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		http.SuccessResponse(response, "User created successfully"),
	)
}

func (handler *UserHandler) Update(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	id, _ := ctx.ParamsInt("id")
	request := new(user.UpdateRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}
	request.Id = id

	response, err := handler.UserUsecase.Update(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create user")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create user"),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		http.SuccessResponse(response, "User updated successfully"),
	)
}

func (handler *UserHandler) Delete(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	err := handler.UserUsecase.Delete(ctx.UserContext(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create user"),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		http.SuccessResponse("", "User deleted successfully"),
	)
}

func (handler *UserHandler) List(ctx *fiber.Ctx) error {
	sortBy := ctx.Query("sort_by")
	orderBy := ctx.Query("order_by")
	page := ctx.QueryInt("page")
	limit := ctx.QueryInt("limit")

	request := pagination.Request{
		SortBy:  sortBy,
		OrderBy: orderBy,
		Page:    page,
		Limit:   limit,
	}

	response, err := handler.UserUsecase.FindAll(ctx.UserContext(), &request)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to fetch users"),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		http.SuccessResponse(response, "Users fetched successfully"),
	)
}

func (handler *UserHandler) GetByID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	response, err := handler.UserUsecase.FindById(ctx.UserContext(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create user"),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		http.SuccessResponse(response, "User fetched successfully"),
	)
}

func (handler *UserHandler) GetByEmail(ctx *fiber.Ctx) error {
	email := ctx.Params("email")

	response, err := handler.UserUsecase.FindByEmail(ctx.UserContext(), email)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create user"),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		http.SuccessResponse(response, "User fetched successfully"),
	)
}
