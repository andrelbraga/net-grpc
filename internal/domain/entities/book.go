package entities

import "time"

const (
	PrintTypeBook string = "BOOK"
)

type Books struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	BookID        *string `json:"book_id"`
	Language      string  `json:"language"`
	PublishedDate *string `json:"published_date"`
	PageCount     *int32  `json:"page_count"`
	Description   *string `json:"description"`
}

type BookAuthor struct {
	ID        int        `json:"id"`
	Author    string     `json:"author"`
	BookID    int        `json:"book_int_id"`
	CreatedAt *time.Time `json:"created_at"`
}

type BookImageLinks struct {
	ID             int        `json:"id"`
	BookID         int        `json:"book_int_id"`
	SmallThumbnail string     `json:"small_thumbnail"`
	Thumbnail      string     `json:"thumbnail"`
	CreatedAt      *time.Time `json:"created_at"`
}

// BookDetail
type BookDetail struct {
	ID            int            `json:"id"`
	Title         string         `json:"title"`
	Authors       []BookAuthor   `json:"authors"`
	ImageLinks    BookImageLinks `json:"image_links"`
	PrintType     string         `json:"print_type"`
	Language      string         `json:"language"`
	PublishedDate *string        `json:"published_date"`
	PageCount     *int32         `json:"page_count"`
	Description   *string        `json:"description"`
}
