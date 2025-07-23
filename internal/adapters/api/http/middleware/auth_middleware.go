package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"starter/config"
	"starter/internal/adapters/api/http"
)

func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(config.AppConfig.JWTSecret),
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: " + err.Error(),
			})
		},
	})
}

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userToken, ok := ctx.Locals("user").(*jwt.Token)
		if !ok {
			log.Error().Msg("failed to get current user")
			return ctx.Status(fiber.StatusUnauthorized).JSON(
				http.ErrorResponse("Invalid or expired token"),
			)
		}

		claims, ok := userToken.Claims.(jwt.MapClaims)
		if !ok {
			log.Error().Msg("failed to parse token")
			return ctx.Status(fiber.StatusUnauthorized).JSON(
				http.ErrorResponse("Invalid or expired token"),
			)
		}

		role, _ := claims["role"].(string)
		userID, _ := claims["user_id"].(string)

		for _, allowed := range allowedRoles {
			if role == allowed {
				ctx.Locals("user_id", userID)
				ctx.Locals("role", role)
				return ctx.Next()
			}
		}

		return ctx.Status(fiber.StatusForbidden).JSON(
			http.ErrorResponse("Insufficient roles"),
		)
	}
}
