package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
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

	books, err := srv.repo.GetAllRandom()
	if err != nil {
		return err
	}

	for idx, book := range books {
		log.Printf("book index: %d", idx)
		var bookDisplayed = &pb.GetDisplayBooksResponse{
			RandomBook: &pb.Book{
				Id:        strconv.Itoa(book.ID),
				Title:     book.Title,
				Authors:   []string{},
				PrintType: "BOOK",
				Language:  book.Language,
			},
			LastBook: nil,
		}
		time.Sleep(5 * time.Second)
		if err := stream.Send(bookDisplayed); err != nil {
			return err
		}
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
