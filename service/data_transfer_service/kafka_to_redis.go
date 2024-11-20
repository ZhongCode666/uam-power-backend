package data_transfer_service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
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
	wg                         sync.WaitGroup
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
		wg:                         sync.WaitGroup{},
	}
}

func (ser *KafkaToRedis) KafkaStatusToRedis() {
	utils.MsgSuccess("        [KafkaToRedis]start KafkaStatusToRedis successfully!")
	for !ser.StopFlag {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
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
	ser.wg.Done()
}

func (ser *KafkaToRedis) KafkaEventToRedis() {
	utils.MsgSuccess("        [KafkaToRedis]start KafkaEventToRedis successfully!")
	for !ser.StopFlag {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
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
	ser.wg.Done()
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

func main() {
	utils.MsgInfo("    [KafkaToRedis]Process ready to start")
	args := os.Args[1]
	cfg, _ := utils.LoadDBConfig(args)

	// 获取从主进程传递过来的参数
	ser := NewKafkaToRedis(&cfg.KafkaCfg, &cfg.RedisCfg)
	ser.wg.Add(2)
	ser.Start()
	utils.MsgSuccess("    [KafkaToRedis]Process successfully started!")
	ser.wg.Wait()
	//var input string
	//_, err := fmt.Scan(&input)
	//if err != nil {
	//	return
	//}
	ser.Stop()
	utils.MsgSuccess("    [KafkaToRedis]Process successfully stopped!")
}
