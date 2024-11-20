package dbservice

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

// KafkaProducer 封装 Kafka 生产者
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer 创建一个新的 Kafka 生产者
func NewKafkaProducer(addr, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(addr),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		Async:        true,
		BatchSize:    32768, // 批次大小（消息数）
		BatchTimeout: 100 * time.Millisecond,
	}
	return &KafkaProducer{writer: writer}
}

// SendMessage 发送消息到 Kafka
func (p *KafkaProducer) SendMessage(message string) error {
	msg := kafka.Message{
		Value: []byte(message),
	}
	return p.writer.WriteMessages(context.Background(), msg)
}

// Close 关闭 Kafka 生产者
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
