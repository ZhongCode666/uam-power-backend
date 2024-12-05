package routes

import (
	"github.com/gofiber/fiber/v2"
	"uam-power-backend/controller/aircraft_id_controller"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/utils"
)

// SetupAircraftIdRoutes 设置 Aircraft ID 路由
// r 是 Fiber 应用实例
// RedisCfg 是 Redis 配置模型
// MySqlCfg 是 MySQL 配置模型
func SetupAircraftIdRoutes(
	r *fiber.App, RedisCfg *db_config_model.RedisConfigModel,
	MySqlCfg *db_config_model.MySqlConfigModel,
) {
	utils.MsgInfo("    [SetupAircraftIdRoutes]setting up aircraft ID routes...")
	// 创建一个新的 AircraftIdController 实例
	aircraftIDController := aircraft_id_controller.NewAircraftIdController(MySqlCfg, RedisCfg)

	// 创建一个新的路由组 /aircraftID
	uploadApis := r.Group("/aircraftID")
	uploadApis.Post("/info", aircraftIDController.GetAircraftInfo) // 获取 Aircraft 信息
	uploadApis.Post("/create", aircraftIDController.CreateUser)    // 创建用户

	// 打印成功初始化信息
	utils.MsgSuccess("    [AircraftIdRoutes]Successfully init!")
}
