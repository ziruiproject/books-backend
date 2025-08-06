package route

import (
	"github.com/gofiber/fiber/v2"
	"starter/internal/adapters/api/http/handler"
	"starter/internal/adapters/api/http/middleware"
)

type BookRoutes struct {
	bookHandler *handler.BookHandler
}

func NewBookRoutes(bookHandler *handler.BookHandler) *BookRoutes {
	return &BookRoutes{
		bookHandler: bookHandler,
	}
}

func (r *BookRoutes) InstallRoutes(app fiber.Router) {
	bookGroup := app.Group("/books",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("admin", "user"),
	)

	bookGroup.Post("/", r.bookHandler.Create)
	bookGroup.Get("/", r.bookHandler.List)
	bookGroup.Get("/:id", r.bookHandler.GetByID)
	bookGroup.Put("/:id", r.bookHandler.Update)
	bookGroup.Delete("/:id", r.bookHandler.Delete)
}
