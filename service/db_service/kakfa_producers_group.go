package dbservice

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// KafkaManager 管理多个生产者的 Kafka 类
type KafkaManager struct {
	producers []*KafkaProducer
	mu        sync.Mutex
	counter   uint64
}

// NewKafkaManager 创建 KafkaManager
func NewKafkaManager(addr string, topic string, producerCount int) *KafkaManager {
	manager := &KafkaManager{
		producers: make([]*KafkaProducer, 0, producerCount),
	}

	// 初始化多个生产者
	for i := 0; i < producerCount; i++ {
		producer := NewKafkaProducer(addr, topic)
		manager.producers = append(manager.producers, producer)
	}
	return manager
}

// SendMsg 发送消息到指定生产者
func (km *KafkaManager) SendMsg(producerID int, message string) error {
	km.mu.Lock()
	defer km.mu.Unlock()

	if producerID < 0 || producerID >= len(km.producers) {
		return fmt.Errorf("invalid producer ID: %d", producerID)
	}

	producer := km.producers[producerID]
	return producer.SendMessage(message)
}

// SendMsgRoundRobin 使用无锁轮询策略发送消息
func (km *KafkaManager) SendMsgRoundRobin(message string) error {
	// 获取当前生产者的索引
	index := atomic.AddUint64(&km.counter, 1) % uint64(len(km.producers))

	// 选择生产者并发送消息
	producer := km.producers[index]
	return producer.SendMessage(message)
}

// Close 关闭所有生产者
func (km *KafkaManager) Close() error {
	for _, producer := range km.producers {
		if err := producer.Close(); err != nil {
			return err
		}
	}
	return nil
}
