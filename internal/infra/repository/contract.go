package repository

import (
	"context"

	"gorm.io/gorm"
	"net-grpc.com/internal/domain/entities"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) GetByID(bookID int) (*entities.BookDetail, error) {
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

func (r Repository) GetAllRandom(ctx context.Context) (*[]entities.Books, error) {
	var books = []entities.Books{}

	tx := r.db.Find(&books)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &books, nil
}
