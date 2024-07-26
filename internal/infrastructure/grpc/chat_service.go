package infrastructure_grpc

import (
	"log"
	"sync"
	"time"

	"github.com/seedlings-calm/chat-kafka/internal/constants"
	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/kafka"
	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/redis"
	"github.com/seedlings-calm/chat-kafka/proto/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// 定义聊天服务实现proto编译的ChatServiceServer接口
type ChatServer struct {
	types.UnimplementedChatServiceServer
	clients  sync.Map          //客户端流存储
	consumer *kafka.ImConsumer //kafka消费者实例
}

func NewChatService() *ChatServer {
	return &ChatServer{
		clients:  sync.Map{},
		consumer: kafka.NewImConsumer(),
	}
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
	clientID := ids[0]
	//通讯客户端存储
	cs.clients.Store(clientID, stream)
	defer cs.clients.Delete(clientID)

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
			kafka.New().PushPrivate(reqMsg)
		case constants.Group:
			kafka.New().PushGroup(reqMsg)
		case constants.Broadcast:
			kafka.New().PushBroadcast(reqMsg)
		default:
			log.Printf("未知的消息类型:%s", reqMsg.ChatType)
			continue
		}
		cs.privateMsg(stream)
		go cs.groupMsg(stream)
		go cs.broadcastMsg(stream)
	}
}

// privateMsg 从 Kafka 消费消息并通过 gRPC 推送到客户端
func (cs *ChatServer) privateMsg(stream types.ChatService_ChatServer) {
	go func() {

		for {
			consumer := cs.consumer.GetPrivateConsumer()
			msg, err := consumer.Pull()
			if err != nil {
				log.Printf("Failed to pull message: %v", err)
				continue
			}

			req := &types.ChatServiceRequest{}
			err = proto.Unmarshal([]byte(msg.Value), req)
			if err != nil {
				log.Println("proto unmarshal err:", err)
				continue
			}
			response := &types.ChatServiceResponse{
				From:           req.From,
				To:             []string{req.To},
				Msg:            req.GetMsg(),
				ChatType:       "private", // or whatever type is appropriate
				RoomId:         "",
				FromClientType: req.FromClientType, // or whatever is appropriate
				ToClientType:   "",                 // or whatever is appropriate
			}

			for _, v := range response.To {
				cs.clients.Range(func(key, value interface{}) bool {
					if key == v {
						stream, ok := value.(types.ChatService_ChatServer)
						if !ok {
							log.Println("获取客户端流 出错")
							return true
						}
						err := stream.Send(response)
						if err != nil {
							log.Println("send fail:", err)
							return true
						}

						log.Printf("%s 的消息发送给了 %s", response.From, v)
					}
					return true
				})
			}
			redis.GetChatRedis().Set(stream.Context(), redis.ChatPrivateKey(response.From), response, 7*24*time.Hour)

			consumer.Ack(msg.Session)
		}
	}()

}

// groupMsg 从 Kafka 消费消息并通过 gRPC 推送到客户端
func (cs *ChatServer) groupMsg(stream types.ChatService_ChatServer) {
	for {
		consumer := cs.consumer.GetGroupConsumer()
		msg, err := consumer.Pull()
		if err != nil {
			log.Printf("Failed to pull message: %v", err)
			continue
		}

		req := &types.ChatServiceRequest{}
		err = proto.Unmarshal([]byte(msg.Value), req)
		if err != nil {
			log.Println("proto unmarshal err:", err)
			continue
		}
		//TODO:
		//群组聊天接收者，需要根据request.To 字符串解析出群组成员，赋值到response.To里面

		// 将 Kafka 消息发送回客户端
		response := &types.ChatServiceResponse{
			From:           req.From,
			To:             []string{"1", "2", "3"},
			Msg:            req.GetMsg(),
			ChatType:       "group", // or whatever type is appropriate
			RoomId:         req.GetRoomId(),
			FromClientType: req.FromClientType, // or whatever is appropriate
			ToClientType:   "",                 // or whatever is appropriate
		}

		for _, v := range response.To {
			cs.clients.Range(func(key, value interface{}) bool {
				if key == v {
					stream, ok := value.(types.ChatService_ChatServer)
					if !ok {
						log.Println("获取客户端流 出错")
						return true
					}
					err := stream.Send(response)
					if err != nil {
						log.Println("send fail:", err)
						return true
					}
					log.Printf("%s 的群组消息，发送给了 %s", response.From, v)
				}
				return true
			})
		}
	}
}

// broadcastMsg 从 Kafka 消费消息并通过 gRPC 推送到客户端
func (cs *ChatServer) broadcastMsg(stream types.ChatService_ChatServer) {
	for {
		consumer := cs.consumer.GetBroadcastConsumer()
		msg, err := consumer.Pull()
		if err != nil {
			log.Printf("Failed to pull message: %v", err)
			continue
		}

		req := &types.ChatServiceRequest{}
		err = proto.Unmarshal([]byte(msg.Value), req)
		if err != nil {
			log.Println("proto unmarshal err:", err)
			continue
		}
		//TODO:
		//广播信息，发送给所有在线用户，也可以定义发送用户群体

		// 将 Kafka 消息发送回客户端
		response := &types.ChatServiceResponse{
			From:           req.From,
			To:             []string{"2", "3"},
			Msg:            req.GetMsg(),
			ChatType:       "broadcast", // or whatever type is appropriate
			RoomId:         "",
			FromClientType: req.FromClientType, // or whatever is appropriate
			ToClientType:   req.ToClientType,   // or whatever is appropriate
		}

		for _, v := range response.To {
			cs.clients.Range(func(key, value interface{}) bool {
				if key == v {
					stream, ok := value.(types.ChatService_ChatServer)
					if !ok {
						log.Println("获取客户端流 出错")
						return true
					}
					err := stream.Send(response)
					if err != nil {
						log.Println("send fail:", err)
						return true
					}
					log.Printf("%s 的广播消息，广播给了 %s", response.From, v)
				}
				return true
			})
		}
	}
}
