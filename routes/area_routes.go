package routes

import (
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/controller/area_data_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

func SetupAreaRoutes(
	r *fiber.App,
	mongoCfg *db_config_model.MongoConfigModel,
	mysqlCfg *db_config_model.MySqlConfigModel,
) {
	areaController := area_data_controller.NewAreaController(mongoCfg, mysqlCfg)

	recApis := r.Group("/area")
	recApis.Post("/create", areaController.CreateArea)
	recApis.Post("/delete", areaController.DeleteAreaData)
	recApis.Post("/getInfo", areaController.GetAreaData)
	recApis.Post("/updateOccupied", areaController.UpdateRasterDataOcc)
	recApis.Post("/updateOK", areaController.UpdateRasterDataOK)
	recApis.Post("/updateBarrier", areaController.UpdateRasterDataBarrier)
	recApis.Post("/updateBan", areaController.UpdateRasterDataBan)

	utils.MsgSuccess("    [SetupDataAreaRoutes]Successfully init!")
}
