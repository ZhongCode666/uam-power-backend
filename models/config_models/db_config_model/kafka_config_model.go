package db_config_model

type KafkaConfigModel struct {
	Addr               string `yaml:"Addr"`
	AircraftDataTopic  string `yaml:"AircraftDataTopic"`
	AircraftEventTopic string `yaml:"AircraftEventTopic"`
}
