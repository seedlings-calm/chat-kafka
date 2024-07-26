package main

import (
	"flag"
	"log"
	"net"

	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/config"
	infrastructure_grpc "github.com/seedlings-calm/chat-kafka/internal/infrastructure/grpc"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	var cfgPath string
	// 定义命令行参数
	flag.StringVar(&cfgPath, "c", "./config/config.yml", "choose config file.")
	flag.Parse()

	config.Setup(cfgPath)

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
