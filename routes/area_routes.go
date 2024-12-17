package routes

import (
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/controller/area_data_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

// SetupAreaRoutes 设置 Area 路由
// r 是 Fiber 应用实例
// mongoCfg 是 MongoDB 配置模型
// mysqlCfg 是 MySQL 配置模型
func SetupAreaRoutes(
	r *fiber.App,
	mongoCfg *db_config_model.MongoConfigModel,
	mysqlCfg *db_config_model.MySqlConfigModel,
) {
	utils.MsgInfo("    [SetupAreaRoutes]setting up area routes...")
	// 创建一个新的 AreaController 实例
	areaController := area_data_controller.NewAreaController(mongoCfg, mysqlCfg)

	// 创建一个新的路由组 /area
	recApis := r.Group("/area")
	recApis.Post("/create", areaController.CreateArea)                     // 创建区域
	recApis.Post("/generateAreaID", areaController.CreateAreaID)           // 生成 AreaID
	recApis.Post("/uploadAreaData", areaController.UploadArea)             // 上传区域数据
	recApis.Post("/delete", areaController.DeleteAreaData)                 // 删除区域数据
	recApis.Post("/getInfo", areaController.GetAreaData)                   // 获取区域数据
	recApis.Post("/updateOccupied", areaController.UpdateRasterDataOcc)    // 更新占用数据
	recApis.Post("/updateOK", areaController.UpdateRasterDataOK)           // 更新 OK 数据
	recApis.Post("/updateBarrier", areaController.UpdateRasterDataBarrier) // 更新障碍数据
	recApis.Post("/updateBan", areaController.UpdateRasterDataBan)         // 更新禁用数据
	recApis.Post("/rasterData", areaController.GetRasterData)              // 获取栅格数据

	// 打印成功初始化信息
	utils.MsgSuccess("    [SetupDataAreaRoutes]Successfully init!")
}
