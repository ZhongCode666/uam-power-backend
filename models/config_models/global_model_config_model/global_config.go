package global_model_config_model

type GlobalConfig struct {
	DBConfigPath      string `yaml:"DBConfigPath"`
	Port              int    `yaml:"Port"`
	SaveLogToFile     bool   `yaml:"SaveLogToFile"`
	KafkaPartitionNum int    `yaml:"KafkaPartitionNum"`
}
