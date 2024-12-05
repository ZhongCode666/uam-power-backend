package data_transfer_service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/aircraft_task_model"
	"uam-power-backend/models/controller_models/data_flow_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

// KafkaToMysql 结构体定义
type KafkaToMysql struct {
	KafkaEventConsumerService  *dbservice.KafkaConsumer          // Kafka 事件消费者服务
	KafkaStatusConsumerService *dbservice.KafkaConsumer          // Kafka 状态消费者服务
	MysqlStatusService         *dbservice.MySQLWithBufferService // MySQL 状态服务
	MysqlEventService          *dbservice.MySQLWithBufferService // MySQL 事件服务
	RedisService               *dbservice.RedisDict              // Redis 服务
	StopFlag                   bool                              // 停止标志
	StatusDone                 chan bool                         // 状态完成通道
	EventDone                  chan bool                         // 事件完成通道
	wg                         sync.WaitGroup                    // 同步等待组
}

// NewKafkaToMysql 创建并初始化 KafkaToMysql 实例
func NewKafkaToMysql(
	KafkaConfig *db_config_model.KafkaConfigModel, // Kafka 配置
	MySqlConfig *db_config_model.MySqlConfigModel, // MySQL 配置
	RedisConfig *db_config_model.RedisConfigModel, // Redis 配置
) *KafkaToMysql {
	// 创建 Kafka 消费者服务
	kafkaStatus := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftDataTopic, "KafkaToMysql")
	kafkaEvent := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftEventTopic, "KafkaToMysql")

	// 创建 MySQL 连接字符串
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlConfig.Usr, MySqlConfig.Psw, MySqlConfig.Host, MySqlConfig.Port,
		MySqlConfig.FlightDB,
	)
	// 初始化 MySQL 状态服务
	FlightMysqlService, FlightErr := dbservice.NewMySQLWithBufferService(mysqlLink, MySqlConfig.FlightInterval, MySqlConfig.FlightColumn)
	if FlightErr != nil {
		utils.MsgError("        [KafkaToMysql]Failed to init MySQLWithBufferService!")
		return nil
	}

	// 创建 MySQL 事件服务连接字符串
	mysqlLink = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlConfig.Usr, MySqlConfig.Psw, MySqlConfig.Host, MySqlConfig.Port,
		MySqlConfig.EventDB,
	)
	// 初始化 MySQL 事件服务
	EventMysqlService, EventErr := dbservice.NewMySQLWithBufferService(mysqlLink, MySqlConfig.EventInterval, MySqlConfig.EventColumn)
	if EventErr != nil {
		utils.MsgError("        [KafkaToMysql]Failed to init MySQLWithBufferService!")
		return nil
	}

	// 创建 Redis 服务
	RedisInfo := dbservice.NewRedisDict(RedisConfig.Host, RedisConfig.Port, RedisConfig.TaskInfoDBno)

	// 打印初始化成功信息
	utils.MsgSuccess("        [KafkaToMysql]Successfully init!")

	// 返回 KafkaToMysql 实例
	return &KafkaToMysql{
		KafkaEventConsumerService:  kafkaEvent,         // Kafka 事件消费者服务
		KafkaStatusConsumerService: kafkaStatus,        // Kafka 状态消费者服务
		MysqlStatusService:         FlightMysqlService, // MySQL 状态服务
		MysqlEventService:          EventMysqlService,  // MySQL 事件服务
		RedisService:               RedisInfo,          // Redis 服务
		StatusDone:                 make(chan bool),    // 状态完成通道
		EventDone:                  make(chan bool),    // 事件完成通道
		StopFlag:                   false,              // 停止标志
		wg:                         sync.WaitGroup{},   // 同步等待组
	}
}

// KafkaStatusToMysql 将 Kafka 状态消息传输到 MySQL
func (ser *KafkaToMysql) KafkaStatusToMysql() {
	utils.MsgSuccess("        [KafkaToMysql]start KafkaStatusToMysql successfully!")
	for !ser.StopFlag {
		// 设置上下文超时时间为 1 秒
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		// 接收 Kafka 消息
		KafkaRe, err := ser.KafkaStatusConsumerService.ReceiveMessage(ctx)
		if err != nil {
			utils.MsgError("        [KafkaToMysql]receive msg error")
			continue
		}

		// 反序列化 Kafka 消息
		var reStruct data_flow_model.AircraftStatus
		err = json.Unmarshal([]byte(KafkaRe), &reStruct)
		if err != nil {
			utils.MsgError("        [KafkaToMysql]invalid json")
			continue
		}
		// 从 Redis 获取数据
		re, redisErr := ser.RedisService.Get(strconv.Itoa(reStruct.AircraftID))
		if redisErr != nil {
			utils.MsgError("        [KafkaToMysql]Can not hit redis!")
			continue
		}
		if re == nil {
			utils.MsgError("        [KafkaToMysql]Can not hit redis!")
			continue
		}
		// 将 Redis 数据序列化为 JSON
		jsonData, _ := json.Marshal(re)
		var mysqlData aircraft_task_model.MysqlAircraftTask
		// 反序列化 JSON 数据
		err = json.Unmarshal(jsonData, &mysqlData)
		if err != nil {
			utils.MsgError("        [KafkaToMysql]Invalid Json!")
			continue
		}
		// 将数据添加到 MySQL
		arr := []interface{}{reStruct.Longitude, reStruct.Latitude, reStruct.Altitude, reStruct.Yaw, reStruct.TimeString}
		ser.MysqlStatusService.Add(mysqlData.TrackTable, arr)
		utils.MsgSuccess("        [KafkaToMysql]successfully insert status!")
	}
	// 通知状态任务完成
	ser.StatusDone <- true
	// 等待组任务完成
	ser.wg.Done()
}

// KafkaEventToMysql 将 Kafka 事件消息传输到 MySQL
func (ser *KafkaToMysql) KafkaEventToMysql() {
	// 打印启动成功信息
	utils.MsgSuccess("        [KafkaToMysql]start KafkaEventToMysql successfully!")
	for !ser.StopFlag {
		// 设置上下文超时时间为 1 秒
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		// 接收 Kafka 消息
		KafkaRe, err := ser.KafkaEventConsumerService.ReceiveMessage(ctx)
		if err != nil {
			utils.MsgError("        [KafkaToMysql]KafkaEventToMysql receive msg error err>" + err.Error())
			continue
		}
		// 反序列化 Kafka 消息
		var reStruct data_flow_model.AircraftEvent
		err = json.Unmarshal([]byte(KafkaRe), &reStruct)
		if err != nil {
			utils.MsgError("        [KafkaToMysql]KafkaEventToMysql invalid json")
			continue
		}
		// 从 Redis 获取数据
		re, redisErr := ser.RedisService.Get(strconv.Itoa(reStruct.AircraftID))
		if redisErr != nil {
			utils.MsgError("        [KafkaToMysql]KafkaEventToMysql failed to hit!")
			continue
		}
		// 将 Redis 数据序列化为 JSON
		jsonData, _ := json.Marshal(re)
		var mysqlData aircraft_task_model.MysqlAircraftTask
		// 反序列化 JSON 数据
		err = json.Unmarshal(jsonData, &mysqlData)
		// 将数据添加到 MySQL
		arr := []interface{}{reStruct.TimeString, reStruct.Event}
		ser.MysqlEventService.Add(mysqlData.EventTable, arr)
		// 打印插入成功信息
		utils.MsgSuccess("        [KafkaToMysql]KafkaEventToMysql successfully insert!")
	}
	// 通知事件任务完成
	ser.EventDone <- true
	// 等待组任务完成
	ser.wg.Done()
}

// Stop 停止 KafkaToMysql 服务
func (ser *KafkaToMysql) Stop() {
	ser.StopFlag = true // 设置停止标志
	<-ser.StatusDone    // 等待状态任务完成
	<-ser.EventDone     // 等待事件任务完成
}

// Start 启动 KafkaToMysql 服务
func (ser *KafkaToMysql) Start() {
	// 启动 Kafka 状态消息传输到 MySQL 的协程
	go ser.KafkaStatusToMysql()
	// 启动 Kafka 事件消息传输到 MySQL 的协程
	go ser.KafkaEventToMysql()
}
