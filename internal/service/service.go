package service

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"net-grpc.com/internal/domain/entities"
	pb "net-grpc.com/internal/infra/grpc"
	"net-grpc.com/internal/infra/repository"
)

// BookService
type BookService struct {
	pb.UnimplementedPrivateBookServiceServer
	repo *repository.BookRepository
}

// NewBookService
func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

// GetRandomBook
func (srv *BookService) GetRandomBook(req *pb.GetBookRandomRequest, stream pb.PrivateBookService_GetRandomBookServer) error {
	var booksDiplayed []entities.Books
	var currentBook = &pb.GetDisplayBooksResponse{}
	var lastBook = &pb.Book{}
	var itemsPerRequest = 10
	var totalItems = 40 //Por eu ja ter o numero de itens no banco resolvi deixar setado aqui mas poderia incluir um 'select count(*)'

	randomIds := getRandomBooks(itemsPerRequest, totalItems)
	log.Printf("random ids: %d", randomIds)
	books, err := srv.repo.GetBooksByIDs(randomIds)
	if err != nil {
		return err
	}

	for idx, book := range books {
		log.Printf("index: %d bookId: %d", idx, book.ID)

		if len(booksDiplayed) > 0 {
			var last = booksDiplayed[len(booksDiplayed)-1]
			lastBook = &pb.Book{
				Id:        strconv.Itoa(last.ID),
				Title:     last.Title,
				Authors:   []string{},
				PrintType: entities.PrintTypeBook,
				Language:  last.Language,
			}
		}

		currentBook = &pb.GetDisplayBooksResponse{
			RandomBook: &pb.Book{
				Id:        strconv.Itoa(book.ID),
				Title:     book.Title,
				Authors:   []string{},
				PrintType: entities.PrintTypeBook,
				Language:  book.Language,
			},
			LastBook: lastBook,
		}
		time.Sleep(5 * time.Second)
		if err := stream.Send(currentBook); err != nil {
			return err
		}
		booksDiplayed = append(booksDiplayed, book)

	}
	return nil
}

// GetBookDetail
func (srv *BookService) GetBookDetail(ctx context.Context, req *pb.GetBookDetailsRequest) (*pb.GetBookDetailsResponse, error) {
	var listAuthors []string

	bookId, err := strconv.Atoi(req.BookId)
	if err != nil {
		return nil, err
	}

	log.Printf("searching book-detail by id %d", bookId)
	result, err := srv.repo.GetBookDetailByID(bookId)
	if err != nil {
		return nil, err
	}

	if len(result.Authors) > 0 {
		for _, author := range result.Authors {
			listAuthors = append(listAuthors, author.Author)
		}
	}

	return &pb.GetBookDetailsResponse{
		BookId: req.BookId,
		Book: &pb.Book{
			Id:            strconv.Itoa(result.ID),
			Title:         result.Title,
			Authors:       listAuthors,
			ImageLinks:    &pb.BookImageLinks{SmallThumbnail: result.ImageLinks.SmallThumbnail, Thumbnail: result.ImageLinks.Thumbnail},
			PrintType:     result.PrintType,
			Language:      result.Language,
			PublishedDate: nil,
			PageCount:     result.PageCount,
			Description:   result.Description,
		},
	}, nil
}

func getRandomBooks(itemsPerRequest int, totalItems int) []int {

	// Gera uma base apartir dos milesegundos atuais para poder gerar numeros aleatorios apartir dai
	source := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(source)

	selectedItems := make(map[int]bool)

	for len(selectedItems) < itemsPerRequest {
		item := rnd.Intn(totalItems) + 1
		selectedItems[item] = true
	}

	var randomIds []int
	for item := range selectedItems {
		randomIds = append(randomIds, item)
	}

	return randomIds
}
