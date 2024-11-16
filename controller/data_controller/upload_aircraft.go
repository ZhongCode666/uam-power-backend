package data_controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/data_flow_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

type UploadAircraftController struct {
	kafkaStatusService *dbservice.KafkaProducer
	kafkaEventService  *dbservice.KafkaProducer
}

func NewUploadAircraftController(kafkaConfig *db_config_model.KafkaConfigModel) *UploadAircraftController {
	kafkaStatusService := dbservice.NewKafkaProducer(kafkaConfig.Addr, kafkaConfig.AircraftDataTopic)
	kafkaEventService := dbservice.NewKafkaProducer(kafkaConfig.Addr, kafkaConfig.AircraftEventTopic)
	utils.MsgSuccess("        [UploadAircraftController]init successfully!")
	return &UploadAircraftController{
		kafkaStatusService: kafkaStatusService,
		kafkaEventService:  kafkaEventService,
	}
}

func (controller *UploadAircraftController) UploadData(c *gin.Context) {
	var aircraftData data_flow_model.AircraftStatus
	if err := c.ShouldBindJSON(&aircraftData); err != nil {
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid JSON data rec >" + err.Error())
		c.JSON(400, gin.H{"msg": "Invalid JSON data"})
		return
	}
	if !utils.IsValidSqlTimeFormat(aircraftData.TimeString) {
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid time format")
		c.JSON(403, gin.H{"msg": "Invalid time format"})
		return
	}
	jStr, err := json.Marshal(aircraftData)
	if err != nil {
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid JSON data tran_str >" + err.Error())
		c.JSON(400, gin.H{"msg": "Invalid JSON data"})
		return
	}
	err = controller.kafkaStatusService.SendMessage(string(jStr))
	if err != nil {
		utils.MsgInfo(err.Error())
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid JSON data >" + err.Error())
		c.JSON(400, gin.H{"msg": "Invalid JSON data"})
		return
	}
	utils.MsgSuccess("        [UploadAircraftController]UploadData successfully!")
	c.JSON(200, gin.H{"msg": "Successfully send to Kafka!"})
}

func (controller *UploadAircraftController) UploadEvent(c *gin.Context) {
	var aircraftEvent data_flow_model.AircraftEvent

	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&aircraftEvent); err != nil {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid JSON data >" + err.Error())
		c.JSON(400, gin.H{"msg": "Invalid JSON data"})
		return
	}
	if !utils.IsValidSqlTimeFormat(aircraftEvent.TimeString) {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid time format")
		c.JSON(403, gin.H{"msg": "Invalid time format"})
		return
	}
	jStr, err := json.Marshal(aircraftEvent)
	if err != nil {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid JSON data >" + err.Error())
		c.JSON(400, gin.H{"msg": "Invalid JSON data"})
		return
	}
	err = controller.kafkaEventService.SendMessage(string(jStr))
	if err != nil {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid JSON data >" + err.Error())
		c.JSON(400, gin.H{"msg": "Invalid JSON data"})
		return
	}
	utils.MsgSuccess("        [UploadAircraftController]UploadEvent successfully!")
	c.JSON(200, gin.H{"msg": "Successfully send to Kafka!"})
}

func (controller *UploadAircraftController) Close() {
	// 返回响应
	err := controller.kafkaEventService.Close()
	if err != nil {
		return
	}
	err = controller.kafkaStatusService.Close()
	if err != nil {
		return
	}
	utils.MsgSuccess("        [UploadAircraftController]Close successfully!")
}
