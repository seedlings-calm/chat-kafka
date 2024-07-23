package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/seedlings-calm/chat-kafka/common"
	"github.com/seedlings-calm/chat-kafka/internal/proto/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/anypb"
)

func main() {
	//建立无认证的连接
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	client := types.NewChatServiceClient(conn)

	// go chat(client)
	// time.Sleep(5 * time.Second)
	broadcast(client)
}

// 模拟客户端发送数据和接受
func chat(client types.ChatServiceClient) {

	var wg sync.WaitGroup
	var err error
	streams := make([]types.ChatService_ChatClient, 5)
	for i := 0; i < 5; i++ {
		// Create a context with metadata
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
			"client-id", strconv.Itoa(i+1),
			"auth-token", "example-auth-token",
		))
		stream, err := client.Chat(ctx)
		if err != nil {
			panic(err)
		}
		streams[i] = stream
		wg.Add(1)
		go func(i int, stream types.ChatService_ChatClient) {
			defer wg.Done()
			for {
				msg, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Fatal(err)
				}
				log.Printf("我是 %d,收到:%s 的消息\n", i+1, msg.From)
			}
		}(i, stream)
	}

	for i, stream := range streams {
		content := &types.TextMsg{
			Value: fmt.Sprintf("hello world ! by %d", i+1),
		}
		msg := new(types.MsgType)
		msg.MsgType = "text"
		v, _ := anypb.New(content)
		msg.MsgContent = &types.AnyContent{
			Value: v,
		}
		// data := &types.ChatServiceRequest{
		// 	From:     strconv.Itoa(i + 1),
		// 	To:       "",
		// 	Msg:      []*types.MsgType{msg},
		// 	ChatType: common.Group,
		// }
		data := &types.ChatServiceRequest{
			From:     strconv.Itoa(i + 1),
			To:       "1",
			Msg:      []*types.MsgType{msg},
			ChatType: common.Private,
		}
		err = stream.Send(data)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("客户%d 发送了消息", i)
	}
	wg.Wait()
}

func broadcast(client types.ChatServiceClient) {

	var wg sync.WaitGroup
	var mu sync.RWMutex
	var err error
	streams := make([]types.ChatService_BroadcastClient, 5)
	for i := 0; i < 5; i++ {
		// Create a context with metadata
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
			"client-id", strconv.Itoa(i+1),
			"auth-token", "example-auth-token",
		))
		stream, err := client.Broadcast(ctx)
		if err != nil {
			panic(err)
		}
		mu.RLock()
		streams[i] = stream
		mu.RUnlock()
		wg.Add(1)
		go func(i int, stream types.ChatService_BroadcastClient) {
			defer wg.Done()
			for {
				msg, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Fatal(err)
				}
				log.Printf("我是 %d,收到:%s 的消息\n", i+1, msg.From)
			}
		}(i, stream)
	}

	for i, stream := range streams {
		content := &types.TextMsg{
			Value: fmt.Sprintf("hello world ! by %d", i+1),
		}
		msg := new(types.MsgType)
		msg.MsgType = "text"
		v, _ := anypb.New(content)
		msg.MsgContent = &types.AnyContent{
			Value: v,
		}

		data := &types.ChatServiceRequest{
			From:     "system",
			To:       strconv.Itoa(i + 1),
			Msg:      []*types.MsgType{msg},
			ChatType: common.Broadcast,
		}
		err = stream.Send(data)
		if err != nil {
			log.Fatal(err)
		}
	}
	wg.Wait()
}
