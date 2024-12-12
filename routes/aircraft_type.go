package routes

import (
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/controller/aircraft_type_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

func SetupAircraftTypeRoutes(
	r *fiber.App, MySqlCfg *db_config_model.MySqlConfigModel,
) {
	utils.MsgInfo("    [SetupAircraftTypeRoutes]setting up aircraft ID routes...")
	// 创建一个新的 AircraftIdController 实例
	aircraftTypeController := aircraft_type_controller.NewAircraftTypeController(MySqlCfg)

	// 创建一个新的路由组 /aircraftID
	typeApis := r.Group("/aircraftType")
	typeApis.Post("/info", aircraftTypeController.GetAircraftType)         // 获取 Aircraft 信息
	typeApis.Post("/create", aircraftTypeController.CreateAircraftType)    // 创建type
	typeApis.Post("/change", aircraftTypeController.ChangeType)            // 修改type
	typeApis.Post("/delete", aircraftTypeController.DeleteType)            // 删除type
	typeApis.Post("/infoByID", aircraftTypeController.GetTypeByAircraftID) // 根据飞机ID获取type

	// 打印成功初始化信息
	utils.MsgSuccess("    [SetupAircraftTypeRoutes]Successfully init!")
}
