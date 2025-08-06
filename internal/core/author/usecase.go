package author

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
	DB               *gorm.DB
	Validator        ivalidator.Validator
	AuthorRepository Repository
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

	author := request.ToEntity()
	err := usecase.AuthorRepository.Save(tx, author)
	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return Response{}, errors.New("something went wrong")
	}

	return *ToResponse(author), nil
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

	author, err := usecase.AuthorRepository.FindByID(tx, request.Id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find author by id: %+v", request.Id)
		return nil, errors.New("author not found")
	}
	updated := helper.Differ(author, *request.ToEntity()).(Author)

	err = usecase.AuthorRepository.Update(tx, &updated)
	if err != nil {
		log.Error().Err(err).Msgf("failed to update author")
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

	err := usecase.AuthorRepository.Delete(tx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to update author")
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
	authors, count, err := usecase.AuthorRepository.FindAll(tx, *request)
	if err != nil {
		log.Error().Err(err).Msgf("failed to fetch authors")
		return pagination.Page[Response]{}, errors.New("something went wrong")
	}

	var response []Response
	for _, author := range authors {
		response = append(response, *ToResponse(&author))
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

	author, err := usecase.AuthorRepository.FindByID(tx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find author with id: %d", id)
		return nil, errors.New("author not found")
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return nil, errors.New("something went wrong")
	}
	return ToResponse(&author), nil
}
