package grpc

import (
	"log"
	"net"

	pb "github.com/dagangilat/go-commerce/pkg/api/grpc"
	"google.golang.org/grpc"
)

func StartServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterECommerceServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}

// Implement your gRPC services here
