package global_model_config_model

type GlobalConfig struct {
	DBConfigPath           string `yaml:"DBConfigPath"`
	TaskPort               int    `yaml:"TaskPort"`
	IDPort                 int    `yaml:"IDPort"`
	DataUploadPort         int    `yaml:"DataUploadPort"`
	DataReceivePort        int    `yaml:"DataReceivePort"`
	LanePort               int    `yaml:"LanePort"`
	AreaPort               int    `yaml:"AreaPort"`
	SaveLogToFile          bool   `yaml:"SaveLogToFile"`
	EventKafkaPartitionNum int    `yaml:"EventKafkaPartitionNum"`
	DataKafkaPartitionNum  int    `yaml:"DataKafkaPartitionNum"`
}
