package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"starter/internal/adapters/api/http"
	"starter/internal/core/category"
	"starter/internal/core/pagination"
	ivalidator "starter/internal/core/validator"
)

type CategoryHandler struct {
	CategoryUsecase category.Usecase
}

func NewCategoryHandler(categoryUsecase category.Usecase) *CategoryHandler {
	return &CategoryHandler{
		CategoryUsecase: categoryUsecase,
	}
}

func (handler *CategoryHandler) Create(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	request := new(category.CreateRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}

	response, err := handler.CategoryUsecase.Save(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create category")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to create category")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create category"),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		http.SuccessResponse(response, "Category created successfully"),
	)
}

func (handler *CategoryHandler) Update(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	id, _ := ctx.ParamsInt("id")
	request := new(category.UpdateRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}
	request.Id = id

	response, err := handler.CategoryUsecase.Update(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create category")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to create category")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create category"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Category updated successfully"),
	)
}

func (handler *CategoryHandler) Delete(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	err := handler.CategoryUsecase.Delete(ctx.UserContext(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create category")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create category"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse("", "Category deleted successfully"),
	)
}

func (handler *CategoryHandler) List(ctx *fiber.Ctx) error {
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

	response, err := handler.CategoryUsecase.FindAll(ctx.UserContext(), &request)
	if err != nil {
		log.Error().Err(err).Msg("failed to create category")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to fetch categories"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Categories fetched successfully"),
	)
}

func (handler *CategoryHandler) GetByID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	response, err := handler.CategoryUsecase.FindById(ctx.UserContext(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create category")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create category"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Category fetched successfully"),
	)
}
