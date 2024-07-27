package kafka

import (
	"context"
	"sync"

	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/config"
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

func GetImConsumer() *ImConsumer {
	return newImConsumer()
}

func newImConsumer() *ImConsumer {
	imConsumerOnce.Do(func() {
		private, err := NewConsumer(KafkaConf{
			Addr:    config.GetGlobalConf().Kafka.Private.Addr,
			Group:   config.GetGlobalConf().Kafka.Private.Group,
			Topic:   config.GetGlobalConf().Kafka.Private.Topic,
			Context: context.Background(),
		})
		if err != nil {
			panic(err)
		}
		group, err := NewConsumer(KafkaConf{
			Addr:    config.GetGlobalConf().Kafka.Group.Addr,
			Group:   config.GetGlobalConf().Kafka.Group.Group,
			Topic:   config.GetGlobalConf().Kafka.Group.Topic,
			Context: context.Background(),
		})
		if err != nil {
			panic(err)
		}
		broadcast, err := NewConsumer(KafkaConf{
			Addr:    config.GetGlobalConf().Kafka.Broadcast.Addr,
			Group:   config.GetGlobalConf().Kafka.Broadcast.Group,
			Topic:   config.GetGlobalConf().Kafka.Broadcast.Topic,
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
