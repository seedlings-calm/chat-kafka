package kafka

import (
	"log"
	"sync"

	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/config"
	"github.com/seedlings-calm/chat-kafka/proto/types"
	"google.golang.org/protobuf/proto"
)

var (
	imProduceOnce sync.Once
	imProducer    *ImProducer
)

type ImProducer struct {
	private   *Kafka //私聊
	groups    *Kafka //群聊
	broadcast *Kafka //广播

}

func GetImProducer() *ImProducer {
	return newImProducer()
}

func newImProducer() *ImProducer {
	imProduceOnce.Do(func() {
		private, err := NewSyncProducer(KafkaConf{
			Addr:  config.GetGlobalConf().Kafka.Private.Addr,
			Topic: config.GetGlobalConf().Kafka.Private.Topic,
		})
		if err != nil {
			panic(err)
		}
		group, err := NewSyncProducer(KafkaConf{
			Addr:  config.GetGlobalConf().Kafka.Group.Addr,
			Topic: config.GetGlobalConf().Kafka.Group.Topic,
		})
		if err != nil {
			panic(err)
		}
		broadcast, err := NewAsyncProducer(KafkaConf{
			Addr:  config.GetGlobalConf().Kafka.Broadcast.Addr,
			Topic: config.GetGlobalConf().Kafka.Broadcast.Topic,
		})
		if err != nil {
			panic(err)
		}
		imProducer = &ImProducer{
			private:   private,
			groups:    group,
			broadcast: broadcast,
		}
	})
	return imProducer
}

func (i *ImProducer) PushPrivate(msg ...*types.ChatServiceRequest) {
	if len(msg) == 0 {
		return
	}
	log.Println("进入私聊，根据grpc信息，推送kafka")
	var messages []string
	for _, m := range msg {
		msgByte, err := proto.Marshal(m)
		if err != nil {
			log.Println("marshal err：", err)
			return
		}

		messages = append(messages, string(msgByte))
	}

	// Send the messages
	err := i.private.SyncPush(messages)
	log.Println("推送kafka,", err)
}

func (i *ImProducer) PushGroup(msg ...*types.ChatServiceRequest) {
	if len(msg) == 0 {
		return
	}
	log.Println("进入群组，根据grpc信息，推送kafka")
	var messages []string
	for _, m := range msg {
		msgByte, err := proto.Marshal(m)
		if err != nil {
			log.Println("marshal err：", err)
			return
		}

		messages = append(messages, string(msgByte))
	}

	// Send the messages
	err := i.groups.SyncPush(messages)
	log.Println("推送kafka,", err)
}

func (i *ImProducer) PushBroadcast(msg ...*types.ChatServiceRequest) {
	if len(msg) == 0 {
		return
	}
	log.Println("进入群组，根据grpc信息，推送kafka")
	var messages []string
	for _, m := range msg {
		msgByte, err := proto.Marshal(m)
		if err != nil {
			log.Println("marshal err：", err)
			return
		}

		messages = append(messages, string(msgByte))
	}

	// Send the messages
	i.broadcast.AsyncPush(messages)
	log.Println("推送kafka,")
}
