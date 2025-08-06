package category

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"starter/internal/core/pagination"
	ivalidator "starter/internal/core/validator"
	"starter/pkg/helper"
)

type UsecaseDependency struct {
	DB                 *gorm.DB
	Validator          ivalidator.Validator
	CategoryRepository Repository
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

	category := request.ToEntity()
	err := usecase.CategoryRepository.Save(tx, category)
	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return Response{}, errors.New("something went wrong")
	}

	return *ToResponse(category), nil
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

	category, err := usecase.CategoryRepository.FindByID(tx, request.Id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find category by id: %+v", request.Id)
		return nil, errors.New("category not found")
	}
	updated := helper.Differ(category, *request.ToEntity()).(Category)

	err = usecase.CategoryRepository.Update(tx, &updated)
	if err != nil {
		log.Error().Err(err).Msgf("failed to update category")
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

	err := usecase.CategoryRepository.Delete(tx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to update category")
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
	categorys, count, err := usecase.CategoryRepository.FindAll(tx, *request)
	if err != nil {
		log.Error().Err(err).Msgf("failed to fetch categories")
		return pagination.Page[Response]{}, errors.New("something went wrong")
	}

	var response []Response
	for _, category := range categorys {
		response = append(response, *ToResponse(&category))
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

	category, err := usecase.CategoryRepository.FindByID(tx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find category with id: %d", id)
		return nil, errors.New("category not found")
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return nil, errors.New("something went wrong")
	}
	return ToResponse(&category), nil
}
