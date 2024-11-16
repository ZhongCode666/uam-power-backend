package db_config_model

type DbConfigModel struct {
	KafkaCfg KafkaConfigModel `yaml:"KafkaCfg"`
	RedisCfg RedisConfigModel `yaml:"RedisCfg"`
	MySqlCfg MySqlConfigModel `yaml:"MySqlCfg"`
}
