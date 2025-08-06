package category

import (
	"gorm.io/gorm"
	"strings"
)

type CreateRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type UpdateRequest struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"max=100"`
}

type Response struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (dto *CreateRequest) ToEntity() *Category {
	return &Category{
		Name: strings.ToUpper(dto.Name),
	}
}

func (dto *UpdateRequest) ToEntity() *Category {
	return &Category{
		Model: gorm.Model{ID: uint(dto.Id)},
		Name:  dto.Name,
	}
}

func ToResponse(entity *Category) *Response {
	return &Response{
		Id:   int(entity.ID),
		Name: entity.Name,
	}
}
