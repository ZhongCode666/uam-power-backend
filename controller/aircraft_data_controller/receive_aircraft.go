package aircraft_data_controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/data_flow_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

type RequestAircraft struct {
	StatusRedisService *dbservice.RedisDict
	EventRedisService  *dbservice.RedisDict
	TaskRedisService   *dbservice.RedisDict
	FlightMysql        *dbservice.MySQLService
	EventMysql         *dbservice.MySQLService
}

func NewReceiveAircraft(
	redisConfig *db_config_model.RedisConfigModel,
	mysqlConfig *db_config_model.MySqlConfigModel,
) *RequestAircraft {
	utils.MsgSuccess("        [ReceiveAircraft]init successfully!")
	redisStatusService := dbservice.NewRedisDict(redisConfig.Host, redisConfig.Port, redisConfig.StatusDBno)
	redisEventService := dbservice.NewRedisDict(redisConfig.Host, redisConfig.Port, redisConfig.EventDBno)
	redisTaskService := dbservice.NewRedisDict(redisConfig.Host, redisConfig.Port, redisConfig.TaskInfoDBno)
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.Usr, mysqlConfig.Psw, mysqlConfig.Host, mysqlConfig.Port,
		mysqlConfig.FlightDB,
	)
	FlightMysqlService, FlightErr := dbservice.NewMySQLService(mysqlLink)
	if FlightErr != nil {
		return nil
	}
	mysqlLink = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.Usr, mysqlConfig.Psw, mysqlConfig.Host, mysqlConfig.Port,
		mysqlConfig.EventDB,
	)
	EventMysqlService, EventErr := dbservice.NewMySQLService(mysqlLink)
	if EventErr != nil {
		return nil
	}
	return &RequestAircraft{
		StatusRedisService: redisStatusService, EventRedisService: redisEventService,
		TaskRedisService: redisTaskService, FlightMysql: FlightMysqlService,
		EventMysql: EventMysqlService,
	}
}

func (receiver *RequestAircraft) RequestAircraftStatus(c *fiber.Ctx) error {
	var aircraftReq data_flow_model.RecAircraftStatusRequest

	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftReq); err != nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftStatus Invalid JSON data! >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}

	rec, err := receiver.StatusRedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil || rec == nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftStatus Invalid JSON data!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A.!"})
	}
	utils.MsgError("        [ReceiveAircraft]RequestAircraftStatus Invalid JSON data!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully requestData!", "data": rec})
}

func (receiver *RequestAircraft) RequestAircraftEvent(c *fiber.Ctx) error {
	var aircraftReq data_flow_model.RecAircraftStatusRequest

	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftReq); err != nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftEvent Invalid JSON data! >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}

	rec, err := receiver.EventRedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil || rec == nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftEvent Invalid JSON data!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A.!"})
	}
	utils.MsgSuccess("        [ReceiveAircraft]RequestAircraftEvent Successfully requestData!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully requestData!", "data": rec})
}

func (receiver *RequestAircraft) ReceiveActiveData(c *fiber.Ctx) error {
	keys, err := receiver.TaskRedisService.Keys()
	if err != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveActiveData Get activate aircraft failed!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Get all active aircraft failed!"})
	}
	if len(keys) == 0 {
		utils.MsgError("        [ReceiveAircraft]ReceiveActiveData No active aircraft!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "No active aircraft!"})
	}
	activateResult, GetValErr := receiver.StatusRedisService.GetVals(keys)
	if GetValErr != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveActiveData Get active aircraft Data failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Get active aircraft Data failed!"})
	}
	utils.MsgSuccess("        [ReceiveAircraft]ReceiveActiveData Successfully request activate Data!")
	returnData := bson.M{"ActivateIDs": keys, "PosData": activateResult}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully request activate Data!", "data": returnData})
}

func (receiver *RequestAircraft) ReceiveActiveEvent(c *fiber.Ctx) error {
	keys, err := receiver.TaskRedisService.Keys()
	if err != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveActiveEvent Get activate aircraft failed!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Get all active aircraft failed!"})
	}
	if len(keys) == 0 {
		utils.MsgError("        [ReceiveAircraft]ReceiveActiveEvent No active aircraft!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "No active aircraft!"})
	}
	activateResult, GetValErr := receiver.EventRedisService.GetVals(keys)
	if GetValErr != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveActiveEvent No active aircraft!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Get active aircraft Event failed!"})
	}
	utils.MsgSuccess("        [ReceiveAircraft]ReceiveActiveEvent Successfully request activate Event!")
	returnData := bson.M{"ActivateIDs": keys, "EventData": activateResult}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully request activate Event!", "data": returnData})
}

func (receiver *RequestAircraft) ReceiveHistoryData(c *fiber.Ctx) error {
	var aircraftReq data_flow_model.RecAircraftStatusRequest

	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftReq); err != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveHistoryData Invalid JSON data! >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	re, err := receiver.TaskRedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveHistoryData no such Task!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Not Found"})
	}
	jsonData, _ := json.Marshal(re)
	var mysqlData data_flow_model.TaskAircraftRedis
	err = json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveHistoryData Unmarshal fail!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Unmarshal fail!"})
	}
	rows, MysqlErr := receiver.FlightMysql.QueryRows(fmt.Sprintf("SELECT * FROM %s;", mysqlData.TrackTable))
	if MysqlErr != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveHistoryData find data from mysql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Find data from mysql failed"})
	}
	utils.MsgSuccess("        [ReceiveAircraft]ReceiveHistoryData Successfully request activate data history!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully request activate Data history!", "data": rows})
}

func (receiver *RequestAircraft) ReceiveActiveIDs(ctx *fiber.Ctx) error {
	keys, err := receiver.TaskRedisService.Keys()
	if err != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveActiveIDs Get activate aircraft failed!")
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Get all active aircraft failed!"})
	}
	if len(keys) == 0 {
		utils.MsgError("        [ReceiveAircraft]ReceiveActiveIDs No active aircraft!")
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "No active aircraft!"})
	}
	utils.MsgSuccess("        [ReceiveAircraft]ReceiveActiveIDs Successfully request activate Data!")
	returnData := bson.M{"ActivateIDs": keys}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully request activate Data!", "data": returnData})
}

func (receiver *RequestAircraft) ReceiveHistoryEvent(c *fiber.Ctx) error {
	var aircraftReq data_flow_model.RecAircraftStatusRequest

	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftReq); err != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveHistoryEvent Invalid JSON data! >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	re, err := receiver.TaskRedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveHistoryEvent no such Task!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Not Found"})
	}
	jsonData, _ := json.Marshal(re)
	var mysqlData data_flow_model.TaskAircraftRedis
	err = json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveHistoryEvent Unmarshal fail!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Unmarshal fail!"})
	}
	rows, MysqlErr := receiver.FlightMysql.QueryRows(fmt.Sprintf("SELECT * FROM %s;", mysqlData.EventTable))
	if MysqlErr != nil {
		utils.MsgError("        [ReceiveAircraft]ReceiveHistoryEvent find Event from mysql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Find Event from mysql failed"})
	}
	utils.MsgSuccess("        [ReceiveAircraft]ReceiveHistoryEvent Successfully request activate Event history!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully request activate Event history!", "data": rows})
}
