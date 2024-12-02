package routes

import (
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/controller/aircraft_data_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

// SetupUploadFlowRoutes 设置上传数据流路由
// r 是 Fiber 应用实例
// kafkaCfg 是 Kafka 配置模型
func SetupUploadFlowRoutes(
	r *fiber.App, kafkaCfg *db_config_model.KafkaConfigModel,
) {
	// 创建一个新的 UploadAircraftController 实例
	aircraftUploadController := aircraft_data_controller.NewUploadAircraftController(kafkaCfg)

	// 设置公共路由
	r.Get("/alive", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})

	// 创建一个新的路由组 /upload
	uploadApis := r.Group("/upload")
	uploadApis.Post("/aircraftData", aircraftUploadController.UploadData)   // 上传 Aircraft 数据
	uploadApis.Post("/aircraftEvent", aircraftUploadController.UploadEvent) // 上传 Aircraft 事件

	// 打印成功初始化信息
	utils.MsgSuccess("    [SetupDataFlowRoutes]Successfully init!")
}

// SetupReceiveFlowRoutes 设置接收数据流路由
// r 是 Fiber 应用实例
// redisCfg 是 Redis 配置模型
// mysqlCfg 是 MySQL 配置模型
func SetupReceiveFlowRoutes(
	r *fiber.App,
	redisCfg *db_config_model.RedisConfigModel, mysqlCfg *db_config_model.MySqlConfigModel,
) {
	// 创建一个新的 ReceiveAircraft 实例
	aircraftReqController := aircraft_data_controller.NewReceiveAircraft(redisCfg, mysqlCfg)

	// 设置公共路由
	r.Get("/alive", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})

	// 创建一个新的路由组 /request
	recApis := r.Group("/request")
	recApis.Post("/aircraftData", aircraftReqController.RequestAircraftStatus)             // 请求 Aircraft 状态
	recApis.Post("/aircraftEvent", aircraftReqController.RequestAircraftEvent)             // 请求 Aircraft 事件
	recApis.Get("/activeAircraftData", aircraftReqController.ReceiveActiveData)            // 接收活动数据
	recApis.Get("/activeAircraftEvent", aircraftReqController.ReceiveActiveEvent)          // 接收活动事件
	recApis.Get("/activeAircraftIDs", aircraftReqController.ReceiveActiveIDs)              // 接收活动 ID
	recApis.Post("/activeAircraftHistoryData", aircraftReqController.ReceiveHistoryData)   // 接收历史数据
	recApis.Post("/activeAircraftHistoryEvent", aircraftReqController.ReceiveHistoryEvent) // 接收历史事件

	// 打印成功初始化信息
	utils.MsgSuccess("    [SetupDataFlowRoutes]Successfully init!")
}
