package area_data_controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/area_model"
	dbservice "uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

// AreaController 结构体定义
type AreaController struct {
	AreaMongoDB *dbservice.MongoDBClient // MongoDB 客户端
	AreaMysql   *dbservice.MySQLService  // MySQL 服务
}

// NewAreaController 创建一个新的 AreaController 实例
// 参数:
// - MongoCfg: MongoDB 配置模型
// - MySqlCfg: MySQL 配置模型
// 返回值:
// - *AreaController: 返回一个新的 AreaController 实例
func NewAreaController(
	MongoCfg *db_config_model.MongoConfigModel, MySqlCfg *db_config_model.MySqlConfigModel,
) *AreaController {
	// 创建 MySQL 连接字符串
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlCfg.Usr, MySqlCfg.Psw, MySqlCfg.Host, MySqlCfg.Port,
		MySqlCfg.DB,
	)
	// 初始化 MySQL 服务
	MysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		utils.MsgError("        [AreaController]初始化 MySQL 服务错误:" + err.Error())
		return nil
	}
	utils.MsgSuccess("        [AreaController]成功初始化 AreaMysql!")

	// 创建 MongoDB 连接字符串
	mongoLink := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		MongoCfg.Usr, MongoCfg.Psw, MongoCfg.Host, MongoCfg.Port,
	)
	// 初始化 MongoDB 客户端
	mongoService, MongoErr := dbservice.NewMongoDBClient(mongoLink, MongoCfg.AreaDB)
	if MongoErr != nil {
		utils.MsgError("        [AreaController]MongoDB 客户端错误:" + MongoErr.Error())
	}
	utils.MsgSuccess("        [AreaController]成功初始化 AreaMongo!")

	// 返回 AreaController 实例
	return &AreaController{
		AreaMysql:   MysqlService,
		AreaMongoDB: mongoService,
	}
}

// GenerateAreaID 生成区域ID
// 参数:
// - Name: 区域名称
// 返回值:
// - int: 生成的区域ID
// - error: 错误信息
func (ac *AreaController) GenerateAreaID(Name string) (int, error) {
	// 获取当前时间字符串
	curStr := utils.GetTimeStr()
	// 生成唯一字符串
	randStr := curStr + "-" + utils.GetUniqueStr()
	// 执行插入区域数据的SQL命令
	_, err := ac.AreaMysql.ExecuteCmd(
		fmt.Sprintf("INSERT INTO systemdb.area_table(Name, TimeStr) VALUES ('%s', '%s');",
			Name, randStr,
		))
	if err != nil {
		utils.MsgError("        [AreaController]GenerateAreaID CreateID Area Failed!")
		return -1, err
	}
	// 查询插入的数据
	mysqlRe, mysqlErr := ac.AreaMysql.QueryRow(
		fmt.Sprintf("Select * from systemdb.area_table where TimeStr = '%s';",
			randStr))
	if mysqlErr != nil {
		utils.MsgError("        [AreaController]GenerateAreaID Query sql failed!")
		return -1, mysqlErr
	}
	// 将查询结果序列化为JSON
	jsonData, _ := json.Marshal(mysqlRe)
	var mysqlData area_model.AreaMysqlModel
	// 反序列化JSON数据
	err = json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [AreaController]GenerateAreaID Unmarshal Json data failed")
		return -1, err
	}
	// 返回生成的区域ID
	return mysqlData.AreaID, nil
}

// UpdateAreaData 更新区域数据
// 参数:
// - AreaID: 区域ID
// - RangeData: 范围数据
// - RasterSize: 栅格大小
// - RasterIndex: 栅格索引
// - RasterData: 栅格数据
// 返回值:
// - error: 错误信息
func (ac *AreaController) UpdateAreaData(
	AreaID int, RangeData []float64, RasterSize []float64, RasterIndex [][][]int,
	RasterData map[int]area_model.SingleRasterData,
) error {
	// 构造 MongoDB 文档
	document := bson.M{
		"AreaID": AreaID, "RangeData": RangeData,
		"RasterSize": RasterSize, "RasterIndex": RasterIndex,
	}
	// 插入文档到 MongoDB
	_, MErr := ac.AreaMongoDB.InsertOne("area_collection", document)
	if MErr != nil {
		utils.MsgError("        [AreaController]GenerateAreaID Insert Mongo failed")
		return MErr
	}
	var values []string
	// 构造每一行的值
	for id, record := range RasterData {
		status := record.Status
		longitude := record.Longitude
		latitude := record.Latitude
		altitude := record.Altitude

		values = append(
			values,
			fmt.Sprintf(
				"(%d, %d, %.12f, %.12f, %.12f, '%s')",
				id, AreaID, longitude, latitude, altitude, status,
			))
	}

	// 拼接最终 SQL
	sql := fmt.Sprintf(
		"INSERT INTO systemdb.raster_table(RasterID, AreaID, Longitude, Latitude, Altitude, Status) VALUES %s;",
		strings.Join(values, ", "))
	// 执行 SQL 命令
	_, MysqlErr := ac.AreaMysql.ExecuteCmd(sql)
	if MysqlErr != nil {
		utils.MsgError("        [AreaController]GenerateAreaID Insert Raster to Mysql failed err>" + MysqlErr.Error())
		return MysqlErr
	}
	return nil
}

// CreateAreaID 处理创建区域ID的请求
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (ac *AreaController) CreateAreaID(c *fiber.Ctx) error {
	var createAreaData area_model.GenerateAreaID
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&createAreaData); err != nil {
		utils.MsgError("        [AreaController]CreateAreaID Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 调用 GenerateAreaID 方法生成区域 ID
	AreaID, GenErr := ac.GenerateAreaID(createAreaData.Name)
	if GenErr != nil {
		utils.MsgError("        [AreaController]CreateAreaID Generate AreaID failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Generate AreaID failed"})
	}
	utils.MsgSuccess("        [AreaController]CreateAreaID Successfully CreateAreaID!")
	// 返回生成的区域 ID
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully CreateAreaID!", "data": bson.M{"AreaID": AreaID}})
}

// UploadArea 处理上传区域数据的请求
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (ac *AreaController) UploadArea(c *fiber.Ctx) error {
	var createAreaData area_model.UploadArea
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&createAreaData); err != nil {
		utils.MsgError("        [AreaController]UploadArea Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 更新区域数据
	err := ac.UpdateAreaData(
		createAreaData.AreaID, createAreaData.RangeData,
		createAreaData.RasterSize, createAreaData.RasterIndex,
		createAreaData.RasterData,
	)
	if err != nil {
		utils.MsgError("        [AreaController]UploadArea Insert data failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert data failed"})
	}
	utils.MsgSuccess("        [AreaController]UploadArea Successfully CreateAreaID!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully UploadArea!"})
}

// CreateArea 处理创建区域的请求
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (ac *AreaController) CreateArea(c *fiber.Ctx) error {
	var createAreaData area_model.CreateArea
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&createAreaData); err != nil {
		utils.MsgError("        [AreaController]CreateArea Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 调用 GenerateAreaID 方法生成区域 ID
	AreaID, GenErr := ac.GenerateAreaID(createAreaData.Name)
	if GenErr != nil {
		utils.MsgError("        [AreaController]CreateArea Generate AreaID failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Generate AreaID failed"})
	}
	// 更新区域数据
	err := ac.UpdateAreaData(AreaID, createAreaData.RangeData, createAreaData.RasterSize, createAreaData.RasterIndex,
		createAreaData.RasterData)
	if err != nil {
		utils.MsgError("        [AreaController]CreateArea Insert data failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert data failed"})
	}
	utils.MsgSuccess("        [AreaController]CreateArea Successfully CreateLane!")
	// 返回生成的区域 ID
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully CreateArea!", "data": bson.M{"AreaID": AreaID}})
}

// GetAreaData 获取区域数据
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (ac *AreaController) GetAreaData(c *fiber.Ctx) error {
	var AreaInfo area_model.GetAreaData
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&AreaInfo); err != nil {
		utils.MsgError("        [AreaController]GetAreaData Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 构造 MongoDB 查询过滤器
	filter := bson.M{"AreaID": AreaInfo.AreaID}
	// 构造 MongoDB 查询字段排除器
	drops := bson.M{"_id": 0}
	// 从 MongoDB 中查询数据
	re, err := ac.AreaMongoDB.FindOneWithDropRow("area_collection", filter, drops)
	if err != nil {
		utils.MsgError("        [AreaController]GetAreaData FindOne error!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Mongo find one failed"})
	}
	if re == nil {
		utils.MsgError("        [AreaController]GetAreaData N.A.!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Not Found"})
	}
	// 从 MySQL 中查询数据
	mysqlRe, mysqlErr := ac.AreaMysql.QueryRows(
		fmt.Sprintf("Select * from systemdb.raster_table where AreaID = %d;",
			AreaInfo.AreaID))
	if mysqlErr != nil {
		utils.MsgError("        [AreaController]GetAreaData Query sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "N.A.!"})
	}
	// 构造返回数据
	returnData := bson.M{"AreaID": AreaInfo.AreaID, "AreaData": re, "RasterData": mysqlRe}
	utils.MsgSuccess("        [AreaController]GetAreaData Successfully GetAreaData!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully GetAreaData!", "data": returnData})
}

// UpdateRasterStatus 更新栅格状态
// 参数:
// - AreaID: 区域ID
// - RasterIDs: 栅格ID列表
// - Status: 新状态
// 返回值:
// - error: 错误信息
func (ac *AreaController) UpdateRasterStatus(AreaID int, RasterIDs []int, Status string) error {
	sql := fmt.Sprintf(
		"Update raster_table set Status = '%s' where AreaID = %d and RasterID in %s and Status != '%s';",
		Status, AreaID,
		"("+strings.Trim(strings.Join(strings.Fields(fmt.Sprint(RasterIDs)), ", "), "[]")+")",
		Status,
	)
	_, MysqlErr := ac.AreaMysql.ExecuteCmd(sql)
	return MysqlErr
}

// UpdateRasterDataOcc 更新栅格数据为占用状态
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (ac *AreaController) UpdateRasterDataOcc(c *fiber.Ctx) error {
	var UpdateRasterData area_model.RasterData

	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&UpdateRasterData); err != nil {
		utils.MsgError("        [AreaController]UpdateRasterDataOcc Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 更新栅格状态为占用
	MysqlErr := ac.UpdateRasterStatus(UpdateRasterData.AreaID, UpdateRasterData.RasterIDs, "Occupied")
	if MysqlErr != nil {
		utils.MsgError("        [AreaController]UpdateRasterDataOcc Insert to mysql failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Insert to mysql failed!"})
	}
	utils.MsgSuccess("        [AreaController]UpdateRasterDataOcc Successfully UpdateRasterDataOcc!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully UpdateRasterDataOcc!"})
}

// UpdateRasterDataBan 更新栅格数据为禁用状态
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (ac *AreaController) UpdateRasterDataBan(c *fiber.Ctx) error {
	var UpdateRasterData area_model.RasterData
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&UpdateRasterData); err != nil {
		utils.MsgError("        [AreaController]UpdateRasterDataBan Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 更新栅格状态为禁用
	MysqlErr := ac.UpdateRasterStatus(UpdateRasterData.AreaID, UpdateRasterData.RasterIDs, "Ban")
	if MysqlErr != nil {
		utils.MsgError("        [AreaController]UpdateRasterDataBan Insert to mysql failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Insert to mysql failed!"})
	}
	utils.MsgSuccess("        [AreaController]UpdateRasterDataBan Successfully UpdateRasterDataBan!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully UpdateRasterDataBan!"})
}

func (ac *AreaController) UpdateRasterDataOK(c *fiber.Ctx) error {
	var UpdateRasterData area_model.RasterData
	if err := c.BodyParser(&UpdateRasterData); err != nil {
		utils.MsgError("        [AreaController]UpdateRasterDataOK Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	MysqlErr := ac.UpdateRasterStatus(UpdateRasterData.AreaID, UpdateRasterData.RasterIDs, "OK")
	if MysqlErr != nil {
		utils.MsgError("        [AreaController]UpdateRasterDataOK Insert to mysql failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Insert to mysql failed!"})
	}
	utils.MsgSuccess("        [AreaController]UpdateRasterDataOK Successfully UpdateRasterDataOK!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully UpdateRasterDataOK!"})
}

// UpdateRasterDataBarrier 更新栅格数据为障碍状态
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (ac *AreaController) UpdateRasterDataBarrier(c *fiber.Ctx) error {
	var UpdateRasterData area_model.RasterData
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&UpdateRasterData); err != nil {
		utils.MsgError("        [AreaController]UpdateRasterDataBarrier Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 更新栅格状态为障碍
	MysqlErr := ac.UpdateRasterStatus(UpdateRasterData.AreaID, UpdateRasterData.RasterIDs, "Barrier")
	if MysqlErr != nil {
		utils.MsgError("        [AreaController]UpdateRasterDataBarrier Insert to mysql failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Insert to mysql failed!"})
	}
	utils.MsgSuccess("        [AreaController]UpdateRasterDataBarrier Successfully UpdateRasterDataBarrier!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully UpdateRasterDataBarrier!"})
}

// DeleteAreaData 处理删除区域数据的请求
// 参数:
// - c: fiber.Ctx 上下文
// 返回值:
// - error: 错误信息
func (ac *AreaController) DeleteAreaData(c *fiber.Ctx) error {
	var AreaInfo area_model.GetAreaData
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&AreaInfo); err != nil {
		utils.MsgError("        [AreaController]DeleteAreaData Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	// 构造 MongoDB 查询过滤器
	filter := bson.M{"AreaID": AreaInfo.AreaID}
	// 从 MongoDB 中删除数据
	_, err := ac.AreaMongoDB.DeleteOne("area_collection", filter)
	if err != nil {
		utils.MsgError("        [AreaController]DeleteAreaData DeleteOne error!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Mongo delete one failed"})
	}

	// 从 MySQL 中删除数据
	_, mysqlErr := ac.AreaMysql.ExecuteCmd(
		fmt.Sprintf("DELETE FROM systemdb.raster_table where AreaID = %d;", AreaInfo.AreaID))
	if mysqlErr != nil {
		utils.MsgError("        [AreaController]DeleteAreaData delete mysql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "delete mysql failed!"})
	}

	utils.MsgSuccess("        [AreaController]DeleteAreaData Successfully DeleteAreaData!")
	// 返回成功信息
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully DeleteAreaData!"})
}

func (ac *AreaController) GetRasterData(c *fiber.Ctx) error {
	var GetRasterData area_model.GetRasterData
	if err := c.BodyParser(&GetRasterData); err != nil {
		utils.MsgError("        [AreaController]GetRasterData Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})

	}
	re, err := ac.AreaMysql.QueryRows(
		fmt.Sprintf("Select * from systemdb.raster_table where AreaID = %d and Status = '%s';",
			GetRasterData.AreaID, GetRasterData.Status))
	if err != nil {
		utils.MsgError("        [AreaController]GetRasterData Query sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "N.A.!"})

	}
	// 删除不必要的字段
	delete(re, "AreaID")
	delete(re, "RasterID")
	delete(re, "Status")

	filter := bson.M{"AreaID": GetRasterData.AreaID}
	// 构造 MongoDB 查询字段排除器
	drops := bson.M{"_id": 0, "AreaID": 0, "RasterData": 0, "RasterIndex": 0}
	// 从 MongoDB 中查询数据
	MRe, MErr := ac.AreaMongoDB.FindOneWithDropRow("area_collection", filter, drops)
	if MErr != nil {
		utils.MsgError("        [AreaController]GetRasterData Find raster size error!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Find raster size failed"})
	}
	var reData area_model.GetRasterSizeMongo
	jsonData, _ := json.Marshal(MRe)
	err = json.Unmarshal(jsonData, &reData)
	if err != nil {
		utils.MsgError("        [AreaController]GetRasterData Unmarshal Json data failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Unmarshal Json data failed"})
	}
	re["RasterSize"] = reData.RasterSize
	utils.MsgSuccess("        [AreaController]GetRasterData Successfully GetRasterData!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully GetRasterData!", "data": re})
}
