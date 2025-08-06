package user

import "gorm.io/gorm"

type CreateRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateRequest struct {
	Id    int    `json:"id" validate:"required"`
	Name  string `json:"name" validate:"max=100"`
	Email string `json:"email" validate:"omitempty,email"`
}

type Response struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (dto *CreateRequest) ToEntity() *User {
	return &User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func (dto *UpdateRequest) ToEntity() *User {
	return &User{
		Model: gorm.Model{ID: uint(dto.Id)},
		Name:  dto.Name,
		Email: dto.Email,
	}
}

func ToResponse(entity *User) *Response {
	return &Response{
		Name:  entity.Name,
		Email: entity.Email,
	}
}
