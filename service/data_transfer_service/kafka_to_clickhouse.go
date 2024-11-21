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

type KafkaToClickhouse struct {
	KafkaEventConsumerService  *dbservice.KafkaConsumer
	KafkaStatusConsumerService *dbservice.KafkaConsumer
	ClickHouseStatusService    *dbservice.ClickHouse
	ClickHouseEventService     *dbservice.ClickHouse
	RedisService               *dbservice.RedisDict
	StopFlag                   bool
	StatusDone                 chan bool
	EventDone                  chan bool
	wg                         sync.WaitGroup
}

func NewKafkaToClickhouse(
	KafkaConfig *db_config_model.KafkaConfigModel,
	ClickHouseConfig *db_config_model.ClickHouseConfigModel,
	RedisConfig *db_config_model.RedisConfigModel,
) *KafkaToClickhouse {
	kafkaStatus := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftDataTopic, "KafkaToMysql")
	kafkaEvent := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftEventTopic, "KafkaToMysql")
	FlightClickHouseService, FlightErr := dbservice.NewClickHouse(
		ClickHouseConfig.Host, ClickHouseConfig.Port, ClickHouseConfig.Username, ClickHouseConfig.Password,
		ClickHouseConfig.BatchSize, ClickHouseConfig.FlushPeriod, ClickHouseConfig.FlightDatabase, true,
		ClickHouseConfig.FlightColumn)
	if FlightErr != nil {
		return nil
	}
	EventClickHouseService, EventErr := dbservice.NewClickHouse(
		ClickHouseConfig.Host, ClickHouseConfig.Port, ClickHouseConfig.Username, ClickHouseConfig.Password,
		ClickHouseConfig.EventBatchSize, ClickHouseConfig.EventFlushPeriod, ClickHouseConfig.EventDatabase, true,
		ClickHouseConfig.EventColumn)
	if EventErr != nil {
		return nil
	}
	RedisInfo := dbservice.NewRedisDict(RedisConfig.Host, RedisConfig.Port, RedisConfig.TaskInfoDBno)
	utils.MsgSuccess("        [KafkaToMysql]Successfully init!")
	return &KafkaToClickhouse{
		KafkaEventConsumerService:  kafkaEvent,
		KafkaStatusConsumerService: kafkaStatus,
		ClickHouseStatusService:    FlightClickHouseService,
		ClickHouseEventService:     EventClickHouseService,
		RedisService:               RedisInfo,
		StatusDone:                 make(chan bool),
		EventDone:                  make(chan bool),
		StopFlag:                   false,
		wg:                         sync.WaitGroup{},
	}
}

func (ser *KafkaToClickhouse) KafkaStatusToClickhouse() {
	utils.MsgSuccess("        [KafkaToMysql]start KafkaStatusToClickhouse successfully!")
	for !ser.StopFlag {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		KafkaRe, err := ser.KafkaStatusConsumerService.ReceiveMessage(ctx)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaStatusToClickhouse receive msg error")
			continue
		}

		var reStruct data_flow_model.AircraftStatus
		err = json.Unmarshal([]byte(KafkaRe), &reStruct)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaStatusToClickhouse invalid json")
			continue
		}
		re, redisErr := ser.RedisService.Get(strconv.Itoa(reStruct.AircraftID))
		if redisErr != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaStatusToClickhouse Can not hit redis!")
			continue
		}
		if re == nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaStatusToClickhouse Can not hit redis!")
			continue
		}
		jsonData, _ := json.Marshal(re)
		var mysqlData aircraft_task_model.MysqlAircraftTask
		err = json.Unmarshal(jsonData, &mysqlData)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaStatusToClickhouse Invalid Json!")
			continue
		}
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
	ser.StatusDone <- true
	ser.wg.Done()
}

func (ser *KafkaToClickhouse) KafkaEventToClickhouse() {
	utils.MsgSuccess("        [KafkaToClickhouse]KafkaEventToClickhouse start KafkaEventToMysql successfully!")
	for !ser.StopFlag {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		KafkaRe, err := ser.KafkaEventConsumerService.ReceiveMessage(ctx)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaEventToClickhouse receive msg error err>" + err.Error())
			continue
		}
		var reStruct data_flow_model.AircraftEvent
		err = json.Unmarshal([]byte(KafkaRe), &reStruct)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaEventToClickhouse invalid json")
			continue
		}
		re, redisErr := ser.RedisService.Get(strconv.Itoa(reStruct.AircraftID))
		if redisErr != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaEventToClickhouse failed to hit!")
			continue
		}
		jsonData, _ := json.Marshal(re)
		var mysqlData aircraft_task_model.MysqlAircraftTask
		err = json.Unmarshal(jsonData, &mysqlData)
		err = ser.ClickHouseStatusService.Add(
			mysqlData.EventTable,
			[]string{"DataTime", "Event"},
			[]interface{}{reStruct.TimeString, reStruct.Event},
		)
		if err != nil {
			utils.MsgError("        [KafkaToClickhouse]KafkaEventToClickhouse failed to insert!")
			continue
		}
		utils.MsgSuccess("        [KafkaToClickhouse]KafkaEventToClickhouse successfully insert!")
	}
	ser.EventDone <- true
	ser.wg.Done()
}

func (ser *KafkaToClickhouse) Stop() {
	ser.StopFlag = true
	<-ser.StatusDone
	<-ser.EventDone
}

func (ser *KafkaToClickhouse) Start() {
	go ser.KafkaStatusToClickhouse()
	go ser.KafkaEventToClickhouse()
}
