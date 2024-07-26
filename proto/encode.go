package proto

import (
	"github.com/seedlings-calm/chat-kafka/proto/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ChatServiceReq struct {
	Req *types.ChatServiceRequest
}

func (req *ChatServiceReq) New(msgType, msgValue string) *types.ChatServiceRequest {
	msg := &types.MsgType{
		MsgType: msgType,
	}

	content := &types.TextMsg{
		Value: msgValue,
	}
	v, _ := anypb.New(content)
	msg.MsgContent = &types.AnyContent{
		Value: v,
	}

	msgN := make([]*types.MsgType, 0)
	msgN = append(msgN, msg)
	req.Req.Msg = msgN
	return req.Req
}

func (req *ChatServiceReq) Marshal() ([]byte, error) {
	return proto.Marshal(req.Req)
}

func (req *ChatServiceReq) Unmarshal(item []byte) (err error) {
	reqNew := &types.ChatServiceRequest{}
	err = proto.Unmarshal(item, reqNew)
	req.Req = reqNew
	return
}
