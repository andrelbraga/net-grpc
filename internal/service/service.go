package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"net-grpc.com/internal/domain/entities"
	pb "net-grpc.com/internal/grpc/proto"
	"net-grpc.com/internal/infra/repository"
)

type BookService struct {
	pb.UnimplementedPrivateBookServiceServer
	repo repository.Repository
}

func NewService(repo repository.Repository) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (srv *BookService) GetRandomBook(emp *emptypb.Empty, stream pb.PrivateBookService_GetRandomBookServer) error {
	var booksDiplayed []entities.Books
	var currentBook = &pb.GetDisplayBooksResponse{}
	var lastBook = &pb.Book{}

	books, err := srv.repo.GetAllRandom()
	if err != nil {
		return err
	}

	for idx, book := range books {
		log.Printf("book index: %d", idx)

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

func (srv *BookService) GetBookDetail(ctx context.Context, req *pb.GetBookDetailsRequest) (*pb.GetBookDetailsResponse, error) {
	var listAuthors []string

	bookId, err := strconv.Atoi(req.BookId)
	if err != nil {
		return nil, err
	}

	log.Printf("searching book-detail by id %d", bookId)
	result, err := srv.repo.GetByID(bookId)
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
