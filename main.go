package main

import (
	"log"
	"net"

	infrastructure_grpc "github.com/seedlings-calm/chat-kafka/internal/infrastructure/grpc"
)

func main() {
	// Initialize repositories and business logic
	grpcServer := infrastructure_grpc.NewGRPCServer()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
