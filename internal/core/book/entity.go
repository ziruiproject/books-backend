package book

import (
	"gorm.io/gorm"
	"starter/internal/core/author"
	"starter/internal/core/category"
	"starter/internal/core/publisher"
	"time"
)

type Book struct {
	Title           string
	Cover           string
	Description     string
	PageCount       int
	AuthorId        int
	Author          author.Author
	Categories      []category.Category `gorm:"many2many:book_category;"`
	PublisherId     int
	Publisher       publisher.Publisher
	PublicationDate time.Time
	gorm.Model
}
