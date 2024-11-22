package routes

import (
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/controller/aircraft_data_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

// SetupDataFlowRoutes 配置所有路由
func SetupDataFlowRoutes(
	r *fiber.App, kafkaCfg *db_config_model.KafkaConfigModel,
	redisCfg *db_config_model.RedisConfigModel,
) {
	aircraftUploadController := aircraft_data_controller.NewUploadAircraftController(kafkaCfg)
	aircraftReqController := aircraft_data_controller.NewReceiveAircraft(redisCfg)
	// 设置公共路由
	r.Get("/alive", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})

	uploadApis := r.Group("/upload")
	uploadApis.Post("/aircraftData", aircraftUploadController.UploadData)
	uploadApis.Post("/aircraftEvent", aircraftUploadController.UploadEvent)

	recApis := r.Group("/request")
	recApis.Post("/aircraftData", aircraftReqController.RequestAircraftStatus)
	recApis.Post("/aircraftEvent", aircraftReqController.RequestAircraftEvent)
	utils.MsgSuccess("    [SetupDataFlowRoutes]Successfully init!")
}
