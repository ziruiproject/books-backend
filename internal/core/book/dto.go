package book

import (
	"gorm.io/gorm"
	"starter/internal/core/category"
	"strings"
	"time"
)

type CreateRequest struct {
	Title           string `json:"title" validate:"required,max=100"`
	Cover           string `json:"cover" validate:"required"`
	Description     string `json:"description" validate:"required,max=500"`
	PageCount       int    `json:"page_count" validate:"required,min=1,max=10000"`
	AuthorId        int    `json:"author_id" validate:"required"`
	Categories      []int  `json:"categories" validate:"required,dive,min=1"`
	PublisherId     int    `json:"publisher_id" validate:"required"`
	PublicationDate string `json:"publication_date" validate:"required,publication_date"`
}

type UpdateRequest struct {
	Id              int    `json:"id" validate:"required"`
	Title           string `json:"title" validate:"max=100"`
	Cover           string `json:"cover"`
	Description     string `json:"description" validate:"max=500"`
	PageCount       int    `json:"page_count,omitempty" validate:"omitempty,min=1,max=10000"`
	AuthorId        int    `json:"author_id"`
	Categories      []int  `json:"categories,omitempty" validate:"omitempty,dive,min=1"`
	PublisherId     int    `json:"publisher_id" validate:"required"`
	PublicationDate string `json:"publication_date,omitempty" validate:"omitempty,publication_date"`
}

type AuthorResponse struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CategoryResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PublisherResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Id              int                `json:"id"`
	Title           string             `json:"title"`
	Cover           string             `json:"cover"`
	Description     string             `json:"description"`
	PageCount       int                `json:"page_count"`
	Author          AuthorResponse     `json:"author"`
	Categories      []CategoryResponse `json:"categories"`
	Publisher       PublisherResponse  `json:"publisher"`
	PublicationDate string             `json:"publication_date"`
}

func (dto *CreateRequest) ToEntity() *Book {
	categories := make([]category.Category, 0, len(dto.Categories))
	for _, v := range dto.Categories {
		categories = append(categories, category.Category{
			Model: gorm.Model{ID: uint(v)},
		})
	}

	publicationDate, _ := time.Parse("2006-01-02", dto.PublicationDate)

	book := &Book{
		Title:           strings.ToUpper(dto.Title),
		Description:     dto.Description,
		Cover:           dto.Cover,
		PageCount:       dto.PageCount,
		AuthorId:        dto.AuthorId,
		PublisherId:     dto.PublisherId,
		Categories:      categories,
		PublicationDate: publicationDate,
	}
	book.Author.ID = uint(dto.AuthorId)
	book.Publisher.ID = uint(dto.PublisherId)

	return book
}

func (dto *UpdateRequest) ToEntity() *Book {
	categories := make([]category.Category, 0, len(dto.Categories))
	for _, v := range dto.Categories {
		categories = append(categories, category.Category{
			Model: gorm.Model{ID: uint(v)},
		})
	}
	publicationDate, _ := time.Parse("2006-01-02", dto.PublicationDate)

	book := &Book{
		Title:           strings.ToUpper(dto.Title),
		Description:     dto.Description,
		Cover:           dto.Cover,
		PageCount:       dto.PageCount,
		AuthorId:        dto.AuthorId,
		Categories:      categories,
		PublicationDate: publicationDate,
	}
	book.ID = uint(dto.Id)
	book.Author.ID = uint(dto.AuthorId)
	book.Publisher.ID = uint(dto.PublisherId)

	return book
}

func ToResponse(entity *Book) *Response {
	categories := make([]CategoryResponse, 0, len(entity.Categories))
	for _, v := range entity.Categories {
		categories = append(categories, CategoryResponse{
			Id:   int(v.ID),
			Name: v.Name,
		})
	}
	publicationDate := entity.PublicationDate.Format("2006-01-02")

	return &Response{
		Id:          int(entity.ID),
		Title:       strings.ToUpper(entity.Title),
		Cover:       entity.Cover,
		Description: entity.Description,
		PageCount:   entity.PageCount,
		Author: AuthorResponse{
			Id:        int(entity.Author.ID),
			FirstName: entity.Author.FirstName,
			LastName:  entity.Author.LastName,
		},
		Categories: categories,
		Publisher: PublisherResponse{
			Id:   int(entity.Publisher.ID),
			Name: entity.Publisher.Name,
		},
		PublicationDate: publicationDate,
	}
}
