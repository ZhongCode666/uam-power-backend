package aircraft_data_controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/controller_models/data_flow_model"
	"uam-power-backend/service/db_service"
	"uam-power-backend/utils"
)

// UploadAircraftController 结构体表示上传飞机数据的控制器
type UploadAircraftController struct {
	kafkaStatusService *dbservice.KafkaManager // Kafka 状态服务
	kafkaEventService  *dbservice.KafkaManager // Kafka 事件服务
}

// NewUploadAircraftController 创建并返回一个新的 UploadAircraftController 实例
func NewUploadAircraftController(kafkaConfig *db_config_model.KafkaConfigModel) *UploadAircraftController {
	kafkaStatusService := dbservice.NewKafkaManager(kafkaConfig.Addr, kafkaConfig.AircraftDataTopic, kafkaConfig.NumDataProducers)
	kafkaEventService := dbservice.NewKafkaManager(kafkaConfig.Addr, kafkaConfig.AircraftEventTopic, kafkaConfig.NumEventProducers)
	utils.MsgSuccess("        [UploadAircraftController]init successfully!")
	return &UploadAircraftController{
		kafkaStatusService: kafkaStatusService,
		kafkaEventService:  kafkaEventService,
	}
}

// UploadData 处理上传飞机数据的请求
// @param c *fiber.Ctx 请求上下文
// @return error 返回错误信息
func (controller *UploadAircraftController) UploadData(c *fiber.Ctx) error {
	// 定义 aircraftData 变量并解析请求体中的 JSON 数据
	var aircraftData data_flow_model.AircraftStatus
	if err := c.BodyParser(&aircraftData); err != nil {
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid JSON data rec >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}

	// 验证时间格式是否有效
	if !utils.IsValidSqlTimeFormat(aircraftData.TimeString) {
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid time format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid time format"})
	}

	// 将 aircraftData 转换为 JSON 字符串
	jStr, err := json.Marshal(aircraftData)
	if err != nil {
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid JSON data tran_str >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}

	// 通过 Kafka 发送 JSON 数据
	err = controller.kafkaStatusService.SendMsgRoundRobin(string(jStr))
	if err != nil {
		utils.MsgInfo(err.Error())
		utils.MsgError("        [UploadAircraftController]UploadData error-Invalid JSON data >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid JSON data"})
	}

	// 返回成功响应
	utils.MsgSuccess("        [UploadAircraftController]UploadData successfully!")
	return c.Status(fiber.StatusOK).JSON(gin.H{"msg": "Successfully send to Kafka!"})
}

// UploadEvent 处理上传飞机事件的请求
// @param c *fiber.Ctx 请求上下文
// @return error 返回错误信息
func (controller *UploadAircraftController) UploadEvent(c *fiber.Ctx) error {
	var aircraftEvent data_flow_model.AircraftEvent

	// 绑定 JSON 数据到结构体
	if err := c.BodyParser(&aircraftEvent); err != nil {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid JSON data >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "无效的 JSON 数据"})
	}

	// 验证时间格式是否有效
	if !utils.IsValidSqlTimeFormat(aircraftEvent.TimeString) {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid time format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "无效的时间格式"})
	}

	// 将 aircraftEvent 转换为 JSON 字符串
	jStr, err := json.Marshal(aircraftEvent)
	if err != nil {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid JSON data >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "无效的 JSON 数据"})
	}

	// 通过 Kafka 发送 JSON 数据
	err = controller.kafkaEventService.SendMsgRoundRobin(string(jStr))
	if err != nil {
		utils.MsgError("        [UploadAircraftController]UploadEvent error-Invalid JSON data >" + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "无效的 JSON 数据"})
	}

	// 返回成功响应
	utils.MsgSuccess("        [UploadAircraftController]UploadEvent successfully!")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "成功发送到 Kafka!"})
}

// Close 关闭 Kafka 服务
func (controller *UploadAircraftController) Close() {
	// 关闭 Kafka 事件服务
	err := controller.kafkaEventService.Close()
	if err != nil {
		utils.MsgError("        [UploadAircraftController]Close error-Event service >" + err.Error())
		return
	}
	// 关闭 Kafka 状态服务
	err = controller.kafkaStatusService.Close()
	if err != nil {
		utils.MsgError("        [UploadAircraftController]Close error-Status service >" + err.Error())
		return
	}
	utils.MsgSuccess("        [UploadAircraftController]Close successfully!")
}
