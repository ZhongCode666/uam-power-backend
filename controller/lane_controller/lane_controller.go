package lane_controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/lane_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

type LaneController struct {
	LaneMongoDB *dbservice.MongoDBClient
	LaneMysql   *dbservice.MySQLService
}

func NewLaneController(
	MongoCfg *db_config_model.MongoConfigModel, MySqlCfg *db_config_model.MySqlConfigModel,
) *LaneController {
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlCfg.Usr, MySqlCfg.Psw, MySqlCfg.Host, MySqlCfg.Port,
		MySqlCfg.DB,
	)
	MysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		utils.MsgError("        [LaneController]ini mysql service error:" + err.Error())
		return nil
	}
	utils.MsgSuccess("        [LaneController]Successfully LaneMysql!")
	mongoLink := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		MongoCfg.Usr, MongoCfg.Psw, MongoCfg.Host, MongoCfg.Port,
	)
	mongoService, MongoErr := dbservice.NewMongoDBClient(mongoLink, MongoCfg.AreaDB)
	if MongoErr != nil {
		utils.MsgError("        [LaneController]MongoDB client error:" + MongoErr.Error())
	}
	utils.MsgSuccess("        [LaneController]Successfully LaneMongo!")
	return &LaneController{
		LaneMysql:   MysqlService,
		LaneMongoDB: mongoService,
	}
}

func (laneModel *LaneController) CreateLane(c *fiber.Ctx) error {
	curStr := utils.GetTimeStr()
	randStr := curStr + "-" + utils.GetUniqueStr()
	var LaneInfo lane_model.LaneCreateModel
	if err := c.BodyParser(&LaneInfo); err != nil {
		utils.MsgError("        [LaneController]CreateLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	_, err := laneModel.LaneMysql.ExecuteCmd(
		fmt.Sprintf("INSERT INTO systemdb.lane_table(IsHide, Name, TimeStr) VALUES (%t, '%s', '%s');",
			LaneInfo.IsHide, LaneInfo.Name, randStr,
		))
	if err != nil {
		utils.MsgError("        [LaneController]CreateLane Create Lane Failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert failed"})
	}
	mysqlRe, mysqlErr := laneModel.LaneMysql.QueryRow(
		fmt.Sprintf("Select * from systemdb.lane_table where TimeStr = '%s';",
			randStr))
	if mysqlErr != nil {
		utils.MsgError("        [LaneController]CreateLane Query sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "N.A.!"})
	}
	jsonData, _ := json.Marshal(mysqlRe)
	var mysqlData lane_model.LaneMysqlData
	err = json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [LaneController]CreateLane Unmarshal Json data failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Invalid Data"})
	}
	document := bson.M{
		"LaneID": mysqlData.LaneID, "PointData": LaneInfo.PointData,
		"RasterData": LaneInfo.RasterData, "IsHide": mysqlData.IsHide == 1,
	}
	_, err = laneModel.LaneMongoDB.InsertOne("lane_data_collection", document)
	if err != nil {
		utils.MsgError("        [LaneController]CreateLane Insert Mongo failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert Mongo failed"})
	}
	utils.MsgSuccess("        [LaneController]CreateLane Successfully CreateLane!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully CreateLane!", "data": mysqlRe})
}

func (laneModel *LaneController) GetLane(c *fiber.Ctx) error {
	var LaneInfo lane_model.FindLaneModel
	if err := c.BodyParser(&LaneInfo); err != nil {
		utils.MsgError("        [LaneController]GetLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	filter := bson.M{"LaneID": LaneInfo.LaneID, "IsHide": false}
	drops := bson.M{"_id": 0}
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully return Lane!", "data": re})
}
