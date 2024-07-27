package kafka

import (
	"log"
	"sync"

	"github.com/seedlings-calm/chat-kafka/internal/constants"
	"github.com/seedlings-calm/chat-kafka/proto/types"
	"google.golang.org/protobuf/proto"
)

func init() {
	newImProducer()
}

var (
	imProduceOnce sync.Once
	imProducer    *ImProducer
	addr          = []string{"127.0.0.1:9194", "127.0.0.1:9195", "127.0.0.1:9196"}
)

type ImProducer struct {
	private   *Kafka //私聊
	groups    *Kafka //群聊
	broadcast *Kafka //广播

}

func GetImProducer() *ImProducer {
	return imProducer
}

func newImProducer() *ImProducer {
	imProduceOnce.Do(func() {
		private, err := NewSyncProducer(KafkaConf{
			Addr:  addr,
			Topic: constants.Private,
		})
		if err != nil {
			panic(err)
		}
		group, err := NewSyncProducer(KafkaConf{
			Addr:  addr,
			Topic: constants.Group,
		})
		if err != nil {
			panic(err)
		}
		broadcast, err := NewAsyncProducer(KafkaConf{
			Addr:  addr,
			Topic: constants.Broadcast,
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
