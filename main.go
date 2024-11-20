package main

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"uam-power-backend/routes"
	"uam-power-backend/utils"
)

func main() {
	GlobalCfg, loadCfgErr := utils.LoadGlobalConfig("config/global_config.yaml")
	if loadCfgErr != nil {
		return
	}
	// 初始化日志
	cfg, loadCfgErr := utils.LoadDBConfig(GlobalCfg.DBConfigPath)
	if loadCfgErr != nil {
		return
	}
	utils.MsgSuccess("[main_server]load DB config successfully!")
	// 创建一个新的Gin实例
	app := fiber.New()

	// 配置路由
	routes.SetupDataFlowRoutes(app, &cfg.KafkaCfg, &cfg.RedisCfg)
	routes.SetupAircraftTaskRoutes(app, &cfg.RedisCfg, &cfg.MySqlCfg)
	routes.SetupAircraftIdRoutes(app, &cfg.RedisCfg, &cfg.MySqlCfg)
	utils.MsgSuccess("[main_server]init routes successfully!")
	// 启动服务器
	if err := app.Listen(":" + strconv.Itoa(GlobalCfg.Port)); err != nil {
		utils.MsgError("[main_server]Failed to run the server: %v")
	}
}
