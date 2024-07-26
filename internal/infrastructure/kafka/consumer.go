package kafka

import (
	"context"
	"sync"

	"github.com/seedlings-calm/chat-kafka/internal/constants"
)

var (
	imConsumerOnce sync.Once
	imConsumer     *ImConsumer
)

type ImConsumer struct {
	private   *Kafka //私聊
	groups    *Kafka //群聊
	broadcast *Kafka //广播

}

func NewImConsumer() *ImConsumer {
	imConsumerOnce.Do(func() {
		private, err := NewConsumer(KafkaConf{
			Addr:    addr,
			Group:   "private-group",
			Topic:   constants.Private,
			Context: context.Background(),
		})
		if err != nil {
			panic(err)
		}
		group, err := NewConsumer(KafkaConf{
			Addr:    addr,
			Group:   "group-group",
			Topic:   constants.Group,
			Context: context.Background(),
		})
		if err != nil {
			panic(err)
		}
		broadcast, err := NewConsumer(KafkaConf{
			Addr:    addr,
			Group:   "broadcast-group",
			Topic:   constants.Broadcast,
			Context: context.Background(),
		})
		if err != nil {
			panic(err)
		}
		imConsumer = &ImConsumer{
			private:   private,
			groups:    group,
			broadcast: broadcast,
		}
	})
	return imConsumer
}

func (im *ImConsumer) GetPrivateConsumer() *Kafka {
	return im.private
}

func (im *ImConsumer) GetGroupConsumer() *Kafka {
	return im.groups
}

func (im *ImConsumer) GetBroadcastConsumer() *Kafka {
	return im.broadcast
}
