package data_transfer_service

import (
	"os"
	"os/exec"
	"path/filepath"
	"uam-power-backend/utils"
)

func StartDataTransferService(ConfigPath string) {
	GlobalCfg, loadCfgErr := utils.LoadGlobalConfig(ConfigPath)
	if loadCfgErr != nil {
		return
	}
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
}
