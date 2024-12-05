package dbservice

import (
	"context"
	"github.com/segmentio/kafka-go"
	"uam-power-backend/utils"
)

// KafkaConsumer 封装 Kafka 消费者
type KafkaConsumer struct {
	// Kafka 消费者读取器
	reader *kafka.Reader
}

// NewKafkaConsumer 创建一个新的 Kafka 消费者
// addr 是 Kafka 代理地址
// topic 是 Kafka 主题
// groupID 是 Kafka 消费者组 ID
func NewKafkaConsumer(addr, topic, groupID string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{addr}, // Kafka 代理地址
		Topic:   topic,          // Kafka 主题
		GroupID: groupID,        // Kafka 消费者组 ID
	})
	return &KafkaConsumer{reader: reader} // 返回 KafkaConsumer 实例
}

// ReceiveMessage 从 Kafka 中接收消息
// ctx 是上下文，用于控制取消和超时
// 返回接收到的消息字符串和可能的错误
func (c *KafkaConsumer) ReceiveMessage(ctx context.Context) (string, error) {
	msg, err := c.reader.ReadMessage(ctx) // 从 Kafka 读取消息
	if err != nil {
		utils.MsgError("        [KafkaReceiveMessage]read message failed: >" + err.Error())
		return "", err // 如果读取失败，返回错误
	}
	return string(msg.Value), nil // 返回消息的值（字符串形式）
}

// Close 关闭 Kafka 消费者
// 关闭 Kafka 消费者读取器并释放相关资源
func (c *KafkaConsumer) Close() error {
	return c.reader.Close() // 关闭读取器
}
