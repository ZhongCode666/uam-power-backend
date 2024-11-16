package main

import (
	"github.com/gin-gonic/gin"
	"uam-power-backend/routes"
	"uam-power-backend/service/data_transfer_service"
	"uam-power-backend/utils"
)

func main() {
	// 初始化日志
	utils.InitLog()

	cfg, loadCfgErr := utils.LoadDBConfig("config/db_config.yaml")
	if loadCfgErr != nil {
		return
	}
	utils.MsgSuccess("[main_server]load DB config successfully!")
	// 创建一个新的Gin实例
	r := gin.Default()

	// 配置路由
	routes.SetupDataFlowRoutes(r, &cfg.KafkaCfg, &cfg.RedisCfg)
	routes.SetupAircraftTaskRoutes(r, &cfg.RedisCfg, &cfg.MySqlCfg)
	routes.SetupAircraftIdRoutes(r, &cfg.RedisCfg, &cfg.MySqlCfg)
	utils.MsgSuccess("[main_server]init routes successfully!")
	transferSer := data_transfer_service.NewKafkaToRedis(&cfg.KafkaCfg, &cfg.RedisCfg)
	transferSer.Start()
	transferSerMysql := data_transfer_service.NewKafkaToMysql(&cfg.KafkaCfg, &cfg.MySqlCfg, &cfg.RedisCfg)
	transferSerMysql.Start()
	utils.MsgSuccess("[main_server]init transfer service successfully!")
	// 启动服务器
	if err := r.Run(":26969"); err != nil {
		utils.MsgError("[main_server]Failed to run the server: %v")
	}
}
