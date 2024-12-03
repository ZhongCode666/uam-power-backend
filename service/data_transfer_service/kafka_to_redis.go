package data_transfer_service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/data_flow_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

// KafkaToRedis 结构体定义
type KafkaToRedis struct {
	KafkaEventConsumerService  *dbservice.KafkaConsumer // Kafka 事件消费者服务
	KafkaStatusConsumerService *dbservice.KafkaConsumer // Kafka 状态消费者服务
	RedisStatusService         *dbservice.RedisDict     // Redis 状态服务
	RedisEventService          *dbservice.RedisDict     // Redis 事件服务
	StopFlag                   bool                     // 停止标志
	StatusDone                 chan bool                // 状态完成通道
	EventDone                  chan bool                // 事件完成通道
	wg                         sync.WaitGroup           // 同步等待组
}

// NewKafkaToRedis 创建并初始化 KafkaToRedis 实例
func NewKafkaToRedis(
	KafkaConfig *db_config_model.KafkaConfigModel, RedisConfig *db_config_model.RedisConfigModel,
) *KafkaToRedis {
	// 初始化 Redis 状态服务
	redisStatus := dbservice.NewRedisDict(RedisConfig.Host, RedisConfig.Port, RedisConfig.StatusDBno)
	// 初始化 Redis 事件服务
	redisEvent := dbservice.NewRedisDict(RedisConfig.Host, RedisConfig.Port, RedisConfig.EventDBno)
	// 初始化 Kafka 状态消费者服务
	kafkaStatus := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftDataTopic, "KafkaToRedis")
	// 初始化 Kafka 事件消费者服务
	kafkaEvent := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftEventTopic, "KafkaToRedis")
	// 打印初始化成功信息
	utils.MsgSuccess("        [KafkaToRedis]init successfully!")
	// 返回 KafkaToRedis 实例
	return &KafkaToRedis{
		KafkaEventConsumerService:  kafkaEvent,       // Kafka 事件消费者服务
		KafkaStatusConsumerService: kafkaStatus,      // Kafka 状态消费者服务
		RedisStatusService:         redisStatus,      // Redis 状态服务
		RedisEventService:          redisEvent,       // Redis 事件服务
		StatusDone:                 make(chan bool),  // 状态完成通道
		EventDone:                  make(chan bool),  // 事件完成通道
		StopFlag:                   false,            // 停止标志
		wg:                         sync.WaitGroup{}, // 同步等待组
	}
}

// KafkaStatusToRedis 将 Kafka 状态消息传输到 Redis
func (ser *KafkaToRedis) KafkaStatusToRedis() {
	utils.MsgSuccess("        [KafkaToRedis]start KafkaStatusToRedis successfully!") // 打印启动成功信息
	for !ser.StopFlag {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second) // 设置上下文超时时间为 1 秒
		KafkaRe, err := ser.KafkaStatusConsumerService.ReceiveMessage(ctx) // 接收 Kafka 消息
		if err != nil {
			utils.MsgError("        [KafkaToRedis]receive msg error") // 打印接收消息错误信息
			continue
		}
		utils.MsgSuccess(fmt.Sprintf("        [KafkaToRedis]KafkaStatusToRedis receive msg -> %s", KafkaRe)) // 打印接收消息成功信息

		var reStruct data_flow_model.AircraftStatus
		err = json.Unmarshal([]byte(KafkaRe), &reStruct) // 反序列化 Kafka 消息
		if err != nil {
			utils.MsgError("        [KafkaToRedis]invalid json") // 打印无效 JSON 错误信息
			continue
		}
		err = ser.RedisStatusService.SetWithDuration(strconv.Itoa(reStruct.AircraftID), KafkaRe, 3600) // 将数据存储到 Redis，设置过期时间为 3600 秒
		if err != nil {
			utils.MsgError("        [KafkaToRedis]invalid json") // 打印存储数据错误信息
			continue
		}
		utils.MsgSuccess("        [KafkaToRedis]KafkaStatusToRedis successfully!") // 打印成功存储数据信息
	}
	ser.StatusDone <- true // 通知状态任务完成
	ser.wg.Done()          // 等待组任务完成
}

// KafkaEventToRedis 将 Kafka 事件消息传输到 Redis
func (ser *KafkaToRedis) KafkaEventToRedis() {
	utils.MsgSuccess("        [KafkaToRedis]start KafkaEventToRedis successfully!") // 打印启动成功信息
	for !ser.StopFlag {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second) // 设置上下文超时时间为 1 秒
		KafkaRe, err := ser.KafkaEventConsumerService.ReceiveMessage(ctx)  // 接收 Kafka 消息
		if err != nil {
			utils.MsgError("        [KafkaToRedis]KafkaEventToRedis receive msg error") // 打印接收消息错误信息
			continue
		}
		utils.MsgSuccess(fmt.Sprintf("        [KafkaToRedis]KafkaEventToRedis receive msg -> %s", KafkaRe)) // 打印接收消息成功信息
		var reStruct data_flow_model.AircraftEvent
		err = json.Unmarshal([]byte(KafkaRe), &reStruct) // 反序列化 Kafka 消息
		if err != nil {
			utils.MsgError("        [KafkaToRedis]KafkaEventToRedis invalid json") // 打印无效 JSON 错误信息
			continue
		}
		err = ser.RedisEventService.SetWithDuration(strconv.Itoa(reStruct.AircraftID), KafkaRe, 3600) // 将数据存储到 Redis，设置过期时间为 3600 秒
		if err != nil {
			utils.MsgError("        [KafkaToRedis]KafkaEventToRedis invalid json") // 打印存储数据错误信息
			continue
		}
		utils.MsgSuccess("        [KafkaToRedis]KafkaEventToRedis successfully!") // 打印成功存储数据信息
	}
	ser.EventDone <- true // 通知事件任务完成
	ser.wg.Done()         // 等待组任务完成
}

// Stop 停止 KafkaToRedis 服务
func (ser *KafkaToRedis) Stop() {
	ser.StopFlag = true // 设置停止标志
	<-ser.StatusDone    // 等待状态任务完成
	<-ser.EventDone     // 等待事件任务完成
}

// Start 启动 KafkaToRedis 服务
func (ser *KafkaToRedis) Start() {
	// 启动 Kafka 状态消息传输到 Redis 的协程
	go ser.KafkaStatusToRedis()
	// 启动 Kafka 事件消息传输到 Redis 的协程
	go ser.KafkaEventToRedis()
}
