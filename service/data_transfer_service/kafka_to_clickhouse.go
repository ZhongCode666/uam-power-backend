package data_transfer_service

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/aircraft_task_model"
	"uam-power-backend/models/controller_models/data_flow_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

// KafkaToClickhouse 负责将 Kafka 消息传输到 ClickHouse
type KafkaToClickhouse struct {
	// KafkaEventConsumerService 处理 Kafka 事件消息的消费者服务
	KafkaEventConsumerService *dbservice.KafkaConsumer
	// KafkaStatusConsumerService 处理 Kafka 状态消息的消费者服务
	KafkaStatusConsumerService *dbservice.KafkaConsumer
	// ClickHouseStatusService 处理 ClickHouse 状态数据的服务
	ClickHouseStatusService *dbservice.ClickHouse
	// ClickHouseEventService 处理 ClickHouse 事件数据的服务
	ClickHouseEventService *dbservice.ClickHouse
	// RedisService 处理 Redis 数据的服务
	RedisService *dbservice.RedisDict
	// StopFlag 用于控制停止任务的标志
	StopFlag bool
	// StatusDone 用于通知状态任务完成的通道
	StatusDone chan bool
	// EventDone 用于通知事件任务完成的通道
	EventDone chan bool
	// wg 用于等待任务完成的同步等待组
	wg sync.WaitGroup
}

// NewKafkaToClickhouse 创建一个新的 KafkaToClickhouse 实例
func NewKafkaToClickhouse(
	KafkaConfig *db_config_model.KafkaConfigModel, // Kafka 配置
	ClickHouseConfig *db_config_model.ClickHouseConfigModel, // ClickHouse 配置
	RedisConfig *db_config_model.RedisConfigModel, // Redis 配置
) *KafkaToClickhouse {
	// 创建 Kafka 消费者服务
	kafkaStatus := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftDataTopic, "KafkaToMysql")
	kafkaEvent := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftEventTopic, "KafkaToMysql")

	// 创建 ClickHouse 服务
	FlightClickHouseService, FlightErr := dbservice.NewClickHouse(
		ClickHouseConfig.Host, ClickHouseConfig.Port, ClickHouseConfig.Username, ClickHouseConfig.Password,
		ClickHouseConfig.BatchSize, ClickHouseConfig.FlushPeriod, ClickHouseConfig.FlightDatabase, true,
		ClickHouseConfig.FlightColumn)
	if FlightErr != nil {
		utils.MsgError("        [KafkaToMysql]Failed to init MySQLWithBufferService!")
		return nil
	}
	EventClickHouseService, EventErr := dbservice.NewClickHouse(
		ClickHouseConfig.Host, ClickHouseConfig.Port, ClickHouseConfig.Username, ClickHouseConfig.Password,
		ClickHouseConfig.EventBatchSize, ClickHouseConfig.EventFlushPeriod, ClickHouseConfig.EventDatabase, true,
		ClickHouseConfig.EventColumn)
	if EventErr != nil {
		utils.MsgError("        [KafkaToMysql]Failed to init MySQLWithBufferService!")
		return nil
	}

	// 创建 Redis 服务
	RedisInfo := dbservice.NewRedisDict(RedisConfig.Host, RedisConfig.Port, RedisConfig.TaskInfoDBno)

	// 打印初始化成功信息
	utils.MsgSuccess("        [KafkaToMysql]Successfully init!")

	// 返回 KafkaToClickhouse 实例
	return &KafkaToClickhouse{
		KafkaEventConsumerService:  kafkaEvent,              // Kafka 事件消费者服务
		KafkaStatusConsumerService: kafkaStatus,             // Kafka 状态消费者服务
		ClickHouseStatusService:    FlightClickHouseService, // ClickHouse 状态服务
		ClickHouseEventService:     EventClickHouseService,  // ClickHouse 事件服务
		RedisService:               RedisInfo,               // Redis 服务
		StatusDone:                 make(chan bool),         // 状态完成通道
		EventDone:                  make(chan bool),         // 事件完成通道
		StopFlag:                   false,                   // 停止标志
		wg:                         sync.WaitGroup{},        // 同步等待组
	}
}

// KafkaStatusToClickhouse 将 Kafka 状态消息传输到 ClickHouse
func (ser *KafkaToClickhouse) KafkaStatusToClickhouse() {
	utils.MsgSuccess("        [KafkaToMysql]start KafkaStatusToClickhouse successfully!")
	for !ser.StopFlag {
		// 设置上下文超时时间为 1 秒
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		// 接收 Kafka 消息
		KafkaRe, err := ser.KafkaStatusConsumerService.ReceiveMessage(ctx)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaStatusToClickhouse receive msg error")
			continue
		}

		// 反序列化 Kafka 消息
		var reStruct data_flow_model.AircraftStatus
		err = json.Unmarshal([]byte(KafkaRe), &reStruct)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaStatusToClickhouse invalid json")
			continue
		}
		// 从 Redis 获取数据
		re, redisErr := ser.RedisService.Get(strconv.Itoa(reStruct.AircraftID))
		if redisErr != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaStatusToClickhouse Can not hit redis!")
			continue
		}
		if re == nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaStatusToClickhouse Can not hit redis!")
			continue
		}
		// 将 Redis 数据序列化为 JSON
		jsonData, _ := json.Marshal(re)
		var mysqlData aircraft_task_model.MysqlAircraftTask
		// 反序列化 JSON 数据
		err = json.Unmarshal(jsonData, &mysqlData)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaStatusToClickhouse Invalid Json!")
			continue
		}
		// 将数据添加到 ClickHouse
		err = ser.ClickHouseStatusService.Add(
			mysqlData.TrackTable,
			[]string{"Longitude", "Latitude", "Altitude", "Yaw", "DataTime"},
			[]interface{}{reStruct.Longitude, reStruct.Latitude, reStruct.Altitude, reStruct.Yaw, reStruct.TimeString},
		)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaStatusToClickhouse Can not insert! err>" + err.Error())
			return
		}
		utils.MsgSuccess("        [KafkaToClickhouse]KafkaStatusToClickhouse successfully insert status!")
	}
	// 通知状态任务完成
	ser.StatusDone <- true
	// 等待组任务完成
	ser.wg.Done()
}

// KafkaEventToClickhouse 将 Kafka 事件消息传输到 ClickHouse
func (ser *KafkaToClickhouse) KafkaEventToClickhouse() {
	// 打印启动成功信息
	utils.MsgSuccess("        [KafkaToClickhouse]KafkaEventToClickhouse start KafkaEventToMysql successfully!")
	for !ser.StopFlag {
		// 设置上下文超时时间为 1 秒
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		// 接收 Kafka 消息
		KafkaRe, err := ser.KafkaEventConsumerService.ReceiveMessage(ctx)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaEventToClickhouse receive msg error err>" + err.Error())
			continue
		}
		// 反序列化 Kafka 消息
		var reStruct data_flow_model.AircraftEvent
		err = json.Unmarshal([]byte(KafkaRe), &reStruct)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaEventToClickhouse invalid json")
			continue
		}
		// 从 Redis 获取数据
		re, redisErr := ser.RedisService.Get(strconv.Itoa(reStruct.AircraftID))
		if redisErr != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaEventToClickhouse failed to hit!")
			continue
		}
		// 将 Redis 数据序列化为 JSON
		jsonData, _ := json.Marshal(re)
		var mysqlData aircraft_task_model.MysqlAircraftTask
		// 反序列化 JSON 数据
		err = json.Unmarshal(jsonData, &mysqlData)
		// 将数据添加到 ClickHouse
		err = ser.ClickHouseStatusService.Add(
			mysqlData.EventTable,
			[]string{"DataTime", "Event"},
			[]interface{}{reStruct.TimeString, reStruct.Event},
		)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaEventToClickhouse failed to insert!")
			continue
		}
		// 打印插入成功信息
		utils.MsgSuccess("        [KafkaToClickhouse]KafkaEventToClickhouse successfully insert!")
	}
	// 通知事件任务完成
	ser.EventDone <- true
	// 等待组任务完成
	ser.wg.Done()
}

// Stop 停止 KafkaToClickhouse 服务
func (ser *KafkaToClickhouse) Stop() {
	ser.StopFlag = true // 设置停止标志
	<-ser.StatusDone    // 等待状态任务完成
	<-ser.EventDone     // 等待事件任务完成
}

// Start 启动 KafkaToClickhouse 服务
func (ser *KafkaToClickhouse) Start() {
	// 启动 Kafka 状态消息传输到 ClickHouse 的协程
	go ser.KafkaStatusToClickhouse()
	// 启动 Kafka 事件消息传输到 ClickHouse 的协程
	go ser.KafkaEventToClickhouse()
}
