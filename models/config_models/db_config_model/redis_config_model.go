package db_config_model

type RedisConfigModel struct {
	Port         int    `yaml:"Port"`
	StatusDBno   int    `yaml:"StatusDBno"`
	EventDBno    int    `yaml:"EventDBno"`
	Host         string `yaml:"Host"`
	AircraftDBno int    `yaml:"AircraftDBno"`
	TaskInfoDBno int    `yaml:"TaskInfoDBno"`
}
