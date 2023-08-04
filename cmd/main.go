package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "net-grpc.com/internal/grpc/proto"
	"net-grpc.com/internal/infra/repository"
	"net-grpc.com/internal/service"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5001))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	con, err := repository.NewConnectDB()
	if err != nil {
		panic(err)
	}
	var repo = repository.NewRepository(con)
	var bsrv = service.NewService(repo)

	pb.RegisterPrivateBookServiceServer(s, bsrv)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
