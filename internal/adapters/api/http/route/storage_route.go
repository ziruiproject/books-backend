package route

import (
	"github.com/gofiber/fiber/v2"
	"starter/internal/adapters/api/http/handler"
	"starter/internal/adapters/api/http/middleware"
)

type StorageRoutes struct {
	storageHandler *handler.StorageHandler
}

func NewStorageRoutes(storageHandler *handler.StorageHandler) *StorageRoutes {
	return &StorageRoutes{
		storageHandler: storageHandler,
	}
}

func (r *StorageRoutes) InstallRoutes(app fiber.Router) {
	storageGroup := app.Group("/storage",
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("admin", "user"),
	)

	storageGroup.Post("/upload", r.storageHandler.Upload)
}
