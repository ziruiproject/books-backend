package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"starter/internal/adapters/api/http"
	"starter/internal/core/pagination"
	"starter/internal/core/publisher"
	ivalidator "starter/internal/core/validator"
)

type PublisherHandler struct {
	PublisherUsecase publisher.Usecase
}

func NewPublisherHandler(publisherUsecase publisher.Usecase) *PublisherHandler {
	return &PublisherHandler{
		PublisherUsecase: publisherUsecase,
	}
}

func (handler *PublisherHandler) Create(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	request := new(publisher.CreateRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}

	response, err := handler.PublisherUsecase.Save(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create publisher")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to create publisher")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create publisher"),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		http.SuccessResponse(response, "Publisher created successfully"),
	)
}

func (handler *PublisherHandler) Update(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	id, _ := ctx.ParamsInt("id")
	request := new(publisher.UpdateRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}
	request.Id = id

	response, err := handler.PublisherUsecase.Update(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create publisher")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to create publisher")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create publisher"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Publisher updated successfully"),
	)
}

func (handler *PublisherHandler) Delete(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	err := handler.PublisherUsecase.Delete(ctx.UserContext(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create publisher")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create publisher"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse("", "Publisher deleted successfully"),
	)
}

func (handler *PublisherHandler) List(ctx *fiber.Ctx) error {
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

	response, err := handler.PublisherUsecase.FindAll(ctx.UserContext(), &request)
	if err != nil {
		log.Error().Err(err).Msg("failed to create publisher")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to fetch publishers"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Publishers fetched successfully"),
	)
}

func (handler *PublisherHandler) GetByID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	response, err := handler.PublisherUsecase.FindById(ctx.UserContext(), id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create publisher")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to create publisher"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Publisher fetched successfully"),
	)
}
