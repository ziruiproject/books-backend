package route

import (
	"github.com/gofiber/fiber/v2"
	"starter/internal/adapters/api/http/handler"
	"starter/internal/adapters/api/http/middleware"
)

type UserRoutes struct {
	userHandler *handler.UserHandler
}

func NewUserRoutes(userHandler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		userHandler: userHandler,
	}
}

func (r *UserRoutes) InstallRoutes(app fiber.Router) {
	userGroup := app.Group("/users",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("admin"),
	)

	userGroup.Post("/", r.userHandler.Create)
	userGroup.Get("/", r.userHandler.List)
	userGroup.Get("/:id", r.userHandler.GetByID)
	userGroup.Get("/:email", r.userHandler.GetByEmail)
	userGroup.Put("/:id", r.userHandler.Update)
	userGroup.Delete("/:id", r.userHandler.Delete)
}
