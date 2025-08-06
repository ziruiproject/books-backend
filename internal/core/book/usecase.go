package book

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"starter/internal/core/filter"
	"starter/internal/core/pagination"
	"starter/internal/core/storage"
	ivalidator "starter/internal/core/validator"
	"starter/pkg/helper"
)

type UsecaseDependency struct {
	DB             *gorm.DB
	Validator      ivalidator.Validator
	Storage        storage.Storage
	BookRepository Repository
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

	book := request.ToEntity()
	err := usecase.BookRepository.Save(tx, book)
	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return Response{}, errors.New("something went wrong")
	}

	return *ToResponse(book), nil
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

	book, err := usecase.BookRepository.FindByID(tx, request.Id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find book by id: %+v", request.Id)
		return nil, errors.New("book not found")
	}
	updated := helper.Differ(book, *request.ToEntity()).(Book)

	err = usecase.BookRepository.Update(tx, &updated)
	if err != nil {
		log.Error().Err(err).Msgf("failed to update book")
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

	err := usecase.BookRepository.Delete(tx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to update book")
		return errors.New("something went wrong")
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return errors.New("something went wrong")
	}

	return nil
}

func (usecase *UsecaseImpl) FindAll(ctx context.Context, request *pagination.Request, query *filter.BookFilter) (pagination.Page[Response], error) {
	tx := usecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	pagination.NewPagination(request)
	filter.NewBookFilter(query)

	validation := usecase.Validator.ValidateStruct(request)
	if validation != nil {
		log.Error().Msgf("validation error: %s", validation)
		return pagination.Page[Response]{}, ivalidator.ValidationErrors{
			Errors: validation,
		}
	}

	validation = usecase.Validator.ValidateStruct(query)
	if validation != nil {
		log.Error().Msgf("validation error: %s", validation)
		return pagination.Page[Response]{}, ivalidator.ValidationErrors{
			Errors: validation,
		}
	}

	books, count, err := usecase.BookRepository.FindAll(tx, *request, *query)
	if err != nil {
		log.Error().Err(err).Msgf("failed to fetch books")
		return pagination.Page[Response]{}, errors.New("something went wrong")
	}

	var response []Response
	for _, book := range books {
		cover := usecase.Storage.Download(ctx, storage.DownloadRequest{
			Bucket: "bucket",
			Name:   book.Cover,
		})

		book.Cover = ""
		if cover != nil {
			book.Cover = cover.Link
		}
		response = append(response, *ToResponse(&book))
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

	book, err := usecase.BookRepository.FindByID(tx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to find book with id: %d", id)
		return nil, errors.New("book not found")
	}

	cover := usecase.Storage.Download(ctx, storage.DownloadRequest{
		Bucket: "bucket",
		Name:   book.Cover,
	})

	book.Cover = ""
	if cover != nil {
		book.Cover = cover.Link
	}

	if err = tx.Commit().Error; err != nil {
		log.Error().Err(err).Msgf("failed to commit transaction: %+v", err)
		return nil, errors.New("something went wrong")
	}
	return ToResponse(&book), nil
}
