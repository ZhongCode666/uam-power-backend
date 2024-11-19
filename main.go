package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"os/exec"
	"path/filepath"
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
	r := gin.Default()

	// 配置路由
	routes.SetupDataFlowRoutes(r, &cfg.KafkaCfg, &cfg.RedisCfg)
	routes.SetupAircraftTaskRoutes(r, &cfg.RedisCfg, &cfg.MySqlCfg)
	routes.SetupAircraftIdRoutes(r, &cfg.RedisCfg, &cfg.MySqlCfg)
	utils.MsgSuccess("[main_server]init routes successfully!")
	RedisTransfer, _ := filepath.Abs("service/data_transfer_service/kafka_to_redis.go")
	MySqlTransfer, _ := filepath.Abs("service/data_transfer_service/kafka_to_mysql.go")
	currentDir, _ := os.Getwd()
	for i := 0; i < GlobalCfg.KafkaPartitionNum; i++ {
		cmd := exec.Command("go", "run", RedisTransfer, GlobalCfg.DBConfigPath)
		cmd.Dir = currentDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			return
		}
	}

	for i := 0; i < GlobalCfg.KafkaPartitionNum; i++ {
		cmd := exec.Command("go", "run", MySqlTransfer, GlobalCfg.DBConfigPath)
		cmd.Dir = currentDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			return
		}
	}

	utils.MsgSuccess("[main_server]init transfer service successfully!")
	// 启动服务器
	if err := r.Run(":" + strconv.Itoa(GlobalCfg.Port)); err != nil {
		utils.MsgError("[main_server]Failed to run the server: %v")
	}
}
