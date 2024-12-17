package routes

import (
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/controller/lane_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

// SetupLaneRoutes 设置航线路由
// r 是 Fiber 应用实例
// mongoCfg 是 MongoDB 配置模型
// mysqlCfg 是 MySQL 配置模型
func SetupLaneRoutes(
	r *fiber.App,
	mongoCfg *db_config_model.MongoConfigModel,
	mysqlCfg *db_config_model.MySqlConfigModel,
) {
	utils.MsgInfo("    [SetupDataLaneRoutes]setting up lane routes...")
	// 创建一个新的 LaneController 实例
	laneController := lane_controller.NewLaneController(mongoCfg, mysqlCfg)

	// 创建一个新的路由组 /lane
	recApis := r.Group("/lane")
	recApis.Post("/create", laneController.CreateLane)         // 创建航线
	recApis.Post("/getInfo", laneController.GetLane)           // 获取航线信息
	recApis.Get("/getAllNotHidden", laneController.GetAllLane) // 获取所有未隐藏航线
	recApis.Get("/list", laneController.LaneList)              // 获取航线列表
	recApis.Post("/hide", laneController.HideLane)             // 隐藏航线
	recApis.Post("/show", laneController.ShowLane)             // 显示航线
	recApis.Post("/delete", laneController.DeleteLane)         // 删除航线
	recApis.Post("/getLaneInfo", laneController.QueryLaneInfo) // 获取航线信息

	// 打印成功初始化信息
	utils.MsgSuccess("    [SetupDataLaneRoutes]Successfully init!")
}
