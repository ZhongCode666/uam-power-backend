package routes

import (
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/controller/aircraft_task_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

func SetupAircraftTaskRoutes(
	r *fiber.App, RedisCfg *db_config_model.RedisConfigModel,
	MySqlCfg *db_config_model.MySqlConfigModel,
	ClickHouse *db_config_model.ClickHouseConfigModel,
) {
	//aircraftTaskController := aircraft_task_controller.NewAircraftTaskModel(RedisCfg, MySqlCfg)
	aircraftTaskController := aircraft_task_controller.NewAircraftTaskModelClickHouse(RedisCfg, MySqlCfg, ClickHouse)
	uploadApis := r.Group("/aircraftTask")
	uploadApis.Post("/end", aircraftTaskController.EndTask)
	uploadApis.Post("/create", aircraftTaskController.CreateTask)
	uploadApis.Post("/check", aircraftTaskController.CheckTaskInfo)
	utils.MsgSuccess("    [SetupAircraftTaskRoutes]Successfully init!")
}
