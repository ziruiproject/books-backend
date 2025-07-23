package auth

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"starter/internal/core/role"
	"starter/internal/core/user"
	ivalidator "starter/internal/core/validator"
	"starter/pkg/hasher"
)

type UsecaseDependency struct {
	Hasher         hasher.Hasher
	TokenGenerator TokenGenerator
	Validator      ivalidator.Validator
	UserRepository user.Repository
	DB             *gorm.DB
}

type UsecaseImpl struct {
	UsecaseDependency
}

func NewUsecase(deps UsecaseDependency) Usecase {
	return &UsecaseImpl{
		deps,
	}
}

func (usecase *UsecaseImpl) Login(ctx context.Context, request LoginRequest) (*Response, error) {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	validation := usecase.Validator.ValidateStruct(request)
	if validation != nil {
		log.Error().Msgf("Validation error: %s", validation)
		return &Response{}, ivalidator.ValidationErrors{
			Errors: validation,
		}
	}

	user, err := usecase.UserRepository.FindByEmail(tx, request.Email)
	if !usecase.Hasher.Check(request.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := usecase.TokenGenerator.Generate(AuthenticatedUser{
		Id:   user.ID,
		Role: user.Role.Name,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to generate token")
		return nil, errors.New("something went wrong")
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return nil, errors.New("something went wrong")
	}
	return &Response{
		Token: token,
	}, nil
}

func (usecase *UsecaseImpl) Register(ctx context.Context, request RegisterRequest) (*Response, error) {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	validation := usecase.Validator.ValidateStruct(request)
	if validation != nil {
		log.Error().Msgf("Validation error: %s", validation)
		return &Response{}, ivalidator.ValidationErrors{
			Errors: validation,
		}
	}

	user := request.ToEntity()
	password, err := usecase.Hasher.Hash(user.Password)
	if err != nil {
		log.Error().Err(err).Msgf("failed to secure password")
		return nil, errors.New("something went wrong")
	}

	user.Password = password
	user.Role = role.Role{
		Name: "user",
	}

	err = usecase.UserRepository.Save(tx, user)
	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return nil, errors.New("something went wrong")
	}

	return nil, nil
}

func (usecase *UsecaseImpl) Current(ctx context.Context) (*CurrentAuthResponse, error) {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	claim := ctx.Value("user").(AuthenticatedUser)

	user, err := usecase.UserRepository.FindByID(tx, int(claim.Id))
	if err != nil {
		log.Error().Err(err).Msgf("failed to get current user")
		return nil, errors.New("user not found")
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return nil, errors.New("something went wrong")
	}
	return ToResponse(&user), err
}

func (usecase *UsecaseImpl) Logout(ctx context.Context) error {
	// Stateless logout → nothing to do
	return nil
}
