package dbservice

import (
	"context"
	"github.com/segmentio/kafka-go"
)

// KafkaConsumer 封装 Kafka 消费者
type KafkaConsumer struct {
	reader *kafka.Reader
}

// NewKafkaConsumer 创建一个新的 Kafka 消费者
func NewKafkaConsumer(addr, topic, groupID string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{addr},
		Topic:   topic,
		GroupID: groupID,
	})
	return &KafkaConsumer{reader: reader}
}

// ReceiveMessage 从 Kafka 中接收消息
func (c *KafkaConsumer) ReceiveMessage(ctx context.Context) (string, error) {
	msg, err := c.reader.ReadMessage(ctx)
	//var re map[string]interface{} map[string]interface{}
	if err != nil {
		return "", err
	}
	//err = json.Unmarshal(msg.Value, &re)
	return string(msg.Value), nil
}

// Close 关闭 Kafka 消费者
func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
