package publisher

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
	DB                  *gorm.DB
	Validator           ivalidator.Validator
	PublisherRepository Repository
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

	publisher := request.ToEntity()
	err := usecase.PublisherRepository.Save(tx, publisher)
	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return Response{}, errors.New("something went wrong")
	}

	return *ToResponse(publisher), nil
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

	publisher, err := usecase.PublisherRepository.FindByID(tx, request.Id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find publisher by id: %+v", request.Id)
		return nil, errors.New("publisher not found")
	}
	updated := helper.Differ(publisher, *request.ToEntity()).(Publisher)

	err = usecase.PublisherRepository.Update(tx, &updated)
	if err != nil {
		log.Error().Err(err).Msgf("failed to update publisher")
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

	err := usecase.PublisherRepository.Delete(tx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to update publisher")
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
	publishers, count, err := usecase.PublisherRepository.FindAll(tx, *request)
	if err != nil {
		log.Error().Err(err).Msgf("failed to fetch publishers")
		return pagination.Page[Response]{}, errors.New("something went wrong")
	}

	var response []Response
	for _, publisher := range publishers {
		response = append(response, *ToResponse(&publisher))
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

	publisher, err := usecase.PublisherRepository.FindByID(tx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find publisher with id: %d", id)
		return nil, errors.New("publisher not found")
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return nil, errors.New("something went wrong")
	}
	return ToResponse(&publisher), nil
}
