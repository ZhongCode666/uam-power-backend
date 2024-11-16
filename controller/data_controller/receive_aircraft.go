package data_controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/data_flow_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

type RequestAircraft struct {
	StatusRedisService *dbservice.RedisDict
	EventRedisService  *dbservice.RedisDict
}

func NewReceiveAircraft(redisConfig *db_config_model.RedisConfigModel) *RequestAircraft {
	utils.MsgSuccess("        [ReceiveAircraft]init successfully!")
	redisStatusService := dbservice.NewRedisDict(redisConfig.Host, redisConfig.Port, redisConfig.StatusDBno)
	redisEventService := dbservice.NewRedisDict(redisConfig.Host, redisConfig.Port, redisConfig.EventDBno)
	return &RequestAircraft{StatusRedisService: redisStatusService, EventRedisService: redisEventService}
}

func (receiver *RequestAircraft) RequestAircraftStatus(c *gin.Context) {
	var aircraftReq data_flow_model.RecAircraftStatusRequest

	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&aircraftReq); err != nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftStatus Invalid JSON data! >" + err.Error())
		c.JSON(400, gin.H{"msg": "Invalid JSON data"})
		return
	}

	rec, err := receiver.StatusRedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil || rec == nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftStatus Invalid JSON data!")
		c.JSON(404, gin.H{"msg": "N.A.!"})
		return
	}
	utils.MsgError("        [ReceiveAircraft]RequestAircraftStatus Invalid JSON data!")
	c.JSON(200, gin.H{"msg": "Successfully requestData!", "data": rec})
	return
}

func (receiver *RequestAircraft) RequestAircraftEvent(c *gin.Context) {
	var aircraftReq data_flow_model.RecAircraftStatusRequest

	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&aircraftReq); err != nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftEvent Invalid JSON data! >" + err.Error())
		c.JSON(400, gin.H{"msg": "Invalid JSON data"})
		return
	}

	rec, err := receiver.EventRedisService.Get(strconv.Itoa(aircraftReq.AircraftID))
	if err != nil || rec == nil {
		utils.MsgError("        [ReceiveAircraft]RequestAircraftEvent Invalid JSON data!")
		c.JSON(404, gin.H{"msg": "N.A.!"})
		return
	}
	utils.MsgSuccess("        [ReceiveAircraft]RequestAircraftEvent Successfully requestData!")
	c.JSON(200, gin.H{"msg": "Successfully requestData!", "data": rec})
	return
}
