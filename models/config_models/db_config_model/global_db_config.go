package db_config_model

// DbConfigModel 结构体表示数据库配置模型
type DbConfigModel struct {
	KafkaCfg      KafkaConfigModel      `yaml:"KafkaCfg"`      // Kafka 配置
	RedisCfg      RedisConfigModel      `yaml:"RedisCfg"`      // Redis 配置
	MySqlCfg      MySqlConfigModel      `yaml:"MySqlCfg"`      // MySQL 配置
	ClickHouseCfg ClickHouseConfigModel `yaml:"ClickHouseCfg"` // ClickHouse 配置
	MongoCfg      MongoConfigModel      `yaml:"MongoCfg"`      // MongoDB 配置
}
