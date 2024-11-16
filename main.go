package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"uam-power-backend/routes"
	"uam-power-backend/service/data_transfer_service"
	"uam-power-backend/utils"
)

func main() {
	GlobalCfg, loadCfgErr := utils.LoadGlobalConfig("config/global_config.yaml")
	if loadCfgErr != nil {
		return
	}
	// 初始化日志
	utils.InitLog(GlobalCfg.SaveLogToFile)

	cfg, loadCfgErr := utils.LoadDBConfig(GlobalCfg.DBConfigPath)
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
	var redisTransferGroup []*data_transfer_service.KafkaToRedis
	for i := 0; i < GlobalCfg.KafkaPartitionNum; i++ {
		redisTransferGroup = append(redisTransferGroup, data_transfer_service.NewKafkaToRedis(&cfg.KafkaCfg, &cfg.RedisCfg))
		redisTransferGroup[i].Start()
	}
	var mysqlTransferGroup []*data_transfer_service.KafkaToMysql
	for i := 0; i < GlobalCfg.KafkaPartitionNum; i++ {
		mysqlTransferGroup = append(mysqlTransferGroup, data_transfer_service.NewKafkaToMysql(&cfg.KafkaCfg, &cfg.MySqlCfg, &cfg.RedisCfg))
		mysqlTransferGroup[i].Start()
	}

	utils.MsgSuccess("[main_server]init transfer service successfully!")
	// 启动服务器
	if err := r.Run(":" + strconv.Itoa(GlobalCfg.Port)); err != nil {
		utils.MsgError("[main_server]Failed to run the server: %v")
	}
}
