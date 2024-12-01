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

func (laneModel *LaneController) GetAllLane(c *fiber.Ctx) error {
	var LaneInfo lane_model.FindLaneModel
	if err := c.BodyParser(&LaneInfo); err != nil {
		utils.MsgError("        [LaneController]GetAllLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	filter := bson.M{"IsHide": false}
	drops := bson.M{"_id": 0}
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully return Lane!", "data": re})
}

func (laneModel *LaneController) LaneList(c *fiber.Ctx) error {
	re, err := laneModel.LaneMysql.QueryRows("Select * from lane_table;")
	if err != nil {
		utils.MsgError("        [LaneController]LaneList failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Get data from mysql error"})
	}
	utils.MsgSuccess("        [LaneController]LaneList Successfully GetList!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully return Lane!", "data": re})
}

func (laneModel *LaneController) HideLane(c *fiber.Ctx) error {
	var showHideLaneModel lane_model.ShowHideLaneModel
	if err := c.BodyParser(&showHideLaneModel); err != nil {
		utils.MsgError("        [LaneController]HideLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	_, err := laneModel.LaneMysql.ExecuteCmd(
		fmt.Sprintf("UPDATE lane_table SET IsHide = 1 WHERE LaneID IN %s;",
			"("+strings.Trim(strings.Join(strings.Fields(fmt.Sprint(showHideLaneModel.LaneIDs)), ", "), "[]")+")",
		))
	if err != nil {
		utils.MsgError("        [LaneController]HideLane Change mysql lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mysql failed!"})
	}
	filter := bson.M{"LaneID": bson.M{"$in": showHideLaneModel.LaneIDs}}
	update := bson.M{"$set": bson.M{"IsHide": true}}
	_, MErr := laneModel.LaneMongoDB.Update("lane_data_collection", filter, update)
	if MErr != nil {
		utils.MsgError("        [LaneController]HideLane Change mongo lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mongo failed!"})
	}
	utils.MsgSuccess("        [LaneController]HideLane Successfully HideLane!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully hide Lanes!"})
}

func (laneModel *LaneController) ShowLane(c *fiber.Ctx) error {
	var showHideLaneModel lane_model.ShowHideLaneModel
	if err := c.BodyParser(&showHideLaneModel); err != nil {
		utils.MsgError("        [LaneController]ShowLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	_, err := laneModel.LaneMysql.ExecuteCmd(
		fmt.Sprintf("UPDATE lane_table SET IsHide = 0 WHERE LaneID IN %s;",
			"("+strings.Trim(strings.Join(strings.Fields(fmt.Sprint(showHideLaneModel.LaneIDs)), ", "), "[]")+")",
		))
	if err != nil {
		utils.MsgError("        [LaneController]ShowLane Change mysql lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mysql failed!"})
	}
	filter := bson.M{"LaneID": bson.M{"$in": showHideLaneModel.LaneIDs}}
	update := bson.M{"$set": bson.M{"IsHide": false}}
	_, MErr := laneModel.LaneMongoDB.Update("lane_data_collection", filter, update)
	if MErr != nil {
		utils.MsgError("        [LaneController]ShowLane Change mongo lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mongo failed!"})
	}
	utils.MsgSuccess("        [LaneController]ShowLane Successfully ShowLane!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully show Lanes!"})
}

func (laneModel *LaneController) DeleteLane(c *fiber.Ctx) error {
	var showHideLaneModel lane_model.ShowHideLaneModel
	if err := c.BodyParser(&showHideLaneModel); err != nil {
		utils.MsgError("        [LaneController]DeleteLane Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	_, err := laneModel.LaneMysql.ExecuteCmd(
		fmt.Sprintf("DELETE FROM lane_table WHERE LaneID IN %s;",
			"("+strings.Trim(strings.Join(strings.Fields(fmt.Sprint(showHideLaneModel.LaneIDs)), ", "), "[]")+")",
		))
	if err != nil {
		utils.MsgError("        [LaneController]DeleteLane Change mysql lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mysql failed!"})
	}
	filter := bson.M{"LaneID": bson.M{"$in": showHideLaneModel.LaneIDs}}
	_, MErr := laneModel.LaneMongoDB.Delete("lane_data_collection", filter)
	if MErr != nil {
		utils.MsgError("        [LaneController]DeleteLane Change mongo lane status failed!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Change mongo failed!"})
	}
	utils.MsgSuccess("        [LaneController]DeleteLane Successfully DeleteLane!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully delete Lanes!"})
}
