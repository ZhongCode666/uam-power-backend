package main

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"uam-power-backend/routes"
	"uam-power-backend/utils"
)

func main() {
	GlobalConfigPath := "config/global_config.yaml"
	GlobalCfg, loadCfgErr := utils.LoadGlobalConfig(GlobalConfigPath)
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
	if err := app.Listen(":" + strconv.Itoa(GlobalCfg.DataPort)); err != nil {
		utils.MsgError("[main_server]Failed to run the server: %v")
	}
}
