package handler

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"starter/internal/adapters/api/http"
	"starter/internal/core/auth"
	ivalidator "starter/internal/core/validator"
)

type AuthHandler struct {
	AuthUsecase auth.Usecase
}

func NewAuthHandler(authUsecase auth.Usecase) *AuthHandler {
	return &AuthHandler{
		AuthUsecase: authUsecase,
	}
}

func (handler *AuthHandler) Login(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	request := new(auth.LoginRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}

	response, err := handler.AuthUsecase.Login(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to login")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to login")
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			http.ErrorResponse("Invalid email or password"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "Login successful"),
	)
}

func (handler *AuthHandler) Register(ctx *fiber.Ctx) error {
	var validationError ivalidator.ValidationErrors

	request := new(auth.RegisterRequest)
	if err := ctx.BodyParser(&request); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			http.ErrorResponse("Invalid request body"),
		)
	}

	response, err := handler.AuthUsecase.Register(ctx.UserContext(), *request)
	if errors.As(err, &validationError) {
		log.Error().Err(err).Msg("failed to create user")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ValidationResponse(validationError),
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to register user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to register user"),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		http.SuccessResponse(response, "User registered successfully"),
	)
}

func (handler *AuthHandler) Current(ctx *fiber.Ctx) error {
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok {
		log.Error().Msg("failed to get current user")
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			http.ErrorResponse("Invalid or expired token"),
		)
	}

	claim, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Error().Msg("failed to parse token")
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			http.ErrorResponse("Invalid or expired token"),
		)
	}

	authCtx := context.WithValue(ctx.Context(), "user", auth.AuthenticatedUser{
		Id:   uint(claim["user_id"].(float64)),
		Role: claim["role"].(string),
	})
	ctx.SetUserContext(authCtx)

	response, err := handler.AuthUsecase.Current(ctx.UserContext())
	if err != nil {
		log.Error().Err(err).Msg("failed to login")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Invalid email or password"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(response, "User fetched successfully"),
	)
}
