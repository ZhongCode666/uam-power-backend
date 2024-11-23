package main

import (
	"uam-power-backend/service/data_transfer_service"
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

	//var redisTransferGroup []*data_transfer_service.KafkaToRedis
	//for i := 0; i < GlobalCfg.KafkaPartitionNum; i++ {
	//	redisTransferGroup = append(redisTransferGroup, data_transfer_service.NewKafkaToRedis(&cfg.KafkaCfg, &cfg.RedisCfg))
	//	redisTransferGroup[i].Start()
	//}
	var mysqlTransferGroup []*data_transfer_service.KafkaToMysql
	for i := 0; i < GlobalCfg.KafkaPartitionNum; i++ {
		mysqlTransferGroup = append(mysqlTransferGroup, data_transfer_service.NewKafkaToMysql(&cfg.KafkaCfg, &cfg.MySqlCfg, &cfg.RedisCfg))
		mysqlTransferGroup[i].Start()
	}
	var clickhouseTransferGroup []*data_transfer_service.KafkaToClickhouse
	for i := 0; i < GlobalCfg.KafkaPartitionNum; i++ {
		clickhouseTransferGroup = append(clickhouseTransferGroup, data_transfer_service.NewKafkaToClickhouse(&cfg.KafkaCfg, &cfg.ClickHouseCfg, &cfg.RedisCfg))
		clickhouseTransferGroup[i].Start()
	}
	utils.MsgSuccess("[main_transfer_server]init transfer service successfully!")
}
