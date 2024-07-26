package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/seedlings-calm/chat-kafka/internal/constants"
	"github.com/seedlings-calm/chat-kafka/proto"
	"github.com/seedlings-calm/chat-kafka/proto/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	//建立无认证的连接
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	client := types.NewChatServiceClient(conn)

	go chat(client)
	select {}
}

// 模拟客户端发送数据和接受
func chat(client types.ChatServiceClient) {

	var wg sync.WaitGroup
	var err error
	var mu sync.RWMutex
	streams := make([]types.ChatService_ChatClient, 5)
	for i := 0; i < 5; i++ {
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
			"client-id", strconv.Itoa(i+1),
			"auth-token", "example-auth-token",
		))
		stream, err := client.Chat(ctx)
		if err != nil {
			panic(err)
		}
		mu.RLock()
		streams[i] = stream
		mu.RUnlock()
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
		se := &proto.ChatServiceReq{
			Req: &types.ChatServiceRequest{
				From:     strconv.Itoa(i + 1),
				To:       "1",
				ChatType: constants.Private,
			},
		}
		err = stream.Send(se.New("text", fmt.Sprintf("收到来自%d的消息", i)))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("客户%d 发送了消息", i)
	}
	wg.Wait()
}
