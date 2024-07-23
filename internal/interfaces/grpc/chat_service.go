package grpc

import (
	"log"
	"sync"

	logic "github.com/seedlings-calm/chat-kafka/internal/logic/grpc"
	"github.com/seedlings-calm/chat-kafka/internal/proto/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 定义聊天服务实现proto编译的ChatServiceServer接口
type ChatServer struct {
	types.UnimplementedChatServiceServer
	// clients map[string]types.ChatService_ChatServer
	clients          sync.Map
	boardcastClients sync.Map
	//消息处理逻辑
	logic logic.ChatLogic
}

func NewChatService(logic logic.ChatLogic) *ChatServer {
	return &ChatServer{
		logic: logic,
		// clients:           make(map[string]types.ChatService_ChatServer),
		clients:          sync.Map{},
		boardcastClients: sync.Map{},
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

	cs.clients.Store(clientID, stream)
	defer cs.clients.Delete(clientID)
	for {
		// 接收消息
		reqMsg, err := stream.Recv()
		if err != nil {
			if status.Code(err) == codes.Canceled || status.Code(err) == codes.DeadlineExceeded {
				log.Printf("Client %s disconnected: %v", clientID, err)
			} else {
				log.Printf("Error receiving message from client %s: %v", clientID, err)
			}
			return err
		}

		//调用消息处理中间件处理接收到的消息，然后返回处理结果，根据结果把消息推送给特定用户
		resMsg := cs.logic.HandleChatMessage(reqMsg)
		err = cs.SendChatMessage(resMsg)
		if err != nil {
			log.Println("发送失败:", err.Error())
			return err
		}

	}
}
func (cs *ChatServer) Broadcast(stream types.ChatService_BroadcastServer) error {
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

	cs.boardcastClients.Store(clientID, stream)
	defer cs.boardcastClients.Delete(clientID)

	for {
		// 接收消息
		reqMsg, err := stream.Recv()
		if err != nil {
			if status.Code(err) == codes.Canceled || status.Code(err) == codes.DeadlineExceeded {
				log.Printf("Client %s disconnected: %v", clientID, err)
			} else {
				log.Printf("Error receiving message from client %s: %v", clientID, err)
			}
			return err
		}

		//调用消息处理中间件处理接收到的消息，然后返回处理结果，根据结果把消息推送给特定用户
		resMsg := cs.logic.HandleBoardcast(reqMsg)
		err = cs.SendBoardcastMessage(resMsg)
		if err != nil {
			log.Println("发送失败:", err.Error())
			return err
		}

	}
}

func (cs *ChatServer) SendChatMessage(msg *types.ChatServiceResponse) error {

	log.Println("接收方::", msg.To)

	for _, v := range msg.To {
		cs.clients.Range(func(key, value interface{}) bool {
			if key == v {
				stream, ok := value.(types.ChatService_ChatServer)
				if !ok {
					log.Println("获取客户端流 出错")
					return true
				}
				err := stream.Send(msg)
				if err != nil {
					log.Println("send fail:", err)
					return true
				}
				log.Printf("%s 的消息发送给了 %s", msg.From, v)
			}
			return true
		})
	}
	return nil
}

func (cs *ChatServer) SendBoardcastMessage(msg *types.ChatServiceResponse) error {
	cs.boardcastClients.Range(func(key, value interface{}) bool {
		stream, ok := value.(types.ChatService_BroadcastServer)
		if !ok {
			log.Println("获取广播流错误")
			return true
		}
		stream.Send(msg)
		log.Printf("广播消息给 %s", key)
		return true
	})
	return nil
}