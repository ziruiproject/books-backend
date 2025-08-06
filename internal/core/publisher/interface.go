package publisher

import (
	"context"
	"gorm.io/gorm"
	"starter/internal/core/pagination"
)

type Repository interface {
	Save(db *gorm.DB, Publisher *Publisher) error
	Update(db *gorm.DB, Publisher *Publisher) error
	Delete(db *gorm.DB, id int) error
	FindAll(db *gorm.DB, params pagination.Request) ([]Publisher, int64, error)
	FindByID(db *gorm.DB, id int) (Publisher, error)
}

type Usecase interface {
	Save(ctx context.Context, request CreateRequest) (Response, error)
	Update(ctx context.Context, request UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context, request *pagination.Request) (pagination.Page[Response], error)
	FindById(ctx context.Context, id int) (*Response, error)
}
