package database

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"starter/internal/core/book"
	"starter/internal/core/filter"
	"starter/internal/core/pagination"
)

type BookRepository struct {
}

func NewBookRepository() book.Repository {
	return &BookRepository{}
}

func (repository *BookRepository) Save(db *gorm.DB, book *book.Book) error {
	result := db.Create(book)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save book")
		return result.Error
	}
	result = db.
		Preload("Author").
		Preload("Categories").
		First(&book, book.ID)

	return nil
}

func (repository *BookRepository) Update(db *gorm.DB, book *book.Book) error {
	result := db.Debug().Updates(book)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save book")
		return result.Error
	}

	return nil
}

func (repository *BookRepository) Delete(db *gorm.DB, id int) error {
	result := db.Delete(&book.Book{}, id)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to delete book")
		return result.Error
	}

	return nil
}

func (repository *BookRepository) FindAll(db *gorm.DB, params pagination.Request, filter filter.BookFilter) ([]book.Book, int64, error) {
	var books []book.Book
	var count int64

	query := db.
		Model(&book.Book{}).
		Joins("LEFT JOIN authors ON authors.id = books.author_id").
		Joins("LEFT JOIN publishers ON publishers.id = books.publisher_id").
		Preload("Author").
		Preload("Publisher").
		Preload("Categories")

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where(`
		books.title ILIKE ? OR
		authors.first_name ILIKE ? OR
		authors.last_name ILIKE ? OR
		publishers.name ILIKE ?
	`, search, search, search, search)
	}

	if filter.StartDate != "" {
		query = query.Where("books.publication_date >= ?", filter.From)
	}
	if filter.EndDate != "" {
		query = query.Where("books.publication_date <= ?", filter.To)
	}
	if filter.From != "" {
		query = query.Where("books.created_at >= ?", filter.StartDate)
	}
	if filter.To != "" {
		query = query.Where("books.created_at <= ?", filter.EndDate)
	}

	if len(filter.Categories) > 0 {
		query = query.
			Joins("JOIN book_category bc ON bc.book_id = books.id").
			Where("bc.category_id IN ?", filter.Categories)
	}

	query.Count(&count)

	result := query.
		Debug().
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.SortBy)).
		Limit(params.Limit).
		Offset((params.Page - 1) * params.Limit).
		Find(&books)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to find all books")
	}

	return books, count, result.Error
}

func (repository *BookRepository) FindByID(db *gorm.DB, id int) (book.Book, error) {
	var book book.Book
	result := db.
		Preload("Author").
		Preload("Publisher").
		Preload("Categories").
		First(&book, id)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to find book")
		return book, result.Error
	}

	return book, nil
}
