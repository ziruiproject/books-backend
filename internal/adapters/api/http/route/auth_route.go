package route

import (
	"github.com/gofiber/fiber/v2"
	"starter/internal/adapters/api/http/handler"
	"starter/internal/adapters/api/http/middleware"
)

type AuthRoutes struct {
	authHandler *handler.AuthHandler
}

func NewAuthRoutes(authHandler *handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		authHandler: authHandler,
	}
}

func (r *AuthRoutes) InstallRoutes(app fiber.Router) {
	authGroup := app.Group("/auth")

	authGroup.Post("/login", r.authHandler.Login)
	authGroup.Post("/register", r.authHandler.Register)

	authGroup.Get("/current",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("admin", "user"),
		r.authHandler.Current,
	)
}
