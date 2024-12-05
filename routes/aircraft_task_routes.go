package routes

import (
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/controller/aircraft_task_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

// SetupAircraftTaskRoutes 设置 Aircraft Task 路由
// r 是 Fiber 应用实例
// RedisCfg 是 Redis 配置模型
// MySqlCfg 是 MySQL 配置模型
// ClickHouse 是 ClickHouse 配置模型
func SetupAircraftTaskRoutes(
	r *fiber.App, RedisCfg *db_config_model.RedisConfigModel,
	MySqlCfg *db_config_model.MySqlConfigModel,
	ClickHouse *db_config_model.ClickHouseConfigModel,
) {
	utils.MsgInfo("    [SetupAircraftTaskRoutes]setting up aircraft task routes...")
	// 创建一个新的 AircraftTaskController 实例
	aircraftTaskController := aircraft_task_controller.NewAircraftTaskModel(RedisCfg, MySqlCfg)
	// aircraftTaskController := aircraft_task_controller.NewAircraftTaskModelClickHouse(RedisCfg, MySqlCfg, ClickHouse)

	// 创建一个新的路由组 /aircraftTask
	uploadApis := r.Group("/aircraftTask")
	uploadApis.Post("/end", aircraftTaskController.EndTask)         // 结束任务
	uploadApis.Post("/create", aircraftTaskController.CreateTask)   // 创建任务
	uploadApis.Post("/check", aircraftTaskController.CheckTaskInfo) // 检查任务信息

	// 打印成功初始化信息
	utils.MsgSuccess("    [SetupAircraftTaskRoutes]Successfully init!")
}
