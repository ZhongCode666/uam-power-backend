package main

import (
	"time"
	"uam-power-backend/service/data_transfer_service"
	"uam-power-backend/utils"
)

func main() {
	GlobalConfigPath := "config/global_config.yaml"
	GlobalCfg, loadCfgErr := utils.LoadGlobalConfig(GlobalConfigPath)
	if loadCfgErr != nil {
		utils.MsgError("[transfer]setting up area routes failed: >" + loadCfgErr.Error())
		return
	}
	// 初始化日志
	cfg, loadCfgErr := utils.LoadDBConfig(GlobalCfg.DBConfigPath)
	if loadCfgErr != nil {
		utils.MsgError("[transfer]setting up area routes failed: >" + loadCfgErr.Error())
		return
	}
	utils.MsgSuccess("[main_server]load DB config successfully!")

	var redisTransferGroup []*data_transfer_service.KafkaToRedis
	for i := 0; i < GlobalCfg.DataKafkaPartitionNum; i++ {
		redisTransferGroup = append(redisTransferGroup, data_transfer_service.NewKafkaToRedis(&cfg.KafkaCfg, &cfg.RedisCfg))
		redisTransferGroup[i].Start()
	}
	utils.MsgSuccess("[main_transfer_server]init transfer service successfully!")
	for true {
		time.Sleep(10 * time.Second)
		utils.MsgSuccess("[main_transfer_server]heart beating...")
	}
	utils.MsgError("[main_transfer_server]Quit!")
}
