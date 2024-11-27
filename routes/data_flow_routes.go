package routes

import (
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/controller/aircraft_data_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

func SetupUploadFlowRoutes(
	r *fiber.App, kafkaCfg *db_config_model.KafkaConfigModel,
) {
	aircraftUploadController := aircraft_data_controller.NewUploadAircraftController(kafkaCfg)
	// 设置公共路由
	r.Get("/alive", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})

	uploadApis := r.Group("/upload")
	uploadApis.Post("/aircraftData", aircraftUploadController.UploadData)
	uploadApis.Post("/aircraftEvent", aircraftUploadController.UploadEvent)

	utils.MsgSuccess("    [SetupDataFlowRoutes]Successfully init!")
}

func SetupReceiveFlowRoutes(
	r *fiber.App,
	redisCfg *db_config_model.RedisConfigModel, mysqlCfg *db_config_model.MySqlConfigModel,
) {
	aircraftReqController := aircraft_data_controller.NewReceiveAircraft(redisCfg, mysqlCfg)
	// 设置公共路由
	r.Get("/alive", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})

	recApis := r.Group("/request")
	recApis.Post("/aircraftData", aircraftReqController.RequestAircraftStatus)
	recApis.Post("/aircraftEvent", aircraftReqController.RequestAircraftEvent)
	recApis.Get("/activeAircraftData", aircraftReqController.ReceiveActiveData)
	recApis.Get("/activeAircraftEvent", aircraftReqController.ReceiveActiveEvent)
	recApis.Get("/activeAircraftIDs", aircraftReqController.ReceiveActiveIDs)
	recApis.Post("/activeAircraftHistoryData", aircraftReqController.ReceiveHistoryData)
	recApis.Post("/activeAircraftHistoryEvent", aircraftReqController.ReceiveHistoryEvent)
	utils.MsgSuccess("    [SetupDataFlowRoutes]Successfully init!")
}
