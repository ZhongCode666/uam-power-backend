package main

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"uam-power-backend/routes"
	"uam-power-backend/utils"
)

func main() {
	GlobalConfigPath := "config/global_config.yaml"
	GlobalCfg, loadCfgErr := utils.LoadGlobalConfig(GlobalConfigPath)
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
	utils.MsgSuccess("[main_transfer_server]init transfer service successfully!")
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
