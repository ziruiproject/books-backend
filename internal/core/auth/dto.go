package auth

import (
	"starter/internal/core/user"
	"strconv"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type CurrentAuthResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Response struct {
	Token string `json:"token"`
}

type AuthenticatedUser struct {
	Id   uint   `json:"user_id"`
	Role string `json:"role"`
}

func ToResponse(entity *user.User) *CurrentAuthResponse {
	return &CurrentAuthResponse{
		Id:    strconv.Itoa(int(entity.ID)),
		Name:  entity.Name,
		Email: entity.Email,
		Role:  entity.Role.Name,
	}
}

func (request *RegisterRequest) ToEntity() *user.User {
	return &user.User{
		Name:     strings.ToUpper(request.Name),
		Email:    request.Email,
		Password: request.Password,
	}
}
