package route

import (
	"github.com/gofiber/fiber/v2"
	"starter/internal/adapters/api/http/handler"
	"starter/internal/adapters/api/http/middleware"
)

type PublisherRoutes struct {
	publisherHandler *handler.PublisherHandler
}

func NewPublisherRoutes(publisherHandler *handler.PublisherHandler) *PublisherRoutes {
	return &PublisherRoutes{
		publisherHandler: publisherHandler,
	}
}

func (r *PublisherRoutes) InstallRoutes(app fiber.Router) {
	publisherGroup := app.Group("/publishers",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("admin", "user"),
	)

	publisherGroup.Post("/", r.publisherHandler.Create)
	publisherGroup.Get("/", r.publisherHandler.List)
	publisherGroup.Get("/:id", r.publisherHandler.GetByID)
	publisherGroup.Put("/:id", r.publisherHandler.Update)
	publisherGroup.Delete("/:id", r.publisherHandler.Delete)
}
