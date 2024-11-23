package routes

import (
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/controller/lane_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

func SetupLaneRoutes(
	r *fiber.App,
	mongoCfg *db_config_model.MongoConfigModel,
	mysqlCfg *db_config_model.MySqlConfigModel,
) {
	laneController := lane_controller.NewLaneController(mongoCfg, mysqlCfg)

	recApis := r.Group("/lane")
	recApis.Post("/create", laneController.CreateLane)
	recApis.Post("/getInfo", laneController.GetLane)
	utils.MsgSuccess("    [SetupDataLaneRoutes]Successfully init!")
}
