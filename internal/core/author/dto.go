package author

import (
	"gorm.io/gorm"
	"strings"
)

type CreateRequest struct {
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"required,max=100"`
}

type UpdateRequest struct {
	Id        int    `json:"id" validate:"required"`
	FirstName string `json:"first_name" validate:"max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
}

type Response struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (dto *CreateRequest) ToEntity() *Author {
	return &Author{
		FirstName: strings.ToUpper(dto.FirstName),
		LastName:  strings.ToUpper(dto.LastName),
	}
}

func (dto *UpdateRequest) ToEntity() *Author {
	return &Author{
		Model:     gorm.Model{ID: uint(dto.Id)},
		FirstName: strings.ToUpper(dto.FirstName),
		LastName:  strings.ToUpper(dto.LastName),
	}
}

func ToResponse(entity *Author) *Response {
	return &Response{
		Id:        int(entity.ID),
		FirstName: strings.ToUpper(entity.FirstName),
		LastName:  strings.ToUpper(entity.LastName),
	}
}
