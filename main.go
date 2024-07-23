package main

import (
	"log"
	"net"

	interfaces_grpc "github.com/seedlings-calm/chat-kafka/internal/interfaces/grpc"
	"github.com/seedlings-calm/chat-kafka/internal/interfaces/repository"
	logic "github.com/seedlings-calm/chat-kafka/internal/logic/grpc"
	"github.com/seedlings-calm/chat-kafka/internal/proto/types"
	"google.golang.org/grpc"
)

func main() {
	// Initialize repositories and business logic
	repo := repository.NewMySQLChatRepository()
	chatLogic := logic.NewChatLogic(repo)
	chatService := interfaces_grpc.NewChatService(chatLogic)

	// Setup gRPC server
	grpcServer := grpc.NewServer()
	types.RegisterChatServiceServer(grpcServer, chatService)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
