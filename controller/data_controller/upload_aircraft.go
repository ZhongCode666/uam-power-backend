package data_controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
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

func (controller *UploadAircraftController) UploadData(c *fiber.Ctx) error {
	var aircraftData data_flow_model.AircraftStatus
	if err := c.BodyParser(&aircraftData); err != nil {
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid JSON data rec >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	if !utils.IsValidSqlTimeFormat(aircraftData.TimeString) {
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid time format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid time format"})
	}
	jStr, err := json.Marshal(aircraftData)
	if err != nil {
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid JSON data tran_str >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	err = controller.kafkaStatusService.SendMessage(string(jStr))
	if err != nil {
		utils.MsgInfo(err.Error())
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid JSON data >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	utils.MsgSuccess("        [UploadAircraftController]UploadData successfully!")
	return c.Status(fiber.StatusOK).JSON(gin.H{"msg": "Successfully send to Kafka!"})
}

func (controller *UploadAircraftController) UploadEvent(c *fiber.Ctx) error {
	var aircraftEvent data_flow_model.AircraftEvent

	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftEvent); err != nil {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid JSON data >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	if !utils.IsValidSqlTimeFormat(aircraftEvent.TimeString) {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid time format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid time format"})
	}
	jStr, err := json.Marshal(aircraftEvent)
	if err != nil {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid JSON data >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	err = controller.kafkaEventService.SendMessage(string(jStr))
	if err != nil {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid JSON data >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}
	utils.MsgSuccess("        [UploadAircraftController]UploadEvent successfully!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Successfully send to Kafka!"})
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
