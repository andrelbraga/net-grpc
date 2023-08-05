package repository

import (
	"gorm.io/gorm"
	"net-grpc.com/internal/domain/entities"
)

// BookRepository
type BookRepository struct {
	db *gorm.DB
}

// NewBookRepository
func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

// GetBookDetailByID
func (r *BookRepository) GetBookDetailByID(bookID int) (*entities.BookDetail, error) {
	var books = entities.Books{ID: bookID}

	tx := r.db.Find(&books)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var authors = []entities.BookAuthor{}
	tx = r.db.Find(&authors, "book_int_id = ?", bookID)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var bookImageLinks = entities.BookImageLinks{BookID: bookID}
	tx = r.db.Find(&bookImageLinks)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &entities.BookDetail{
		ID:            books.ID,
		Title:         books.Title,
		Authors:       authors,
		ImageLinks:    bookImageLinks,
		PrintType:     entities.PrintTypeBook,
		Language:      books.Language,
		PublishedDate: books.PublishedDate,
		PageCount:     books.PageCount,
		Description:   books.Description,
	}, nil
}

// GetBooksByIDs
func (r *BookRepository) GetBooksByIDs(randomIds []int) ([]entities.Books, error) {
	var books = []entities.Books{}

	tx := r.db.Where(randomIds).Find(&books)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return books, nil
}
