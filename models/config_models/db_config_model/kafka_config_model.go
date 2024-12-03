package db_config_model

// KafkaConfigModel 结构体表示 Kafka 的配置模型
type KafkaConfigModel struct {
	Addr               string `yaml:"Addr"`               // 地址
	AircraftDataTopic  string `yaml:"AircraftDataTopic"`  // 航空器数据主题
	AircraftEventTopic string `yaml:"AircraftEventTopic"` // 航空器事件主题
	NumEventProducers  int    `yaml:"NumEventProducers"`  // 事件生产者数量
	NumDataProducers   int    `yaml:"NumDataProducers"`   // 数据生产者数量
}
