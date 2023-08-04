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
	log.Printf("GetRandomBook")

	for _, i := range []string{"0", "1", "2", "3"} {
		log.Printf("Item %s", i)
		var bookDisplayed = &pb.GetDisplayBooksResponse{
			RandomBook: nil,
			LastBook:   nil,
		}
		time.Sleep(10 * time.Second)
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
