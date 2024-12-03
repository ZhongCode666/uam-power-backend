package dbservice

import (
	"context"
	"github.com/segmentio/kafka-go"
)

// KafkaProducer 封装 Kafka 生产者
type KafkaProducer struct {
	// Kafka 生产者写入器
	writer *kafka.Writer
}

// NewKafkaProducer 创建一个新的 Kafka 生产者
// addr 是 Kafka 代理地址
// topic 是 Kafka 主题
func NewKafkaProducer(addr, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(addr),     // Kafka 代理地址
		Topic:    topic,               // Kafka 主题
		Balancer: &kafka.LeastBytes{}, // 使用最少字节分配策略
		Async:    true,                // 异步写入
	}
	return &KafkaProducer{writer: writer} // 返回 KafkaProducer 实例
}

// SendMessage 发送消息到 Kafka
// message 是要发送的消息字符串
func (p *KafkaProducer) SendMessage(message string) error {
	msg := kafka.Message{
		Value: []byte(message), // 将消息字符串转换为字节数组
	}
	return p.writer.WriteMessages(context.Background(), msg) // 发送消息到 Kafka
}

// Close 关闭 Kafka 生产者
// 关闭 Kafka 生产者写入器并释放相关资源
func (p *KafkaProducer) Close() error {
	return p.writer.Close() // 关闭写入器
}
