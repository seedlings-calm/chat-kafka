package infrastructure_grpc

import (
	"github.com/seedlings-calm/chat-kafka/proto/types"
	"google.golang.org/grpc"
)

func NewGRPCServer() *grpc.Server {
	// 注册服务
	grpcServer := grpc.NewServer()

	//初始化实现了ChatService服务
	chatService := NewChatService()

	types.RegisterChatServiceServer(grpcServer, chatService)

	return grpcServer
}
