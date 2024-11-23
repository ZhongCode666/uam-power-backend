package global_model_config_model

type GlobalConfig struct {
	DBConfigPath      string `yaml:"DBConfigPath"`
	TaskPort          int    `yaml:"TaskPort"`
	IDPort            int    `yaml:"IDPort"`
	DataPort          int    `yaml:"DataPort"`
	LanePort          int    `yaml:"LanePort"`
	SaveLogToFile     bool   `yaml:"SaveLogToFile"`
	KafkaPartitionNum int    `yaml:"KafkaPartitionNum"`
}
