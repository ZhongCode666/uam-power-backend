package data_transfer_service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/data_flow_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

type KafkaToRedis struct {
	KafkaEventConsumerService  *dbservice.KafkaConsumer
	KafkaStatusConsumerService *dbservice.KafkaConsumer
	RedisStatusService         *dbservice.RedisDict
	RedisEventService          *dbservice.RedisDict
	StopFlag                   bool
	StatusDone                 chan bool
	EventDone                  chan bool
}

func NewKafkaToRedis(
	KafkaConfig *db_config_model.KafkaConfigModel, RedisConfig *db_config_model.RedisConfigModel,
) *KafkaToRedis {
	redisStatus := dbservice.NewRedisDict(RedisConfig.Host, RedisConfig.Port, RedisConfig.StatusDBno)
	redisEvent := dbservice.NewRedisDict(RedisConfig.Host, RedisConfig.Port, RedisConfig.EventDBno)
	kafkaStatus := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftDataTopic, "KafkaToRedis")
	kafkaEvent := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftEventTopic, "KafkaToRedis")
	utils.MsgSuccess("        [KafkaToRedis]init successfully!")
	return &KafkaToRedis{
		KafkaEventConsumerService:  kafkaEvent,
		KafkaStatusConsumerService: kafkaStatus,
		RedisStatusService:         redisStatus,
		RedisEventService:          redisEvent,
		StatusDone:                 make(chan bool),
		EventDone:                  make(chan bool),
		StopFlag:                   false,
	}
}

func (ser *KafkaToRedis) KafkaStatusToRedis() {
	utils.MsgSuccess("        [KafkaToRedis]start KafkaStatusToRedis successfully!")
	for !ser.StopFlag {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		KafkaRe, err := ser.KafkaStatusConsumerService.ReceiveMessage(ctx)
		if err != nil {
			utils.MsgError("        [KafkaToRedis]receive msg error")
			continue
		}
		utils.MsgSuccess(fmt.Sprintf("        [KafkaToRedis]KafkaStatusToRedis receive msg -> %s", KafkaRe))

		var reStruct data_flow_model.AircraftStatus
		err = json.Unmarshal([]byte(KafkaRe), &reStruct)
		if err != nil {
			utils.MsgError("        [KafkaToRedis]invalid json")
			continue
		}
		err = ser.RedisStatusService.Set(strconv.Itoa(reStruct.AircraftID), KafkaRe)
		if err != nil {
			utils.MsgError("        [KafkaToRedis]invalid json")
			continue
		}
		utils.MsgSuccess("        [KafkaToRedis]KafkaStatusToRedis successfully!")
	}
	ser.StatusDone <- true
}

func (ser *KafkaToRedis) KafkaEventToRedis() {
	for !ser.StopFlag {
		utils.MsgSuccess("        [KafkaToRedis]start KafkaEventToRedis successfully!")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		KafkaRe, err := ser.KafkaEventConsumerService.ReceiveMessage(ctx)
		if err != nil {
			utils.MsgError("        [KafkaToRedis]KafkaEventToRedis receive msg error")
			continue
		}
		utils.MsgSuccess(fmt.Sprintf("        [KafkaToRedis]KafkaEventToRedis receive msg -> %s", KafkaRe))
		var reStruct data_flow_model.AircraftEvent
		err = json.Unmarshal([]byte(KafkaRe), &reStruct)
		if err != nil {
			utils.MsgError("        [KafkaToRedis]KafkaEventToRedis invalid json")
			continue
		}
		err = ser.RedisEventService.Set(strconv.Itoa(reStruct.AircraftID), KafkaRe)
		if err != nil {
			utils.MsgError("        [KafkaToRedis]KafkaEventToRedis invalid json")
			continue
		}
		utils.MsgSuccess("        [KafkaToRedis]KafkaEventToRedis successfully!")
	}
	ser.EventDone <- true
}

func (ser *KafkaToRedis) Stop() {
	ser.StopFlag = true
	<-ser.StatusDone
	<-ser.EventDone
}

func (ser *KafkaToRedis) Start() {
	go ser.KafkaStatusToRedis()
	go ser.KafkaEventToRedis()
}
