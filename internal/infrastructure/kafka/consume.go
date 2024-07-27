package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/seedlings-calm/chat-kafka/proto/types"
	"google.golang.org/protobuf/proto"
)

// 实现sarama的ConsumerGroupHandler 接口
type consume struct {
	C       chan ConsumerMessage //通道
	session sarama.ConsumerGroupSession
	quit    chan int
	Close   chan int
}

// 方法在开始时运行
func (c *consume) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// 方法在结束时运行
func (c *consume) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// 消费循环
func (c *consume) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.session = session

	for {
		select {
		case <-c.Close:
			log.Printf("kafka收到退出信号")
			return nil
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			log.Printf("Message claimed: key = %s, topic = %s,offset = %d ", string(message.Key), message.Topic, message.Offset)
			r := &types.ChatServiceRequest{}
			proto.Unmarshal(message.Value, r)
			log.Println(r)
			if c.C != nil {
				c.C <- ConsumerMessage{
					Value:     message.Value,
					Timestamp: message.Timestamp,
					Key:       message.Key,
					Topic:     message.Topic,
					Partition: message.Partition,
					Offset:    message.Offset,
					Session:   message,
				}
			} else {
				log.Println("管道已关闭")
			}
		}
	}
}
