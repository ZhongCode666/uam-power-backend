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
	"uam-power-backend/models/controller_models/aircraft_task_model"
	"uam-power-backend/models/controller_models/data_flow_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

type KafkaToMysql struct {
	KafkaEventConsumerService  *dbservice.KafkaConsumer
	KafkaStatusConsumerService *dbservice.KafkaConsumer
	MysqlStatusService         *dbservice.MySQLService
	MysqlEventService          *dbservice.MySQLService
	RedisService               *dbservice.RedisDict
	StopFlag                   bool
	StatusDone                 chan bool
	EventDone                  chan bool
	wg                         sync.WaitGroup
}

func NewKafkaToMysql(
	KafkaConfig *db_config_model.KafkaConfigModel,
	MySqlConfig *db_config_model.MySqlConfigModel,
	RedisConfig *db_config_model.RedisConfigModel,
) *KafkaToMysql {
	kafkaStatus := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftDataTopic, "KafkaToMysql")
	kafkaEvent := dbservice.NewKafkaConsumer(KafkaConfig.Addr, KafkaConfig.AircraftEventTopic, "KafkaToMysql")
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlConfig.Usr, MySqlConfig.Psw, MySqlConfig.Host, MySqlConfig.Port,
		MySqlConfig.FlightDB,
	)
	FlightMysqlService, FlightErr := dbservice.NewMySQLService(mysqlLink)
	if FlightErr != nil {
		return nil
	}
	mysqlLink = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlConfig.Usr, MySqlConfig.Psw, MySqlConfig.Host, MySqlConfig.Port,
		MySqlConfig.EventDB,
	)
	EventMysqlService, EventErr := dbservice.NewMySQLService(mysqlLink)
	if EventErr != nil {
		return nil
	}
	RedisInfo := dbservice.NewRedisDict(RedisConfig.Host, RedisConfig.Port, RedisConfig.TaskInfoDBno)
	utils.MsgSuccess("        [KafkaToMysql]Successfully init!")
	return &KafkaToMysql{
		KafkaEventConsumerService:  kafkaEvent,
		KafkaStatusConsumerService: kafkaStatus,
		MysqlStatusService:         FlightMysqlService,
		MysqlEventService:          EventMysqlService,
		RedisService:               RedisInfo,
		StatusDone:                 make(chan bool),
		EventDone:                  make(chan bool),
		StopFlag:                   false,
		wg:                         sync.WaitGroup{},
	}
}

func (ser *KafkaToMysql) KafkaStatusToMysql() {
	utils.MsgSuccess("        [KafkaToMysql]start KafkaStatusToMysql successfully!")
	for !ser.StopFlag {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		KafkaRe, err := ser.KafkaStatusConsumerService.ReceiveMessage(ctx)
		if err != nil {
			utils.MsgError("        [KafkaToMysql]receive msg error")
			continue
		}

		var reStruct data_flow_model.AircraftStatus
		err = json.Unmarshal([]byte(KafkaRe), &reStruct)
		if err != nil {
			utils.MsgError("        [KafkaToMysql]invalid json")
			continue
		}
		re, redisErr := ser.RedisService.Get(strconv.Itoa(reStruct.AircraftID))
		if redisErr != nil {
			utils.MsgError("        [KafkaToMysql]Can not hit redis!")
			continue
		}
		if re == nil {
			utils.MsgError("        [KafkaToMysql]Can not hit redis!")
			continue
		}
		jsonData, _ := json.Marshal(re)
		var mysqlData aircraft_task_model.MysqlAircraftTask
		err = json.Unmarshal(jsonData, &mysqlData)
		if err != nil {
			utils.MsgError("        [KafkaToMysql]Invalid Json!")
			continue
		}
		sql := fmt.Sprintf("INSERT INTO flightdb.%s (Longitude, Latitude, Altitude, Yaw, DataTime) VALUES (%f, %f, %f, %f, '%s');",
			mysqlData.TrackTable, reStruct.Longitude, reStruct.Latitude, reStruct.Altitude, reStruct.Yaw,
			reStruct.TimeString,
		)
		//utils.MsgInfo("SQL>" + sql)
		_, err = ser.MysqlStatusService.ExecuteCmd(sql)
		if err != nil {
			utils.MsgError("        [KafkaToMysql]Can not insert! err>" + err.Error())
			continue
		}
		utils.MsgSuccess("        [KafkaToMysql]successfully insert status!")
	}
	ser.StatusDone <- true
	ser.wg.Done()
}

func (ser *KafkaToMysql) KafkaEventToMysql() {
	utils.MsgSuccess("        [KafkaToMysql]start KafkaEventToMysql successfully!")
	for !ser.StopFlag {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		KafkaRe, err := ser.KafkaEventConsumerService.ReceiveMessage(ctx)
		if err != nil {
			utils.MsgError("        [KafkaToMysql]KafkaEventToMysql receive msg error err>" + err.Error())
			continue
		}
		var reStruct data_flow_model.AircraftEvent
		err = json.Unmarshal([]byte(KafkaRe), &reStruct)
		if err != nil {
			utils.MsgError("        [KafkaToMysql]KafkaEventToMysql invalid json")
			continue
		}
		re, redisErr := ser.RedisService.Get(strconv.Itoa(reStruct.AircraftID))
		if redisErr != nil {
			utils.MsgError("        [KafkaToMysql]KafkaEventToMysql failed to hit!")
			continue
		}
		jsonData, _ := json.Marshal(re)
		var mysqlData aircraft_task_model.MysqlAircraftTask
		err = json.Unmarshal(jsonData, &mysqlData)
		_, err = ser.MysqlEventService.ExecuteCmd(
			fmt.Sprintf("INSERT INTO %s(DataTime, Event) VALUES ('%s', '%s')",
				mysqlData.EventTable, reStruct.TimeString, reStruct.Event,
			))
		if err != nil {
			utils.MsgError("        [KafkaToMysql]KafkaEventToMysql failed to insert!")
			continue
		}
		utils.MsgSuccess("        [KafkaToMysql]KafkaEventToMysql successfully insert!")
	}
	ser.EventDone <- true
	ser.wg.Done()
}

func (ser *KafkaToMysql) Stop() {
	ser.StopFlag = true
	<-ser.StatusDone
	<-ser.EventDone
}

func (ser *KafkaToMysql) Start() {
	go ser.KafkaStatusToMysql()
	go ser.KafkaEventToMysql()
}

func main() {
	utils.MsgInfo("    [KafkaToMysql]Process ready to start")
	//print("    [KafkaToMysql]Process ready to start")
	args := os.Args[1]

	cfg, _ := utils.LoadDBConfig(args)

	// 获取从主进程传递过来的参数
	ser := NewKafkaToMysql(&cfg.KafkaCfg, &cfg.MySqlCfg, &cfg.RedisCfg)
	ser.wg.Add(2)
	ser.Start()
	ser.wg.Wait()
	utils.MsgSuccess("    [KafkaToMysql]Process successfully started!")
	ser.Stop()
	utils.MsgSuccess("    [KafkaToMysql]Process successfully stopped!")

}
