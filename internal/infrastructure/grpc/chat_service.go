package infrastructure_grpc

import (
	"log"

	"github.com/seedlings-calm/chat-kafka/internal/constants"
	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/kafka"
	"github.com/seedlings-calm/chat-kafka/proto/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 定义聊天服务实现proto编译的ChatServiceServer接口
type ChatServer struct {
	types.UnimplementedChatServiceServer
}

func (cs *ChatServer) Chat(stream types.ChatService_ChatServer) error {
	// Extract unique identifier from metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "missing metadata")
	}
	ids := md["client-id"]
	if len(ids) == 0 {
		return grpc.Errorf(codes.InvalidArgument, "missing client ID")
	}
	//TODO:发送方唯一标识， 用户接收消息的标识
	clientID := ids[0]
	sc := NewStreamChat()
	sc.SetClient(clientID, stream)

	for {
		// 接收消息
		reqMsg, err := stream.Recv()
		if err != nil {
			if status.Code(err) == codes.Canceled || status.Code(err) == codes.DeadlineExceeded {
				//可以做用户的连接，断开日志
				log.Printf("Client %s disconnected: %v", clientID, err)
			} else {
				//异常错误
				log.Printf("Error receiving message from client %s: %v", clientID, err)
			}
			return err
		}
		switch reqMsg.ChatType {
		case constants.Private:
			kafka.GetImProducer().PushPrivate(reqMsg)
		case constants.Group:
			kafka.GetImProducer().PushGroup(reqMsg)
		case constants.Broadcast:
			kafka.GetImProducer().PushBroadcast(reqMsg)
		default:
			log.Printf("未知的消息类型:%s", reqMsg.ChatType)
			continue
		}
		sc.Start()
	}
}
