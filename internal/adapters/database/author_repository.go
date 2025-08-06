package database

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"starter/internal/core/author"
	"starter/internal/core/pagination"
)

type AuthorRepository struct {
}

func NewAuthorRepository() author.Repository {
	return &AuthorRepository{}
}

func (repository *AuthorRepository) Save(db *gorm.DB, author *author.Author) error {
	result := db.Create(author)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save author")
		return result.Error
	}

	return nil
}

func (repository *AuthorRepository) Update(db *gorm.DB, author *author.Author) error {
	result := db.Updates(author)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save author")
		return result.Error
	}

	return nil
}

func (repository *AuthorRepository) Delete(db *gorm.DB, id int) error {
	result := db.Delete(&author.Author{}, id)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to delete author")
		return result.Error
	}

	return nil
}

func (repository *AuthorRepository) FindAll(db *gorm.DB, params pagination.Request) ([]author.Author, int64, error) {
	var authors []author.Author
	var count int64

	db.Model(&authors).Count(&count)

	result := db.
		Debug().
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.SortBy)).
		Limit(params.Limit).
		Offset((params.Page - 1) * params.Limit).
		Find(&authors)

	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to find all authors")
	}

	return authors, count, result.Error
}

func (repository *AuthorRepository) FindByID(db *gorm.DB, id int) (author.Author, error) {
	var author author.Author
	result := db.First(&author, id)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to find author")
		return author, result.Error
	}

	return author, nil
}
