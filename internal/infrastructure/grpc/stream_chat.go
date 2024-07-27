package infrastructure_grpc

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/seedlings-calm/chat-kafka/internal/constants"
	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/kafka"
	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/redis"
	"github.com/seedlings-calm/chat-kafka/proto/types"
	"google.golang.org/protobuf/proto"
)

var steamOnce sync.Once

type StreamChat struct {
	clients sync.Map //客户端流存储
}

// 使用单例设计模式，实现kafka信息和grpc交互
func NewStreamChat() (sc *StreamChat) {
	steamOnce.Do(func() {
		sc = &StreamChat{
			clients: sync.Map{},
		}
	})
	return
}

func (sc *StreamChat) SetClient(it string, stream types.ChatService_ChatServer) {
	sc.clients.Store(it, stream)
}

// processMessage 处理从 Kafka 拉取的消息
func (sc *StreamChat) processMessage(msgType string, consumer *kafka.Kafka) {
	for {
		msg, err := consumer.Pull()
		if err != nil {
			log.Printf("Failed to pull %s message: %v", msgType, err)
			continue
		}

		req := &types.ChatServiceRequest{}
		if err := proto.Unmarshal([]byte(msg.Value), req); err != nil {
			log.Println("proto unmarshal error:", err)
			continue
		}

		response := &types.ChatServiceResponse{
			From:           req.From,
			Msg:            req.GetMsg(),
			ChatType:       msgType,
			RoomId:         req.GetRoomId(),
			FromClientType: req.FromClientType,
			ToClientType:   req.ToClientType,
		}

		sc.sendMessageToClients(response)
		consumer.Ack(msg.Session)
		redis.GetChatRedis().Set(context.Background(), redis.ChatPrivateKey(response.From), response, 7*24*time.Hour)
	}
}

// sendMessageToClients 将消息发送到所有客户端
func (sc *StreamChat) sendMessageToClients(response *types.ChatServiceResponse) {
	sc.clients.Range(func(key, value interface{}) bool {
		clientID, ok := key.(string)
		if !ok {
			log.Println("Error casting client ID")
			return true
		}

		stream, ok := value.(types.ChatService_ChatServer)
		if !ok {
			log.Println("Error casting stream")
			return true
		}

		if contains(response.To, clientID) {
			if err := stream.Send(response); err != nil {
				log.Println("Send failed:", err)
				return true
			}
			log.Printf("Message from %s sent to %s", response.From, clientID)
		}
		return true
	})
}

// contains 检查切片中是否包含指定值
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (sc *StreamChat) Start() {
	go sc.processMessage(constants.Private, kafka.GetImConsumer().GetPrivateConsumer())
	go sc.processMessage(constants.Group, kafka.GetImConsumer().GetGroupConsumer())
	go sc.processMessage(constants.Broadcast, kafka.GetImConsumer().GetBroadcastConsumer())
}
