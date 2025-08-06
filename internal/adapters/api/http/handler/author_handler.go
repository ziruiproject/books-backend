package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"starter/internal/adapters/api/http"
	"starter/internal/core/author"
	"starter/internal/core/pagination"
	ivalidator "starter/internal/core/validator"
)

type AuthorHandler struct {
	AuthorUsecase author.Usecase
}

func NewAuthorHandler(authorUsecase author.Usecase) *AuthorHandler {
	return &AuthorHandler{
		AuthorUsecase: authorUsecase,
	}
}

func (handler *AuthorHandler) Create(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	request := new(author.CreateRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}

	response, err := handler.AuthorUsecase.Save(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create author")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to create author")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create author"),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		http.SuccessResponse(response, "Author created successfully"),
	)
}

func (handler *AuthorHandler) Update(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	id, _ := ctx.ParamsInt("id")
	request := new(author.UpdateRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}
	request.Id = id

	response, err := handler.AuthorUsecase.Update(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create author")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to create author")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create author"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Author updated successfully"),
	)
}

func (handler *AuthorHandler) Delete(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	err := handler.AuthorUsecase.Delete(ctx.UserContext(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create author")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create author"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse("", "Author deleted successfully"),
	)
}

func (handler *AuthorHandler) List(ctx *fiber.Ctx) error {
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

	response, err := handler.AuthorUsecase.FindAll(ctx.UserContext(), &request)
	if err != nil {
		log.Error().Err(err).Msg("failed to create author")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to fetch authors"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Authors fetched successfully"),
	)
}

func (handler *AuthorHandler) GetByID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	response, err := handler.AuthorUsecase.FindById(ctx.UserContext(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create author")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create author"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Author fetched successfully"),
	)
}
