package service_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net-grpc.com/internal/domain/entities"
	pb "net-grpc.com/internal/infra/grpc"
	"net-grpc.com/internal/infra/repository"
	"net-grpc.com/internal/service"
)

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) GetRandomBook(randomIds []int) ([]entities.Books, error) {
	args := m.Called(randomIds)
	return args.Get(0).([]entities.Books), args.Error(1)
}

func (m *MockBookRepository) GetBookDetail(bookID int) (*entities.BookDetail, error) {
	args := m.Called(bookID)
	return args.Get(0).(*entities.BookDetail), args.Error(1)
}

type BookServiceStub struct {
	suite.Suite
	repo    *repository.BookRepository
	service *service.BookService
}

func TestBookTestSuite(t *testing.T) {
	suite.Run(t, new(BookServiceStub))
}

func (s *BookServiceStub) SetupSuite() {
	con, err := repository.NewConnectDB()
	if err != nil {
		panic("connect db")
	}
	//var server = grpc.NewServer()
	var repo = repository.NewBookRepository(con)
	s.repo = repo
	s.service = service.NewBookService(repo)
}

func (srv *BookServiceStub) TestGetbookDetailByID() {
	ctx := context.Background()
	expectedBookID := "1"
	expectedResult := &entities.BookDetail{ID: 1, Title: "Test Book"}

	param := &pb.GetBookDetailsRequest{
		BookId: expectedBookID,
		ApiKey: "api_key",
	}

	result, _ := srv.service.GetBookDetail(ctx, param)
	srv.Equal(strconv.Itoa(expectedResult.ID), result.Book.Id)
}
