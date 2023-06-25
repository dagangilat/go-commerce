package grpc

import (
	"log"
	"net"

	// Update this import to match the go_package option in your .proto file
	// Replace this with the correct package for your generated code

	"google.golang.org/grpc"
)

type server struct{}

func StartServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	grpc.RegisterCustomerServiceServer(s, &server{}) // Use the imported package
	grpc.RegisterProductServiceServer(s, &server{})  // Use the imported package

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Implement your gRPC services here
