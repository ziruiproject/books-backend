package database

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"starter/internal/core/pagination"
	"starter/internal/core/publisher"
)

type PublisherRepository struct {
}

func NewPublisherRepository() publisher.Repository {
	return &PublisherRepository{}
}

func (repository *PublisherRepository) Save(db *gorm.DB, publisher *publisher.Publisher) error {
	result := db.Create(publisher)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save publisher")
		return result.Error
	}

	return nil
}

func (repository *PublisherRepository) Update(db *gorm.DB, publisher *publisher.Publisher) error {
	result := db.Updates(publisher)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save publisher")
		return result.Error
	}

	return nil
}

func (repository *PublisherRepository) Delete(db *gorm.DB, id int) error {
	result := db.Delete(&publisher.Publisher{}, id)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to delete publisher")
		return result.Error
	}

	return nil
}

func (repository *PublisherRepository) FindAll(db *gorm.DB, params pagination.Request) ([]publisher.Publisher, int64, error) {
	var publishers []publisher.Publisher
	var count int64

	db.Model(&publishers).Count(&count)

	result := db.
		Debug().
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.SortBy)).
		Limit(params.Limit).
		Offset((params.Page - 1) * params.Limit).
		Find(&publishers)

	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to find all publishers")
	}

	return publishers, count, result.Error
}

func (repository *PublisherRepository) FindByID(db *gorm.DB, id int) (publisher.Publisher, error) {
	var publisher publisher.Publisher
	result := db.First(&publisher, id)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to find publisher")
		return publisher, result.Error
	}

	return publisher, nil
}
