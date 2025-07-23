package user

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"starter/internal/core/pagination"
	"starter/internal/core/role"
	ivalidator "starter/internal/core/validator"
	"starter/pkg/hasher"
	"starter/pkg/helper"
)

type UsecaseDependency struct {
	DB             *gorm.DB
	Hasher         hasher.Hasher
	Validator      ivalidator.Validator
	UserRepository Repository
}

type UsecaseImpl struct {
	UsecaseDependency
}

func NewUsecase(deps UsecaseDependency) Usecase {
	return &UsecaseImpl{
		deps,
	}
}

func (usecase *UsecaseImpl) Save(ctx context.Context, request CreateRequest) (Response, error) {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	validation := usecase.Validator.ValidateStruct(request)
	if validation != nil {
		log.Error().Msgf("Validation error: %s", validation)
		return Response{}, ivalidator.ValidationErrors{
			Errors: validation,
		}
	}

	user := request.ToEntity()
	password, err := usecase.Hasher.Hash(user.Password)
	if err != nil {
		log.Error().Err(err).Msgf("failed to secure password")
		return Response{}, errors.New("something went wrong")
	}
	user.Password = password
	user.Role = role.Role{
		Name: "user",
	}

	err = usecase.UserRepository.Save(tx, user)
	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return Response{}, errors.New("something went wrong")
	}

	return *ToResponse(user), nil
}

func (usecase *UsecaseImpl) Update(ctx context.Context, request UpdateRequest) (*Response, error) {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	validation := usecase.Validator.ValidateStruct(request)
	if validation != nil {
		log.Error().Msgf("Validation error: %s", validation)
		return &Response{}, ivalidator.ValidationErrors{
			Errors: validation,
		}
	}

	user, err := usecase.UserRepository.FindByID(tx, request.Id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find user by id: %+v", request.Id)
		return nil, errors.New("user not found")
	}
	updated := helper.Differ(user, *request.ToEntity()).(User)

	if updated.Password != "" {
		password, err := usecase.Hasher.Hash(user.Password)
		if err != nil {
			log.Error().Err(err).Msgf("failed to secure password")
			return &Response{}, errors.New("something went wrong")
		}
		user.Password = password
	}

	err = usecase.UserRepository.Update(tx, &updated)
	if err != nil {
		log.Error().Err(err).Msgf("failed to update user")
		return nil, errors.New("something went wrong")
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return &Response{}, errors.New("something went wrong")
	}
	return ToResponse(&updated), nil
}

func (usecase *UsecaseImpl) Delete(ctx context.Context, id int) error {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := usecase.UserRepository.Delete(tx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to update user")
		return errors.New("something went wrong")
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return errors.New("something went wrong")
	}

	return nil
}

func (usecase *UsecaseImpl) FindAll(ctx context.Context, request *pagination.Request) (pagination.Page[Response], error) {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	pagination.NewPagination(request)
	users, count, err := usecase.UserRepository.FindAll(tx, *request)
	if err != nil {
		log.Error().Err(err).Msgf("failed to fetch users")
		return pagination.Page[Response]{}, errors.New("something went wrong")
	}

	var response []Response
	for _, user := range users {
		response = append(response, *ToResponse(&user))
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return pagination.Page[Response]{}, errors.New("something went wrong")
	}

	return *pagination.NewPage[Response](*request, count, response), nil
}

func (usecase *UsecaseImpl) FindById(ctx context.Context, id int) (*Response, error) {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user, err := usecase.UserRepository.FindByID(tx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find user with id: %d", id)
		return nil, errors.New("user not found")
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return nil, errors.New("something went wrong")
	}
	return ToResponse(&user), nil
}

func (usecase *UsecaseImpl) FindByEmail(ctx context.Context, email string) (*Response, error) {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user, err := usecase.UserRepository.FindByEmail(tx, email)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find user with email: %s", email)
		return nil, errors.New("user not found")
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return nil, errors.New("something went wrong")
	}
	return ToResponse(&user), nil
}
