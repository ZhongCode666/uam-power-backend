package aircraft_id_controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/aircraft_id_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

// AircraftIdController 结构体定义
type AircraftIdController struct {
	IDMySql   *dbservice.MySQLService // MySQL 服务
	RedisInfo *dbservice.RedisDict    // Redis 信息服务
}

// NewAircraftIdController 创建并初始化 AircraftIdController 实例
func NewAircraftIdController(
	MySqlCfg *db_config_model.MySqlConfigModel, RedisCfg *db_config_model.RedisConfigModel,
) *AircraftIdController {
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
	// 初始化 Redis 信息服务
	RedisInfo := dbservice.NewRedisDict(RedisCfg.Host, RedisCfg.Port, RedisCfg.AircraftDBno)
	// 打印初始化成功信息
	utils.MsgInfo("        [NewAircraftIdController]Successfully init!")
	// 返回 AircraftIdController 实例
	return &AircraftIdController{IDMySql: MysqlService, RedisInfo: RedisInfo}
}

// GetAircraftInfo 获取飞机信息
// 处理传入的 JSON 请求，首先从 Redis 获取数据，如果未命中则从 MySQL 获取数据并存储到 Redis
func (a *AircraftIdController) GetAircraftInfo(c *fiber.Ctx) error {
	var RequestID aircraft_id_model.GetAircraftInfoID
	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&RequestID); err != nil {
		utils.MsgError("        [NewAircraftIdController]GetAircraftInfo Request Invalid JSON data!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 从 Redis 获取数据
	re, _ := a.RedisInfo.Get(strconv.Itoa(RequestID.AircraftID))
	if re != nil {
		// 打印命中 Redis 信息并返回数据
		utils.MsgSuccess("        [NewAircraftIdController]GetAircraftInfo Hit Redis auto Return!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Successfully GetAircraftInfo!", "data": re})
	}
	// 从 MySQL 获取数据
	mysqlRe, mysqlErr := a.IDMySql.QueryRow(fmt.Sprintf("Select * from systemdb.aircraft_identity_table where AircraftID = %d;", RequestID.AircraftID))
	if mysqlErr != nil {
		// 打印错误信息并返回错误响应
		utils.MsgError("        [NewAircraftIdController]GetAircraftInfo No such Aircraft!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A.!"})
	}
	// 将 MySQL 数据序列化为 JSON
	jsonData, _ := json.Marshal(mysqlRe)
	// 将数据存储到 Redis
	err := a.RedisInfo.Set(strconv.Itoa(RequestID.AircraftID), string(jsonData))
	if err != nil {
		// 打印错误信息并返回错误响应
		utils.MsgError("        [NewAircraftIdController]Set Redis Failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Redis failed!"})
	}
	// 打印成功信息并返回数据
	utils.MsgSuccess("        [NewAircraftIdController]Successfully GetAircraftInfo!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully GetAircraftInfo!", "data": mysqlRe})
}

// CreateUser 创建用户
// 处理传入的 JSON 请求，将数据存储到 MySQL 并缓存到 Redis
func (a *AircraftIdController) CreateUser(c *fiber.Ctx) error {
	curStr := utils.GetTimeStr()                   // 获取当前时间字符串
	randStr := curStr + "-" + utils.GetUniqueStr() // 生成唯一字符串
	var RequestInfo aircraft_id_model.SetAircraftInfo
	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&RequestInfo); err != nil {
		utils.MsgError("        [NewAircraftIdController]CreateUser invalid Requests Json!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 执行 MySQL 插入命令
	_, err := a.IDMySql.ExecuteCmd(
		fmt.Sprintf("INSERT INTO systemdb.aircraft_identity_table(Type, Company, Name, TimeStr) VALUES ('%s', '%s', '%s', '%s')",
			RequestInfo.Type, RequestInfo.Company, RequestInfo.Name, randStr,
		))
	if err != nil {
		utils.MsgError("        [NewAircraftIdController]CreateUser failed to Mysql!, err>" + err.Error())
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Send to Mysql Failed"})
	}
	// 从 MySQL 查询插入的数据
	mysqlRe, mysqlErr := a.IDMySql.QueryRow(
		fmt.Sprintf("Select * from systemdb.aircraft_identity_table where TimeStr = '%s';",
			randStr))
	if mysqlErr != nil {
		utils.MsgError("        [NewAircraftIdController]CreateUser data in MySql not found!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A.!"})
	}
	// 将 MySQL 数据序列化为 JSON
	jsonData, _ := json.Marshal(mysqlRe)
	var mysqlData aircraft_id_model.MysqlAircraftInfo
	// 反序列化 JSON 数据
	err = json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [NewAircraftIdController]CreateUser failed to redis!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "failed to send Redis!"})
	}
	// 将数据存储到 Redis
	err = a.RedisInfo.Set(strconv.Itoa(mysqlData.AircraftID), string(jsonData))
	if err != nil {
		utils.MsgError("        [NewAircraftIdController]CreateUser failed to redis!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "failed to send Redis!"})
	}
	utils.MsgSuccess("        [NewAircraftIdController]Successfully CreateUser!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully CreateUser!", "data": mysqlRe})
}
