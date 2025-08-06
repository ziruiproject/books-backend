package database

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"starter/internal/core/category"
	"starter/internal/core/pagination"
)

type CategoryRepository struct {
}

func NewCategoryRepository() category.Repository {
	return &CategoryRepository{}
}

func (repository *CategoryRepository) Save(db *gorm.DB, category *category.Category) error {
	result := db.Create(category)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save category")
		return result.Error
	}

	return nil
}

func (repository *CategoryRepository) Update(db *gorm.DB, category *category.Category) error {
	result := db.Updates(category)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save category")
		return result.Error
	}

	return nil
}

func (repository *CategoryRepository) Delete(db *gorm.DB, id int) error {
	result := db.Delete(&category.Category{}, id)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to delete category")
		return result.Error
	}

	return nil
}

func (repository *CategoryRepository) FindAll(db *gorm.DB, params pagination.Request) ([]category.Category, int64, error) {
	var categories []category.Category
	var count int64

	db.Model(&categories).Count(&count)

	result := db.
		Debug().
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.SortBy)).
		Limit(params.Limit).
		Offset((params.Page - 1) * params.Limit).
		Find(&categories)

	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to find all categories")
	}

	return categories, count, result.Error
}

func (repository *CategoryRepository) FindByID(db *gorm.DB, id int) (category.Category, error) {
	var category category.Category
	result := db.First(&category, id)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to find category")
		return category, result.Error
	}

	return category, nil
}
