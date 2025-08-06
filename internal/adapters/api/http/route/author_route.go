package route

import (
	"github.com/gofiber/fiber/v2"
	"starter/internal/adapters/api/http/handler"
	"starter/internal/adapters/api/http/middleware"
)

type AuthorRoutes struct {
	authorHandler *handler.AuthorHandler
}

func NewAuthorRoutes(authorHandler *handler.AuthorHandler) *AuthorRoutes {
	return &AuthorRoutes{
		authorHandler: authorHandler,
	}
}

func (r *AuthorRoutes) InstallRoutes(app fiber.Router) {
	authorGroup := app.Group("/authors",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("admin", "user"),
	)

	authorGroup.Post("/", r.authorHandler.Create)
	authorGroup.Get("/", r.authorHandler.List)
	authorGroup.Get("/:id", r.authorHandler.GetByID)
	authorGroup.Put("/:id", r.authorHandler.Update)
	authorGroup.Delete("/:id", r.authorHandler.Delete)
}
