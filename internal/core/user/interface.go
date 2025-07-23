package user

import (
	"context"
	"gorm.io/gorm"
	"starter/internal/core/pagination"
)

type Repository interface {
	Save(db *gorm.DB, user *User) error
	Update(db *gorm.DB, user *User) error
	Delete(db *gorm.DB, id int) error
	FindAll(db *gorm.DB, params pagination.Request) ([]User, int64, error)
	FindByID(db *gorm.DB, id int) (User, error)
	FindByEmail(db *gorm.DB, email string) (User, error)
}

type Usecase interface {
	Save(ctx context.Context, request CreateRequest) (Response, error)
	Update(ctx context.Context, request UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context, request *pagination.Request) (pagination.Page[Response], error)
	FindById(ctx context.Context, id int) (*Response, error)
	FindByEmail(ctx context.Context, email string) (*Response, error)
}
