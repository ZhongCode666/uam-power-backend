package aircraft_id_controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/aircraft_id_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

type AircraftIdController struct {
	IDMySql   *dbservice.MySQLService
	RedisInfo *dbservice.RedisDict
}

func NewAircraftIdController(
	MySqlCfg *db_config_model.MySqlConfigModel, RedisCfg *db_config_model.RedisConfigModel,
) *AircraftIdController {
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlCfg.Usr, MySqlCfg.Psw, MySqlCfg.Host, MySqlCfg.Port,
		MySqlCfg.DB,
	)
	MysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		return nil
	}
	RedisInfo := dbservice.NewRedisDict(RedisCfg.Host, RedisCfg.Port, RedisCfg.AircraftDBno)
	utils.MsgInfo("        [NewAircraftIdController]Successfully init!")
	return &AircraftIdController{IDMySql: MysqlService, RedisInfo: RedisInfo}
}

func (a *AircraftIdController) GetAircraftInfo(c *gin.Context) {
	var RequestID aircraft_id_model.GetAircraftInfoID
	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&RequestID); err != nil {
		utils.MsgError("        [NewAircraftIdController]GetAircraftInfo Request Invalid JSON data!")
		c.JSON(400, gin.H{"msg": "Invalid JSON data"})
		return
	}
	re, _ := a.RedisInfo.Get(strconv.Itoa(RequestID.AircraftID))
	if re != nil {
		utils.MsgSuccess("        [NewAircraftIdController]GetAircraftInfo Hit Redis auto Return!")
		c.JSON(200, gin.H{"msg": "Successfully GetAircraftInfo!", "data": re})
		return
	}
	mysqlRe, mysqlErr := a.IDMySql.QueryRow(fmt.Sprintf("Select * from systemdb.aircraft_identity_table where AircraftID = %d;", RequestID.AircraftID))
	if mysqlErr != nil {
		utils.MsgError("        [NewAircraftIdController]GetAircraftInfo No such Aircraft!")
		c.JSON(404, gin.H{"msg": "N.A.!"})
		return
	}
	jsonData, _ := json.Marshal(mysqlRe)
	err := a.RedisInfo.Set(strconv.Itoa(RequestID.AircraftID), string(jsonData))
	if err != nil {
		utils.MsgError("        [NewAircraftIdController]Set Redis Failed!")
		c.JSON(403, gin.H{"msg": "Redis failed!"})
		return
	}
	utils.MsgSuccess("        [NewAircraftIdController]Successfully GetAircraftInfo!")
	c.JSON(200, gin.H{"msg": "Successfully GetAircraftInfo!", "data": mysqlRe})
}

func (a *AircraftIdController) CreateUser(c *gin.Context) {
	curStr := utils.GetTimeStr()
	var RequestInfo aircraft_id_model.SetAircraftInfo
	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&RequestInfo); err != nil {
		utils.MsgError("        [NewAircraftIdController]CreateUser invalid Requests Json!")
		c.JSON(400, gin.H{"msg": "Invalid JSON data"})
		return
	}
	_, err := a.IDMySql.ExecuteCmd(
		fmt.Sprintf("INSERT INTO systemdb.aircraft_identity_table(Type, Company, Name, TimeStr) VALUES ('%s', '%s', '%s', '%s')",
			RequestInfo.Type, RequestInfo.Company, RequestInfo.Name, curStr,
		))
	if err != nil {
		utils.MsgError("        [NewAircraftIdController]CreateUser failed to Mysql!")
		c.JSON(403, gin.H{"msg": "Send to Mysql Failed"})
		return
	}
	mysqlRe, mysqlErr := a.IDMySql.QueryRow(
		fmt.Sprintf("Select * from systemdb.aircraft_identity_table where TimeStr = '%s';",
			curStr))
	if mysqlErr != nil {
		utils.MsgError("        [NewAircraftIdController]CreateUser data in MySql not found!")
		c.JSON(404, gin.H{"msg": "N.A.!"})
		return
	}
	jsonData, _ := json.Marshal(mysqlRe)
	var mysqlData aircraft_id_model.MysqlAircraftInfo
	err = json.Unmarshal(jsonData, &mysqlData)
	if err != nil {
		utils.MsgError("        [NewAircraftIdController]CreateUser failed to redis!")
		c.JSON(403, gin.H{"msg": "failed to send Redis!"})
		return
	}
	err = a.RedisInfo.Set(strconv.Itoa(mysqlData.AircraftID), string(jsonData))
	if err != nil {
		utils.MsgError("        [NewAircraftIdController]CreateUser failed to redis!")
		c.JSON(403, gin.H{"msg": "failed to send Redis!"})
		return
	}
	utils.MsgSuccess("        [NewAircraftIdController]Successfully CreateUser!")
	c.JSON(200, gin.H{"msg": "Successfully CreateUser!", "data": mysqlRe})
}
