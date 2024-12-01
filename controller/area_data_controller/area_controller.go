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

type AreaController struct {
	AreaMongoDB *dbservice.MongoDBClient
	AreaMysql   *dbservice.MySQLService
}

func NewAreaController(
	MongoCfg *db_config_model.MongoConfigModel, MySqlCfg *db_config_model.MySqlConfigModel,
) *AreaController {
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlCfg.Usr, MySqlCfg.Psw, MySqlCfg.Host, MySqlCfg.Port,
		MySqlCfg.DB,
	)
	MysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		utils.MsgError("        [AreaController]ini mysql service error:" + err.Error())
		return nil
	}
	utils.MsgSuccess("        [AreaController]Successfully AreaMysql!")
	mongoLink := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		MongoCfg.Usr, MongoCfg.Psw, MongoCfg.Host, MongoCfg.Port,
	)
	mongoService, MongoErr := dbservice.NewMongoDBClient(mongoLink, MongoCfg.AreaDB)
	if MongoErr != nil {
		utils.MsgError("        [AreaController]MongoDB client error:" + MongoErr.Error())
	}
	utils.MsgSuccess("        [AreaController]Successfully AreaMongo!")
	return &AreaController{
		AreaMysql:   MysqlService,
		AreaMongoDB: mongoService,
	}
}

func (ac *AreaController) CreateArea(c *fiber.Ctx) error {
	curStr := utils.GetTimeStr()
	randStr := curStr + "-" + utils.GetUniqueStr()
	var createAreaData area_model.CreateArea
	if err := c.BodyParser(&createAreaData); err != nil {
		utils.MsgError("        [AreaController]CreateArea Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	_, err := ac.AreaMysql.ExecuteCmd(
		fmt.Sprintf("INSERT INTO systemdb.area_table(Name, TimeStr) VALUES ('%s', '%s');",
			createAreaData.Name, randStr,
		))
	if err != nil {
		utils.MsgError("        [AreaController]CreateArea Create Area Failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert MySql failed"})
	}
	mysqlRe, mysqlErr := ac.AreaMysql.QueryRow(
		fmt.Sprintf("Select * from systemdb.area_table where TimeStr = '%s';",
			randStr))
	if mysqlErr != nil {
		utils.MsgError("        [AreaController]CreateArea Query sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "N.A.!"})
	}
	jsonData, _ := json.Marshal(mysqlRe)
	var mysqlData area_model.AreaMysqlModel
	err = json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [AreaController]CreateArea Unmarshal Json data failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Invalid MySql Data"})
	}
	document := bson.M{
		"AreaID": mysqlData.AreaID, "RangeData": createAreaData.RangeData,
		"RasterSize": createAreaData.RasterSize, "RasterIndex": createAreaData.RasterIndex,
	}
	_, err = ac.AreaMongoDB.InsertOne("area_collection", document)
	if err != nil {
		utils.MsgError("        [AreaController]CreateArea Insert Mongo failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert Mongo failed"})
	}
	var values []string
	for id, record := range createAreaData.RasterData {
		status := record.Status
		longitude := record.Longitude
		latitude := record.Latitude
		altitude := record.Altitude

		// 构造每一行的值
		values = append(
			values,
			fmt.Sprintf(
				"(%d, %d, %.12f, %.12f, %.12f, '%s')",
				id, mysqlData.AreaID, longitude, latitude, altitude, status,
			))
	}

	// 拼接最终 SQL
	sql := fmt.Sprintf(
		"INSERT INTO systemdb.raster_table(RasterID, AreaID, Longitude, Latitude, Altitude, Status) VALUES %s;",
		strings.Join(values, ", "))
	_, MysqlErr := ac.AreaMysql.ExecuteCmd(sql)
	if MysqlErr != nil {
		utils.MsgError("        [AreaController]CreateArea Insert Raster to Mysql failed err>" + MysqlErr.Error())
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert Raster to Mysql failed"})
	}
	utils.MsgSuccess("        [AreaController]CreateArea Successfully CreateLane!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully CreateArea!", "data": mysqlRe})
}

func (ac *AreaController) GetAreaData(c *fiber.Ctx) error {
	var AreaInfo area_model.GetAreaData
	if err := c.BodyParser(&AreaInfo); err != nil {
		utils.MsgError("        [AreaController]GetAreaData Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	filter := bson.M{"AreaID": AreaInfo.AreaID}
	drops := bson.M{"_id": 0}
	re, err := ac.AreaMongoDB.FindOneWithDropRow("area_collection", filter, drops)
	if err != nil {
		utils.MsgError("        [AreaController]GetAreaData FindOne error!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Mongo find one failed"})
	}
	if re == nil {
		utils.MsgError("        [AreaController]GetAreaData N.A.!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Not Found"})
	}
	mysqlRe, mysqlErr := ac.AreaMysql.QueryRows(
		fmt.Sprintf("Select * from systemdb.raster_table where AreaID = %d;",
			AreaInfo.AreaID))
	if mysqlErr != nil {
		utils.MsgError("        [AreaController]GetAreaData Query sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "N.A.!"})
	}
	returnData := bson.M{"AreaID": AreaInfo.AreaID, "AreaData": re, "RasterData": mysqlRe}
	utils.MsgSuccess("        [AreaController]GetAreaData Successfully GetAreaData!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully GetAreaData!", "data": returnData})
}

func (ac *AreaController) UpdateRasterData(c *fiber.Ctx) error {
	var UpdateRasterData area_model.RasterData
	if err := c.BodyParser(&UpdateRasterData); err != nil {
		utils.MsgError("        [AreaController]UpdateRasterData Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	sql := fmt.Sprintf(
		"Update raster_table set Status = 'Occupied' where AreaID = %d and RasterID in %s;",
		UpdateRasterData.AreaID,
		"("+strings.Trim(strings.Join(strings.Fields(fmt.Sprint(UpdateRasterData.RasterIDs)), ", "), "[]")+")",
	)
	//utils.MsgSuccess(sql)
	_, MysqlErr := ac.AreaMysql.ExecuteCmd(sql)
	if MysqlErr != nil {
		utils.MsgError("        [AreaController]UpdateRasterData Update Raster to Mysql failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Update Raster to Mysql failed"})
	}
	utils.MsgSuccess("        [AreaController]UpdateRasterData Successfully UpdateRasterData!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully UpdateRasterData!"})
}

func (ac *AreaController) DeleteAreaData(c *fiber.Ctx) error {
	var AreaInfo area_model.GetAreaData
	if err := c.BodyParser(&AreaInfo); err != nil {
		utils.MsgError("        [AreaController]DeleteAreaData Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	filter := bson.M{"AreaID": AreaInfo.AreaID}
	_, err := ac.AreaMongoDB.DeleteOne("area_collection", filter)
	if err != nil {
		utils.MsgError("        [AreaController]DeleteAreaData DeleteOne error!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Mongo delete one failed"})
	}

	_, mysqlErr := ac.AreaMysql.ExecuteCmd(
		fmt.Sprintf("DELETE FROM systemdb.raster_table where AreaID = %d;",
			AreaInfo.AreaID))
	if mysqlErr != nil {
		utils.MsgError("        [AreaController]DeleteAreaData delete mysql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "delete mysql failed!"})
	}

	utils.MsgSuccess("        [AreaController]DeleteAreaData Successfully DeleteAreaData!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully DeleteAreaData!"})
}
