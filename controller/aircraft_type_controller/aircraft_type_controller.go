package aircraft_type_controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/aircraft_id_model"
	"uam-power-backend/models/controller_models/aircraft_type_model"
	dbservice "uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

type AircraftTypeController struct {
	TypeMysql *dbservice.MySQLService // MySQL 服务
}

func NewAircraftTypeController(
	MySqlCfg *db_config_model.MySqlConfigModel,
) *AircraftTypeController {
	// 构建 MySQL 连接字符串
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlCfg.Usr, MySqlCfg.Psw, MySqlCfg.Host, MySqlCfg.Port,
		MySqlCfg.DB,
	)
	// 初始化 MySQL 服务
	MysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		utils.MsgError("        [NewAircraftTypeController]Failed to init MySQL service!")
		return nil
	}

	// 打印初始化成功信息
	utils.MsgInfo("        [NewAircraftTypeController]Successfully init!")
	// 返回 AircraftIdController 实例
	return &AircraftTypeController{TypeMysql: MysqlService}
}

func (a *AircraftTypeController) FindAircraftType(Type string) (error, *aircraft_type_model.AircraftTypeMysql) {
	TypeMysqlRe, TypeMysqlErr := a.TypeMysql.QueryRow(
		fmt.Sprintf("Select * from systemdb.aircraft_style_table where Type = '%s';",
			Type))
	if TypeMysqlErr != nil {
		utils.MsgError("        [AircraftTypeController]GetAircraftType Query type sql failed!")
		return TypeMysqlErr, nil
	}
	if TypeMysqlRe == nil {
		utils.MsgError("        [AircraftTypeController]GetAircraftType No such type!")
		return nil, nil
	}
	// 将 MySQL 数据序列化为 JSON
	jsonData, _ := json.Marshal(TypeMysqlRe)
	var mysqlData aircraft_type_model.AircraftTypeMysql
	// 反序列化 JSON 数据
	err := json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [AircraftTypeController]GetAircraftType failed to change!")
		return err, nil
	}
	utils.MsgSuccess("        [AircraftTypeController]GetAircraftType Successfully GetAircraftType!")
	return nil, &mysqlData
}

func (a *AircraftTypeController) GetAircraftType(c *fiber.Ctx) error {
	var TypeInfo aircraft_type_model.GetAircraftTypeInfo
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&TypeInfo); err != nil {
		utils.MsgError("        [AircraftTypeController]GetAircraftType Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	err, mysqlData := a.FindAircraftType(TypeInfo.Type)
	if err != nil {
		utils.MsgError("        [AircraftTypeController]GetAircraftType Query type sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "N.A.!"})
	}
	if mysqlData == nil {
		utils.MsgError("        [AircraftTypeController]GetAircraftType No such type!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A.!"})
	}
	returnData := bson.M{
		"Type": mysqlData.Type, "GLB": mysqlData.GLB,
		"Scale": []float32{mysqlData.ScaleX, mysqlData.ScaleY, mysqlData.ScaleZ},
	}
	utils.MsgSuccess("        [AircraftTypeController]GetAircraftType Successfully GetAircraftType!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully GetAircraftType!", "data": returnData})
}

func (a *AircraftTypeController) CreateAircraftType(c *fiber.Ctx) error {
	var TypeInfo aircraft_type_model.AircraftTypeMysql
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&TypeInfo); err != nil {
		utils.MsgError("        [AircraftTypeController]CreateAircraftType Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	_, mysqlData := a.FindAircraftType(TypeInfo.Type)
	if mysqlData != nil {
		utils.MsgError("        [AircraftTypeController]CreateAircraftType Type already exists!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Type already exists!"})
	}
	_, err := a.TypeMysql.ExecuteCmd(
		fmt.Sprintf("Insert into systemdb.aircraft_style_table values ('%s', '%s', %f, %f, %f);",
			TypeInfo.Type, TypeInfo.GLB, TypeInfo.ScaleX, TypeInfo.ScaleY, TypeInfo.ScaleZ))
	if err != nil {
		utils.MsgError("        [AircraftTypeController]CreateAircraftType Insert type sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Insert failed!"})
	}
	utils.MsgSuccess("        [AircraftTypeController]CreateAircraftType Successfully CreateAircraftType!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully CreateAircraftType!"})
}

func (a *AircraftTypeController) GetTypeByAircraftID(c *fiber.Ctx) error {
	var TypeInfo aircraft_type_model.GetAircraftTypeByID
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&TypeInfo); err != nil {
		utils.MsgError("        [AircraftTypeController]GetTypeByAircraftID Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}

	// 从 MySQL 获取数据
	mysqlRe, mysqlErr := a.TypeMysql.QueryRow(
		fmt.Sprintf(
			"Select * from systemdb.aircraft_identity_table where AircraftID = %d;",
			TypeInfo.AircraftID,
		))
	if mysqlErr != nil {
		// 打印错误信息并返回错误响应
		utils.MsgError("        [NewAircraftIdController]GetTypeByAircraftID No such Aircraft!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "N.A.!"})
	}
	var IDMysqlData aircraft_id_model.MysqlAircraftInfo
	// 将 MySQL 数据序列化为 JSON
	jsonData, _ := json.Marshal(mysqlRe)
	IDErr := json.Unmarshal(jsonData, &IDMysqlData)
	if IDErr != nil {
		utils.MsgError("        [AircraftTypeController]GetTypeByAircraftID failed to find!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "ID N.A.!"})
	}

	err, mysqlData := a.FindAircraftType(IDMysqlData.Type)
	if err != nil {
		utils.MsgError("        [AircraftTypeController]GetTypeByAircraftID Query type sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "N.A.!"})
	}
	if mysqlData == nil {
		utils.MsgError("        [AircraftTypeController]GetTypeByAircraftID No such type!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Type N.A.!"})
	}
	returnData := bson.M{
		"Type": mysqlData.Type, "GLB": mysqlData.GLB,
		"Scale": []float32{mysqlData.ScaleX, mysqlData.ScaleY, mysqlData.ScaleZ},
	}
	utils.MsgSuccess("        [AircraftTypeController]GetTypeByAircraftID Successfully GetAircraftType!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully GetAircraftType!", "data": returnData})
}

func (a *AircraftTypeController) DeleteType(c *fiber.Ctx) error {
	var TypeInfo aircraft_type_model.GetAircraftTypeInfo
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&TypeInfo); err != nil {
		utils.MsgError("        [AircraftTypeController]DeleteType Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	_, mysqlData := a.FindAircraftType(TypeInfo.Type)
	if mysqlData == nil {
		utils.MsgError("        [AircraftTypeController]DeleteType No such type!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Type N.A.!"})
	}
	_, err := a.TypeMysql.ExecuteCmd(
		fmt.Sprintf("Delete from systemdb.aircraft_style_table where Type = '%s';",
			TypeInfo.Type))
	if err != nil {
		utils.MsgError("        [AircraftTypeController]DeleteType Delete type sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Delete failed!"})
	}
	utils.MsgSuccess("        [AircraftTypeController]DeleteType Successfully DeleteType!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully DeleteType!"})
}

func (a *AircraftTypeController) ChangeType(c *fiber.Ctx) error {
	var TypeInfo aircraft_type_model.ChangeAircraftTypeModel
	// 解析请求体中的 JSON 数据
	if err := c.BodyParser(&TypeInfo); err != nil {
		utils.MsgError("        [AircraftTypeController]ChangeType Invalid Request JSON data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	_, mysqlData := a.FindAircraftType(TypeInfo.Type)
	if mysqlData == nil {
		utils.MsgError("        [AircraftTypeController]ChangeType No such type!")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Type N.A.!"})
	}
	sql := fmt.Sprintf("Update systemdb.aircraft_style_table set GLB = '%s', ScaleX = %f, ScaleY = %f, ScaleZ = %f where Type = '%s';",
		TypeInfo.GLB, TypeInfo.ScaleX, TypeInfo.ScaleY, TypeInfo.ScaleZ, TypeInfo.Type)
	_, err := a.TypeMysql.ExecuteCmd(sql)
	if err != nil {
		utils.MsgError("        [AircraftTypeController]ChangeType Change type sql failed!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"msg": "Change failed!"})
	}
	utils.MsgSuccess("        [AircraftTypeController]ChangeType Successfully ChangeType!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully ChangeType!"})
}
