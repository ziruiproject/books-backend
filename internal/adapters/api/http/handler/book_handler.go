package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"starter/internal/adapters/api/http"
	"starter/internal/core/book"
	"starter/internal/core/filter"
	"starter/internal/core/pagination"
	ivalidator "starter/internal/core/validator"
	"starter/pkg/helper"
)

type BookHandler struct {
	BookUsecase book.Usecase
}

func NewBookHandler(bookUsecase book.Usecase) *BookHandler {
	return &BookHandler{
		BookUsecase: bookUsecase,
	}
}

func (handler *BookHandler) Create(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	request := new(book.CreateRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}

	response, err := handler.BookUsecase.Save(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create book")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to create book")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create book"),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		http.SuccessResponse(response, "Book created successfully"),
	)
}

func (handler *BookHandler) Update(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	id, _ := ctx.ParamsInt("id")
	request := new(book.UpdateRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}
	request.Id = id

	response, err := handler.BookUsecase.Update(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create book")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to create book")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create book"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Book updated successfully"),
	)
}

func (handler *BookHandler) Delete(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	err := handler.BookUsecase.Delete(ctx.UserContext(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create book")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create book"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse("", "Book deleted successfully"),
	)
}

func (handler *BookHandler) List(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors
	
	ctg := ctx.Query("categories")
	categories := helper.ParseUintSlice(ctg)

	request := pagination.Request{
		SortBy:  ctx.Query("sort_by"),
		OrderBy: ctx.Query("order_by"),
		Page:    ctx.QueryInt("page"),
		Limit:   ctx.QueryInt("limit"),
	}

	filter := filter.BookFilter{
		Default: filter.Default{
			Search:    ctx.Query("search"),
			StartDate: ctx.Query("start_date"),
			EndDate:   ctx.Query("end_date"),
		},
		Categories: categories,
		From:       ctx.Query("from"),
		To:         ctx.Query("to"),
	}

	response, err := handler.BookUsecase.FindAll(ctx.UserContext(), &request, &filter)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create book")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to fetch book")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to fetch categories"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Books fetched successfully"),
	)
}

func (handler *BookHandler) GetByID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	response, err := handler.BookUsecase.FindById(ctx.UserContext(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create book")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create book"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Book fetched successfully"),
	)
}
