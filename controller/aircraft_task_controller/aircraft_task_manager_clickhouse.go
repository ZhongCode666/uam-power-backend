package aircraft_task_controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/aircraft_task_model"
	dbservice "uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

type AircraftTaskModelClickHouse struct {
	MysqlService            *dbservice.MySQLService
	RedisService            *dbservice.RedisDict
	FlightClickHouseService *dbservice.ClickHouse
	EventClickHouseService  *dbservice.ClickHouse
}

func NewAircraftTaskModelClickHouse(
	RedisCfg *db_config_model.RedisConfigModel,
	MySqlCfg *db_config_model.MySqlConfigModel,
	ClickHouseConfig *db_config_model.ClickHouseConfigModel,
) *AircraftTaskModelClickHouse {
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlCfg.Usr, MySqlCfg.Psw, MySqlCfg.Host, MySqlCfg.Port,
		MySqlCfg.DB,
	)
	MysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		return nil
	}
	utils.MsgSuccess("        [AircraftTaskModel]Successfully SystemMysql!")
	FlightClickHouseService, FlightErr := dbservice.NewClickHouse(
		ClickHouseConfig.Host, ClickHouseConfig.Port, ClickHouseConfig.Username, ClickHouseConfig.Password,
		ClickHouseConfig.BatchSize, ClickHouseConfig.FlushPeriod, ClickHouseConfig.FlightDatabase, false,
		ClickHouseConfig.FlightColumn)
	if FlightErr != nil {
		return nil
	}
	utils.MsgSuccess("        [AircraftTaskModel]Successfully FlightMysql!")
	EventClickHouseService, EventErr := dbservice.NewClickHouse(
		ClickHouseConfig.Host, ClickHouseConfig.Port, ClickHouseConfig.Username, ClickHouseConfig.Password,
		ClickHouseConfig.BatchSize, ClickHouseConfig.FlushPeriod, ClickHouseConfig.EventDatabase, false,
		ClickHouseConfig.EventColumn)
	if EventErr != nil {
		return nil
	}
	utils.MsgSuccess("        [AircraftTaskModel]Successfully EventMysql!")
	RedisInfo := dbservice.NewRedisDict(RedisCfg.Host, RedisCfg.Port, RedisCfg.TaskInfoDBno)
	utils.MsgSuccess("        [AircraftTaskModel]Successfully Redis!")
	utils.MsgSuccess("        [AircraftTaskModel]Successfully init!")
	return &AircraftTaskModelClickHouse{
		MysqlService: MysqlService, RedisService: RedisInfo,
		FlightClickHouseService: FlightClickHouseService, EventClickHouseService: EventClickHouseService,
	}
}

func (taskModel *AircraftTaskModelClickHouse) CreateTask(c *fiber.Ctx) error {
	curStr := utils.GetTimeStr()
	var TaskInfo aircraft_task_model.CreateTaskAircraftInfo
	if err := c.BodyParser(&TaskInfo); err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	FlightTable := fmt.Sprintf("%sFlight_AirID%d_Lane%d", curStr, TaskInfo.AircraftID, TaskInfo.LaneID)
	EventTable := fmt.Sprintf("%sEvent_AirID%d_Lane%d", curStr, TaskInfo.AircraftID, TaskInfo.LaneID)
	err := taskModel.FlightClickHouseService.ExecuteCmd(
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (Longitude Float64, Latitude Float64, Altitude Float64, Yaw Float64, DataTime DateTime64(6),  UploadTime DateTime64(6) DEFAULT now64(6)) ENGINE = MergeTree() ORDER BY (DataTime DESC);",
			FlightTable,
		))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask Create Table Failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Create Status Table Failed!"})
	}
	err = taskModel.EventClickHouseService.ExecuteCmd(
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (DataTime DateTime64(6),  CreateTime DateTime64(6) DEFAULT now64(6), Event String)  ENGINE = MergeTree() ORDER BY (DataTime DESC);",
			EventTable,
		))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask Create Table Failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Create Event Table Failed!"})
	}
	_, err = taskModel.MysqlService.ExecuteCmd(
		fmt.Sprintf("INSERT INTO systemdb.flight_task_table(AircraftID, LaneID, TrackTable, EventTable, TimeStr) VALUES (%d, %d, '%s', '%s', '%s');",
			TaskInfo.AircraftID, TaskInfo.LaneID, FlightTable, EventTable, curStr,
		))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask Create Task Failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert failed"})
	}
	mysqlRe, mysqlErr := taskModel.MysqlService.QueryRow(
		fmt.Sprintf("Select * from systemdb.flight_task_table where TimeStr = '%s';",
			curStr))
	if mysqlErr != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask Query sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "N.A.!"})
	}
	jsonData, _ := json.Marshal(mysqlRe)
	var mysqlData aircraft_task_model.MysqlAircraftTask
	err = json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask failed to redis!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Hit redis Failed"})
	}
	err = taskModel.RedisService.Set(strconv.Itoa(mysqlData.AircraftID), string(jsonData))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask failed to redis!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Hit redis Failed"})
	}
	utils.MsgSuccess("        [AircraftTaskModel]Successfully create Task!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully CreateTask!", "data": mysqlRe})
}

func (taskModel *AircraftTaskModelClickHouse) EndTask(c *fiber.Ctx) error {
	var aircraftReq aircraft_task_model.ByAircraftID
	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftReq); err != nil {
		utils.MsgError("        [AircraftTaskModel]EndTask Request Invalid JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Request Invalid JSON data"})
	}
	re, err := taskModel.RedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]EndTask No such Task!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "No such Task!"})
	}
	if re == nil {
		utils.MsgError("        [AircraftTaskModel]EndTask No such Task!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "No such Task!"})
	}
	var mysqlData aircraft_task_model.MysqlAircraftTask
	jsonData, _ := json.Marshal(re)
	err = json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]EndTask Get Task ID failed!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "No such Task!"})
	}
	_, MysqlErr := taskModel.MysqlService.ExecuteCmd(
		fmt.Sprintf("UPDATE systemdb.flight_task_table SET EndTime = '%s' WHERE TaskID = %d;",
			utils.GetMySqlTimeStr(), mysqlData.TaskID,
		),
	)
	if MysqlErr != nil {
		utils.MsgError("        [AircraftTaskModel]EndTask Set Task ID failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Failed to end Task!"})
	}
	DeleteErr := taskModel.RedisService.Delete(strconv.Itoa(aircraftReq.AircraftID))
	if DeleteErr != nil {
		utils.MsgError("        [AircraftTaskModel]EndTask Set to Redis failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Failed to end set Redis!"})
	}
	utils.MsgSuccess("        [AircraftTaskModel]Successfully EndTask!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully EndTask!"})
}

func (taskModel *AircraftTaskModelClickHouse) CheckTaskInfo(c *fiber.Ctx) error {
	var aircraftReq aircraft_task_model.ByAircraftIDAndTaskID

	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftReq); err != nil {
		utils.MsgError("        [AircraftTaskModel]CheckTaskInfo Invalid request JSON data!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	re, err := taskModel.RedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CheckTaskInfo no such Task!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Not Found"})
	}
	if re == nil {
		row, err := taskModel.MysqlService.QueryRow(
			fmt.Sprintf("SELECT * FROM systemdb.flight_task_table WHERE TaskID = %d;",
				aircraftReq.TaskID))
		if err != nil {
			utils.MsgError("        [AircraftTaskModel]CheckTaskInfo no such Task!")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Not Found"})
		}
		jsonData, _ := json.Marshal(row)
		var mysqlData aircraft_task_model.MysqlAircraftTask
		err = json.Unmarshal(jsonData, &mysqlData)
		if err != nil {
			utils.MsgError("        [AircraftTaskModel]CheckTaskInfo no such Task!")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Not Found"})
		}
		err = taskModel.RedisService.Set(strconv.Itoa(mysqlData.AircraftID), string(jsonData))
		if err != nil {
			utils.MsgError("        [AircraftTaskModel]CreateTask failed to redis!")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Hit redis Failed"})
		}
		re, err = taskModel.RedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
		if err != nil {
			utils.MsgError("        [AircraftTaskModel]CreateTask failed to redis!")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Hit redis Failed"})
		}
	}
	utils.MsgSuccess("        [AircraftTaskModel]CheckTaskInfo TaskInfo!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "CheckTaskInfo TaskInfo!", "data": re})
}
