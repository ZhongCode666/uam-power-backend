package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"strconv"
	"uam-power-backend/routes"
	"uam-power-backend/utils"
)

func main() {
	GlobalConfigPath := "config/global_config.yaml"
	GlobalCfg, loadCfgErr := utils.LoadGlobalConfig(GlobalConfigPath)
	if loadCfgErr != nil {
		utils.MsgError("[taskAPIs]setting up area routes failed: >" + loadCfgErr.Error())
		return
	}
	// 初始化日志
	cfg, loadCfgErr := utils.LoadDBConfig(GlobalCfg.DBConfigPath)
	if loadCfgErr != nil {
		utils.MsgError("[taskAPIs]setting up area routes failed: >" + loadCfgErr.Error())
		return
	}
	utils.MsgSuccess("[main_server]load DB config successfully!")

	// 创建一个新的Gin实例
	app := fiber.New()
	app.Use(cors.New())
	// 配置路由
	routes.SetupAircraftTaskRoutes(app, &cfg.RedisCfg, &cfg.MySqlCfg, &cfg.ClickHouseCfg)
	utils.MsgSuccess("[main_server]init routes successfully!")
	// 启动服务器
	if err := app.Listen(":" + strconv.Itoa(GlobalCfg.TaskPort)); err != nil {
		utils.MsgError("[main_server]Failed to run the server: %v")
	}
}
