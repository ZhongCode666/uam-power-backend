package dbservice

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// KafkaManager 管理多个生产者的 Kafka 类
type KafkaManager struct {
	// producers 是 Kafka 生产者的切片
	producers []*KafkaProducer
	// mu 是用于保护共享资源的互斥锁
	mu sync.Mutex
	// counter 是用于轮询策略的计数器
	counter uint64
}

// NewKafkaManager 创建 KafkaManager
// addr 是 Kafka 代理地址
// topic 是 Kafka 主题
// producerCount 是要创建的生产者数量
func NewKafkaManager(addr string, topic string, producerCount int) *KafkaManager {
	manager := &KafkaManager{
		producers: make([]*KafkaProducer, 0, producerCount), // 初始化生产者切片
	}

	// 初始化多个生产者
	for i := 0; i < producerCount; i++ {
		producer := NewKafkaProducer(addr, topic)               // 创建新的 Kafka 生产者
		manager.producers = append(manager.producers, producer) // 将生产者添加到切片中
	}
	return manager // 返回 KafkaManager 实例
}

// SendMsg 发送消息到指定生产者
// producerID 是生产者的 ID
// message 是要发送的消息字符串
// 返回可能的错误
func (km *KafkaManager) SendMsg(producerID int, message string) error {
	km.mu.Lock()         // 加锁以保护共享资源
	defer km.mu.Unlock() // 函数结束时解锁

	if producerID < 0 || producerID >= len(km.producers) {
		return fmt.Errorf("invalid producer ID: %d", producerID) // 如果生产者 ID 无效，返回错误
	}

	producer := km.producers[producerID] // 获取指定的生产者
	return producer.SendMessage(message) // 发送消息并返回可能的错误
}

// SendMsgRoundRobin 使用无锁轮询策略发送消息
// message 是要发送的消息字符串
// 返回可能的错误
func (km *KafkaManager) SendMsgRoundRobin(message string) error {
	// 获取当前生产者的索引
	index := atomic.AddUint64(&km.counter, 1) % uint64(len(km.producers))

	// 选择生产者并发送消息
	producer := km.producers[index]
	return producer.SendMessage(message)
}

// Close 关闭所有生产者
// 关闭 Kafka 生产者写入器并释放相关资源
func (km *KafkaManager) Close() error {
	for _, producer := range km.producers {
		if err := producer.Close(); err != nil {
			return err // 如果关闭失败，返回错误
		}
	}
	return nil // 成功关闭所有生产者，返回 nil
}
