package global_model_config_model

type GlobalConfig struct {
	DBConfigPath           string `yaml:"DBConfigPath"`           // 数据库配置路径
	TaskPort               int    `yaml:"TaskPort"`               // 任务端口
	IDPort                 int    `yaml:"IDPort"`                 // ID端口
	DataUploadPort         int    `yaml:"DataUploadPort"`         // 数据上传端口
	DataReceivePort        int    `yaml:"DataReceivePort"`        // 数据接收端口
	LanePort               int    `yaml:"LanePort"`               // 航线端口
	AreaPort               int    `yaml:"AreaPort"`               // 区域端口
	SaveLogToFile          bool   `yaml:"SaveLogToFile"`          // 保存日志到文件
	EventKafkaPartitionNum int    `yaml:"EventKafkaPartitionNum"` // 事件Kafka分区数
	DataKafkaPartitionNum  int    `yaml:"DataKafkaPartitionNum"`  // 数据Kafka分区数
}
