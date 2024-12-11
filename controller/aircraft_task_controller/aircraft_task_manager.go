package aircraft_task_controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/aircraft_task_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

// AircraftTaskModel 结构体定义
type AircraftTaskModel struct {
	MysqlService       *dbservice.MySQLService // MySQL 服务
	RedisService       *dbservice.RedisDict    // Redis 服务
	FlightMysqlService *dbservice.MySQLService // 航班 MySQL 服务
	EventMysqlService  *dbservice.MySQLService // 事件 MySQL 服务
}

// NewAircraftTaskModel 创建并初始化 AircraftTaskModel 实例
func NewAircraftTaskModel(
	RedisCfg *db_config_model.RedisConfigModel, // Redis 配置
	MySqlCfg *db_config_model.MySqlConfigModel, // MySQL 配置
) *AircraftTaskModel {
	// 构建 MySQL 连接字符串
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlCfg.Usr, MySqlCfg.Psw, MySqlCfg.Host, MySqlCfg.Port,
		MySqlCfg.DB,
	)
	// 初始化 MySQL 服务
	MysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		return nil
	}
	utils.MsgSuccess("        [AircraftTaskModel]Successfully SystemMysql!")

	// 构建航班 MySQL 连接字符串
	mysqlLink = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlCfg.Usr, MySqlCfg.Psw, MySqlCfg.Host, MySqlCfg.Port,
		MySqlCfg.FlightDB,
	)
	// 初始化航班 MySQL 服务
	FlightMysqlService, FlightErr := dbservice.NewMySQLService(mysqlLink)
	if FlightErr != nil {
		return nil
	}
	utils.MsgSuccess("        [AircraftTaskModel]Successfully FlightMysql!")

	// 构建事件 MySQL 连接字符串
	mysqlLink = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlCfg.Usr, MySqlCfg.Psw, MySqlCfg.Host, MySqlCfg.Port,
		MySqlCfg.EventDB,
	)
	// 初始化事件 MySQL 服务
	EventMysqlService, EventErr := dbservice.NewMySQLService(mysqlLink)
	if EventErr != nil {
		return nil
	}
	utils.MsgSuccess("        [AircraftTaskModel]Successfully EventMysql!")

	// 初始化 Redis 服务
	RedisInfo := dbservice.NewRedisDict(RedisCfg.Host, RedisCfg.Port, RedisCfg.TaskInfoDBno)
	utils.MsgSuccess("        [AircraftTaskModel]Successfully Redis!")
	utils.MsgSuccess("        [AircraftTaskModel]Successfully init!")

	// 返回初始化后的 AircraftTaskModel 实例
	return &AircraftTaskModel{
		MysqlService: MysqlService, RedisService: RedisInfo,
		FlightMysqlService: FlightMysqlService, EventMysqlService: EventMysqlService,
	}
}

// CreateTask 创建任务
func (taskModel *AircraftTaskModel) CreateTask(c *fiber.Ctx) error {
	// 获取当前时间字符串
	curStr := utils.GetTimeStr()
	// 生成唯一字符串
	randStr := curStr + "-" + utils.GetUniqueStr()
	// 定义任务信息结构体
	var TaskInfo aircraft_task_model.CreateTaskAircraftInfo
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&TaskInfo); err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	LaneMysqlRe, LaneMysqlErr := taskModel.MysqlService.QueryRow(
		fmt.Sprintf("Select * from systemdb.lane_table where LaneID = %d;",
			TaskInfo.LaneID))
	if LaneMysqlErr != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask Query lane sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "N.A.!"})
	}
	if LaneMysqlRe == nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask No such Lane!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "No data found!"})
	}
	// 构建航班表名
	FlightTable := fmt.Sprintf("%sFlight_AirID%d_Lane%d", curStr, TaskInfo.AircraftID, TaskInfo.LaneID)
	// 构建事件表名
	EventTable := fmt.Sprintf("%sEvent_AirID%d_Lane%d", curStr, TaskInfo.AircraftID, TaskInfo.LaneID)
	// 创建航班表
	_, err := taskModel.FlightMysqlService.ExecuteCmd(
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (Longitude DOUBLE(15, 12), Latitude DOUBLE(15, 12), Altitude DOUBLE(15, 12), Yaw DOUBLE(15, 12), DataTime DATETIME(6),  UploadTime DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6));",
			FlightTable,
		))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask Create Table Failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Create Status Table Failed!"})
	}
	// 创建事件表
	_, err = taskModel.EventMysqlService.ExecuteCmd(
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (DataTime DATETIME(6),  CreateTime DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6), Event char(20) not NULL);",
			EventTable,
		))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask Create Table Failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Create Event Table Failed!"})
	}
	// 插入任务信息到系统表
	_, err = taskModel.MysqlService.ExecuteCmd(
		fmt.Sprintf("INSERT INTO systemdb.flight_task_table(AircraftID, LaneID, TrackTable, EventTable, TimeStr) VALUES (%d, %d, '%s', '%s', '%s');",
			TaskInfo.AircraftID, TaskInfo.LaneID, FlightTable, EventTable, randStr,
		))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask Create Task Failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert failed"})
	}
	// 查询插入的任务信息
	mysqlRe, mysqlErr := taskModel.MysqlService.QueryRow(
		fmt.Sprintf("Select * from systemdb.flight_task_table where TimeStr = '%s';",
			randStr))
	if mysqlErr != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask Query sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "N.A.!"})
	}
	// 将查询结果转换为 JSON
	jsonData, _ := json.Marshal(mysqlRe)
	// 定义 MySQL 任务信息结构体
	var mysqlData aircraft_task_model.MysqlAircraftTask
	// 解析 JSON 数据
	err = json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask failed to redis!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Hit redis Failed"})
	}
	// 将任务信息存储到 Redis
	err = taskModel.RedisService.Set(strconv.Itoa(mysqlData.AircraftID), string(jsonData))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CreateTask failed to redis!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Hit redis Failed"})
	}
	utils.MsgSuccess("        [AircraftTaskModel]Successfully create Task!")
	// 返回成功响应
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully CreateTask!", "data": mysqlRe})
}

// EndTask 结束任务
func (taskModel *AircraftTaskModel) EndTask(c *fiber.Ctx) error {
	var aircraftReq aircraft_task_model.ByAircraftID
	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftReq); err != nil {
		utils.MsgError("        [AircraftTaskModel]EndTask Request Invalid JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Request Invalid JSON data"})
	}
	// 从 Redis 获取任务信息
	re, err := taskModel.RedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]EndTask No such Task!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "No such Task!"})
	}
	if re == nil {
		utils.MsgError("        [AircraftTaskModel]EndTask No such Task!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "No such Task!"})
	}
	//var mysqlData aircraft_task_model.MysqlAircraftTask
	//// 将 Redis 数据转换为 JSON
	//jsonData, _ := json.Marshal(re)
	//// 将 JSON 数据解析为 MySQL 任务信息结构体
	//err = json.Unmarshal(jsonData, &mysqlData)
	//if err != nil {
	//	utils.MsgError("        [AircraftTaskModel]EndTask Get Task ID failed!")
	//	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "No such Task!"})
	//}
	// 更新 MySQL 中的任务结束时间
	_, MysqlErr := taskModel.MysqlService.ExecuteCmd(
		fmt.Sprintf("UPDATE systemdb.flight_task_table SET EndTime = '%s' WHERE AreaID = %d;",
			utils.GetMySqlTimeStr(), aircraftReq.AircraftID,
		),
	)
	if MysqlErr != nil {
		utils.MsgError("        [AircraftTaskModel]EndTask Set Task ID failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Failed to end Task!"})
	}
	// 从 Redis 中删除任务信息
	DeleteErr := taskModel.RedisService.Delete(strconv.Itoa(aircraftReq.AircraftID))
	if DeleteErr != nil {
		utils.MsgError("        [AircraftTaskModel]EndTask Set to Redis failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Failed to end set Redis!"})
	}
	utils.MsgSuccess("        [AircraftTaskModel]Successfully EndTask!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully EndTask!"})
}

// CheckTaskInfo 检查任务信息
func (taskModel *AircraftTaskModel) CheckTaskInfo(c *fiber.Ctx) error {
	var aircraftReq aircraft_task_model.ByAircraftIDAndTaskID

	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftReq); err != nil {
		utils.MsgError("        [AircraftTaskModel]CheckTaskInfo Invalid request JSON data!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 从 Redis 获取任务信息
	re, err := taskModel.RedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil {
		utils.MsgError("        [AircraftTaskModel]CheckTaskInfo no such Task!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Not Found"})
	}
	if re == nil {
		// 从 MySQL 获取任务信息
		row, err := taskModel.MysqlService.QueryRow(
			fmt.Sprintf("SELECT * FROM systemdb.flight_task_table WHERE TaskID = %d;",
				aircraftReq.TaskID))
		if err != nil {
			utils.MsgError("        [AircraftTaskModel]CheckTaskInfo no such Task!")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Not Found"})
		}
		// 将查询结果转换为 JSON
		jsonData, _ := json.Marshal(row)
		var mysqlData aircraft_task_model.MysqlAircraftTask
		// 将 JSON 数据解析为 MySQL 任务信息结构体
		err = json.Unmarshal(jsonData, &mysqlData)
		if err != nil {
			utils.MsgError("        [AircraftTaskModel]CheckTaskInfo no such Task!")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Not Found"})
		}
		// 将任务信息存储到 Redis
		err = taskModel.RedisService.Set(strconv.Itoa(mysqlData.AircraftID), string(jsonData))
		if err != nil {
			utils.MsgError("        [AircraftTaskModel]CreateTask failed to redis!")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Hit redis Failed"})
		}
		// 再次从 Redis 获取任务信息
		re, err = taskModel.RedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
		if err != nil {
			utils.MsgError("        [AircraftTaskModel]CreateTask failed to redis!")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Hit redis Failed"})
		}
	}
	utils.MsgSuccess("        [AircraftTaskModel]CheckTaskInfo TaskInfo!")
	// 返回成功响应
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "CheckTaskInfo TaskInfo!", "data": re})
}
