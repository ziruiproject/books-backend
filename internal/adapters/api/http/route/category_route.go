package route

import (
	"github.com/gofiber/fiber/v2"
	"starter/internal/adapters/api/http/handler"
	"starter/internal/adapters/api/http/middleware"
)

type CategoryRoutes struct {
	categoryHandler *handler.CategoryHandler
}

func NewCategoryRoutes(categoryHandler *handler.CategoryHandler) *CategoryRoutes {
	return &CategoryRoutes{
		categoryHandler: categoryHandler,
	}
}

func (r *CategoryRoutes) InstallRoutes(app fiber.Router) {
	categoryGroup := app.Group("/categories",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("admin", "user"),
	)

	categoryGroup.Post("/", r.categoryHandler.Create)
	categoryGroup.Get("/", r.categoryHandler.List)
	categoryGroup.Get("/:id", r.categoryHandler.GetByID)
	categoryGroup.Put("/:id", r.categoryHandler.Update)
	categoryGroup.Delete("/:id", r.categoryHandler.Delete)
}
