package auth

import (
	"context"
)

type Usecase interface {
	Login(ctx context.Context, request LoginRequest) (*Response, error)
	Register(ctx context.Context, request RegisterRequest) (*Response, error)
	Current(ctx context.Context) (*CurrentAuthResponse, error)
	Logout(ctx context.Context) error
}

type TokenGenerator interface {
	Generate(user AuthenticatedUser) (string, error)
}
