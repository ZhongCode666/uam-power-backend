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
		return
	}
	cfg, loadCfgErr := utils.LoadDBConfig(GlobalCfg.DBConfigPath)
	if loadCfgErr != nil {
		return
	}
	utils.MsgSuccess("[main_server]load DB config successfully!")
	app := fiber.New(fiber.Config{
		// 设置请求体的最大大小（单位：字节）
		BodyLimit: 5 * 1024 * 1024 * 1024, // 50 MB
	})
	app.Use(cors.New())
	// 配置路由
	routes.SetupLaneRoutes(app, &cfg.MongoCfg, &cfg.MySqlCfg)
	utils.MsgSuccess("[main_server]init routes successfully!")
	// 启动服务器
	if err := app.Listen(":" + strconv.Itoa(GlobalCfg.LanePort)); err != nil {
		utils.MsgError("[main_server]Failed to run the server: %v")
	}
}
