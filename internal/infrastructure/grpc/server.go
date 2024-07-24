package infrastructure_grpc

import (
	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/mysql"
	"github.com/seedlings-calm/chat-kafka/internal/usecase"
	"github.com/seedlings-calm/chat-kafka/proto/types"
	"google.golang.org/grpc"
)

func NewGRPCServer() *grpc.Server {
	// Initialize repositories and business logic
	repo := mysql.NewMySQLChatRepository()
	chatLogic := usecase.NewChatLogic(repo)
	chatService := NewChatService(chatLogic)

	// 注册服务
	grpcServer := grpc.NewServer()

	types.RegisterChatServiceServer(grpcServer, chatService)

	return grpcServer
}
