package usecase

import (
	"log"

	"github.com/seedlings-calm/chat-kafka/internal/constants"
	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/mysql"
	"github.com/seedlings-calm/chat-kafka/proto/types"
)

type ChatLogic interface {
	HandleChatMessage(req *types.ChatServiceRequest) *types.ChatServiceResponse
	HandleBroadcast(req *types.ChatServiceRequest) *types.ChatServiceResponse
}

type chatLogic struct {
	repo mysql.ChatRepository
}

func NewChatLogic(repo mysql.ChatRepository) ChatLogic {
	return &chatLogic{
		repo: repo,
	}
}

// TODO: 解决不在线问题，和推送kafka问题
func (uc *chatLogic) HandleChatMessage(req *types.ChatServiceRequest) *types.ChatServiceResponse {
	log.Println("进入处理客户端消息")
	// 处理接收到的聊天消息，根据 chat_type 来判断是私聊还是群聊
	response := &types.ChatServiceResponse{
		Msg:  req.GetMsg(),
		From: req.From,
		To:   make([]string, 0),
	}

	switch req.ChatType {
	case constants.Private:
		// 私聊，发送给指定用户
		response.To = append(response.To, req.To)
	case constants.Group:
		// 群聊，发送给群组中的每个用户
		users := uc.repo.GetGroupUsers(req.To)
		response.To = append(response.To, users...)
	}
	return response
}

func (uc *chatLogic) HandleBroadcast(req *types.ChatServiceRequest) *types.ChatServiceResponse {
	log.Println("进入服务端广播消息处理")
	response := &types.ChatServiceResponse{
		Msg:  req.GetMsg(),
		From: req.From,
		To:   make([]string, 0),
	}
	return response
}