package lane_controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/lane_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

// LaneController 航线控制器
type LaneController struct {
	// LaneMongoDB MongoDB 客户端
	LaneMongoDB *dbservice.MongoDBClient
	// LaneMysql MySQL 服务
	LaneMysql *dbservice.MySQLService
}

// NewLaneController 创建一个新的 LaneController 实例
// 参数:
// - MongoCfg: MongoDB 配置
// - MySqlCfg: MySQL 配置
// 返回值:
// - *LaneController: 返回一个 LaneController 实例
func NewLaneController(
	MongoCfg *db_config_model.MongoConfigModel, MySqlCfg *db_config_model.MySqlConfigModel,
) *LaneController {
	// 构建 MySQL 连接字符串
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlCfg.Usr, MySqlCfg.Psw, MySqlCfg.Host, MySqlCfg.Port,
		MySqlCfg.DB,
	)
	// 创建 MySQL 服务
	MysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		utils.MsgError("        [LaneController]ini mysql service error:" + err.Error())
		return nil
	}
	utils.MsgSuccess("        [LaneController]Successfully LaneMysql!")
	// 构建 MongoDB 连接字符串
	mongoLink := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		MongoCfg.Usr, MongoCfg.Psw, MongoCfg.Host, MongoCfg.Port,
	)
	// 创建 MongoDB 客户端
	mongoService, MongoErr := dbservice.NewMongoDBClient(mongoLink, MongoCfg.AreaDB)
	if MongoErr != nil {
		utils.MsgError("        [LaneController]MongoDB client error:" + MongoErr.Error())
	}
	utils.MsgSuccess("        [LaneController]Successfully LaneMongo!")
	// 返回 LaneController 实例
	return &LaneController{
		LaneMysql:   MysqlService,
		LaneMongoDB: mongoService,
	}
}

// CreateLane 创建航线
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (laneModel *LaneController) CreateLane(c *fiber.Ctx) error {
	// 获取当前时间字符串
	curStr := utils.GetTimeStr()
	// 生成唯一字符串
	randStr := curStr + "-" + utils.GetUniqueStr()
	var LaneInfo lane_model.LaneCreateModel
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&LaneInfo); err != nil {
		utils.MsgError("        [LaneController]CreateLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 插入新航线数据到 MySQL
	_, err := laneModel.LaneMysql.ExecuteCmd(
		fmt.Sprintf("INSERT INTO systemdb.lane_table(IsHide, Name, TimeStr) VALUES (%t, '%s', '%s');",
			LaneInfo.IsHide, LaneInfo.Name, randStr,
		))
	if err != nil {
		utils.MsgError("        [LaneController]CreateLane Create Lane Failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert failed"})
	}
	// 查询插入的航线数据
	mysqlRe, mysqlErr := laneModel.LaneMysql.QueryRow(
		fmt.Sprintf("Select * from systemdb.lane_table where TimeStr = '%s';",
			randStr))
	if mysqlErr != nil {
		utils.MsgError("        [LaneController]CreateLane Query sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "N.A.!"})
	}
	// 将查询结果转换为 JSON
	jsonData, _ := json.Marshal(mysqlRe)
	var mysqlData lane_model.LaneMysqlData
	// 解析 JSON 数据
	err = json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [LaneController]CreateLane Unmarshal Json data failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Invalid Data"})
	}
	// 构造 MongoDB 文档
	document := bson.M{
		"LaneID": mysqlData.LaneID, "PointData": LaneInfo.PointData,
		"RasterData": LaneInfo.RasterData, "IsHide": mysqlData.IsHide == 1,
	}
	// 插入新航线数据到 MongoDB
	_, err = laneModel.LaneMongoDB.InsertOne("lane_data_collection", document)
	if err != nil {
		utils.MsgError("        [LaneController]CreateLane Insert Mongo failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert Mongo failed"})
	}
	utils.MsgSuccess("        [LaneController]CreateLane Successfully CreateLane!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully CreateLane!", "data": mysqlRe})
}

// GetLane 获取航线信息
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (laneModel *LaneController) GetLane(c *fiber.Ctx) error {
	var LaneInfo lane_model.FindLaneModel
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&LaneInfo); err != nil {
		utils.MsgError("        [LaneController]GetLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 构造 MongoDB 查询过滤器
	filter := bson.M{"LaneID": LaneInfo.LaneID, "IsHide": false}
	// 指定要排除的字段
	drops := bson.M{"_id": 0}
	// 从 MongoDB 中查询航线数据
	re, err := laneModel.LaneMongoDB.FindOneWithDropRow("lane_data_collection", filter, drops)
	if err != nil {
		utils.MsgError("        [LaneController]GetLane Find Lane error!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A."})
	}
	if re == nil {
		utils.MsgError("        [LaneController]GetLane no such lane!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A."})
	}
	utils.MsgSuccess("        [LaneController]GetLane Successfully GetLane!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully return Lane!", "data": re})
}

// GetAllLane 获取所有航线信息
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (laneModel *LaneController) GetAllLane(c *fiber.Ctx) error {
	var LaneInfo lane_model.FindLaneModel
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&LaneInfo); err != nil {
		utils.MsgError("        [LaneController]GetAllLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 构造 MongoDB 查询过滤器
	filter := bson.M{"IsHide": false}
	// 指定要排除的字段
	drops := bson.M{"_id": 0}
	// 从 MongoDB 中查询所有航线数据
	re, err := laneModel.LaneMongoDB.FindAllWithDrops("lane_data_collection", filter, drops)
	if err != nil {
		utils.MsgError("        [LaneController]GetAllLane Find Lane error!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A."})
	}
	if re == nil {
		utils.MsgError("        [LaneController]GetAllLane no such lane!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A."})
	}
	utils.MsgSuccess("        [LaneController]GetAllLane Successfully GetLane!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully return Lane!", "data": re})
}

// LaneList 获取航线列表
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (laneModel *LaneController) LaneList(c *fiber.Ctx) error {
	// 从 MySQL 查询所有航线数据
	re, err := laneModel.LaneMysql.QueryRows("Select * from lane_table;")
	if err != nil {
		utils.MsgError("        [LaneController]LaneList failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Get data from mysql error"})
	}
	utils.MsgSuccess("        [LaneController]LaneList Successfully GetList!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully return Lane!", "data": re})
}

// HideLane 隐藏航线
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (laneModel *LaneController) HideLane(c *fiber.Ctx) error {
	var showHideLaneModel lane_model.ShowHideLaneModel
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&showHideLaneModel); err != nil {
		utils.MsgError("        [LaneController]HideLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 更新 MySQL 中的航线状态为隐藏
	_, err := laneModel.LaneMysql.ExecuteCmd(
		fmt.Sprintf("UPDATE lane_table SET IsHide = 1 WHERE LaneID IN %s;",
			"("+strings.Trim(strings.Join(strings.Fields(fmt.Sprint(showHideLaneModel.LaneIDs)), ", "), "[]")+")",
		))
	if err != nil {
		utils.MsgError("        [LaneController]HideLane Change mysql lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mysql failed!"})
	}
	// 构造 MongoDB 查询过滤器
	filter := bson.M{"LaneID": bson.M{"$in": showHideLaneModel.LaneIDs}}
	// 构造 MongoDB 更新操作
	update := bson.M{"$set": bson.M{"IsHide": true}}
	// 更新 MongoDB 中的航线状态为隐藏
	_, MErr := laneModel.LaneMongoDB.Update("lane_data_collection", filter, update)
	if MErr != nil {
		utils.MsgError("        [LaneController]HideLane Change mongo lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mongo failed!"})
	}
	utils.MsgSuccess("        [LaneController]HideLane Successfully HideLane!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully hide Lanes!"})
}

// ShowLane 显示航线
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (laneModel *LaneController) ShowLane(c *fiber.Ctx) error {
	var showHideLaneModel lane_model.ShowHideLaneModel
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&showHideLaneModel); err != nil {
		utils.MsgError("        [LaneController]ShowLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 更新 MySQL 中的航线状态为显示
	_, err := laneModel.LaneMysql.ExecuteCmd(
		fmt.Sprintf("UPDATE lane_table SET IsHide = 0 WHERE LaneID IN %s;",
			"("+strings.Trim(strings.Join(strings.Fields(fmt.Sprint(showHideLaneModel.LaneIDs)), ", "), "[]")+")",
		))
	if err != nil {
		utils.MsgError("        [LaneController]ShowLane Change mysql lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mysql failed!"})
	}
	// 构造 MongoDB 查询过滤器
	filter := bson.M{"LaneID": bson.M{"$in": showHideLaneModel.LaneIDs}}
	// 构造 MongoDB 更新操作
	update := bson.M{"$set": bson.M{"IsHide": false}}
	// 更新 MongoDB 中的航线状态为显示
	_, MErr := laneModel.LaneMongoDB.Update("lane_data_collection", filter, update)
	if MErr != nil {
		utils.MsgError("        [LaneController]ShowLane Change mongo lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mongo failed!"})
	}
	utils.MsgSuccess("        [LaneController]ShowLane Successfully ShowLane!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully show Lanes!"})
}

// DeleteLane 删除航线
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (laneModel *LaneController) DeleteLane(c *fiber.Ctx) error {
	var showHideLaneModel lane_model.ShowHideLaneModel
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&showHideLaneModel); err != nil {
		utils.MsgError("        [LaneController]DeleteLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 从 MySQL 中删除航线数据
	_, err := laneModel.LaneMysql.ExecuteCmd(
		fmt.Sprintf("DELETE FROM lane_table WHERE LaneID IN %s;",
			"("+strings.Trim(strings.Join(strings.Fields(fmt.Sprint(showHideLaneModel.LaneIDs)), ", "), "[]")+")",
		))
	if err != nil {
		utils.MsgError("        [LaneController]DeleteLane Change mysql lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mysql failed!"})
	}
	// 构造 MongoDB 查询过滤器
	filter := bson.M{"LaneID": bson.M{"$in": showHideLaneModel.LaneIDs}}
	// 从 MongoDB 中删除航线数据
	_, MErr := laneModel.LaneMongoDB.Delete("lane_data_collection", filter)
	if MErr != nil {
		utils.MsgError("        [LaneController]DeleteLane Change mongo lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mongo failed!"})
	}
	utils.MsgSuccess("        [LaneController]DeleteLane Successfully DeleteLane!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully delete Lanes!"})
}
